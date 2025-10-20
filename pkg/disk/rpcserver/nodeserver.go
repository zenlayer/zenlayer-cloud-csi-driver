/*
Copyright (C) 2025 Zenlayer, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rpcserver

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"golang.org/x/net/context"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/util/resizefs"
	k8smount "k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
)

type NodeServer struct {
	driver     *driver.DiskDriver
	cloud      cloud.CloudManager
	mounter    *mount.SafeFormatAndMount
	k8smounter k8smount.Interface
	locks      *common.ResourceLocks
}

var _ csi.NodeServer = &NodeServer{}

func NewNodeServer(d *driver.DiskDriver, c cloud.CloudManager, mnt *mount.SafeFormatAndMount) *NodeServer {
	return &NodeServer{
		driver:     d,
		cloud:      c,
		mounter:    mnt,
		k8smounter: k8smount.NewWithoutSystemd(""),
		locks:      common.NewResourceLocks(),
	}
}

/*
action:
This RPC is called by the CO when a workload that wants to use the specified volume is placed (scheduled) on a node. The Plugin SHALL assume that this RPC will be executed on the node where the volume will be used.
*/
func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	funcName := "NodeServer:NodePublishVolume:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume id missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Target path missing in request")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume capabilities missing in request")
	} else if !ns.driver.ValidateVolumeCapability(req.GetVolumeCapability()) {
		return nil, status.Error(codes.FailedPrecondition, ERRORLOG+"Exceed capabilities")
	}

	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.FailedPrecondition, ERRORLOG+"Stageing target path not set.")
	}

	targetPath := req.GetTargetPath()

	volId, serial, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	devicePath, err := ns.getDevPathBySerial(serial)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	stagePath := req.GetStagingTargetPath()

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer ns.locks.Release(volId)

	isBlockMode := req.GetVolumeCapability().GetBlock() != nil
	if isBlockMode {
		stagePath = devicePath
	}

	fsType := "ext4"
	mnt := req.GetVolumeCapability().GetMount()
	if mnt != nil {
		if mnt.GetFsType() != "" {
			fsType = mnt.FsType
		}
	}
	err = checkfsType(fsType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport fsType "+fsType)
	}

	if err := ns.createTargetMountPathIfNotExists(targetPath, isBlockMode); err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}

	notMnt, err := ns.mounter.IsNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}
	if !notMnt {
		klog.Infof("%s TargetPath[%s] is mounted, not need mount again", INFOLOG, targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// set bind mount options
	options := []string{"bind"}
	if req.GetReadonly() {
		options = append(options, "ro")
	}
	if !isBlockMode && mnt != nil {
		for _, item := range common.FeatureGates {
			if item == common.FEAT_MOUNTOPT_ENABLE {
				options = append(options, mnt.MountFlags...)
			}
		}
	}

	//filesystem: stage阶段格式化文件驱动会把卷先挂载或准备到这个stagingpath,通常是节点级别的中间准备路径不特定于某个Pod. publish阶段将已经在staging或者设备已准备好的资源“发布”到 Pod 可访问的位置，也即Pod容器内部所看到的路径
	//filesystem: stagePath=/var/lib/kubelet/plugins/kubernetes.io/csi/disk.csi.zenlayer.com/ce05eddc643b4c39b0469fedc3139363254a87e46f5a367e3b52abdb7b78e99c/globalmount
	//filesystem: targetPath=/var/lib/kubelet/pods/926e3fe7-8f79-4cb3-b2e6-e7cd8e449f43/volumes/kubernetes.io~csi/zeccsi-pv-dc39b10e-048e-4eb6-bc18-1b57cab76750/mount

	//blockmode: stage阶段准备裸设备，publish阶段将裸设备暴露到targetpath,可能的方式是创建一个设备文件链接到/dev/vdb，并把它放在target_path
	//blockmode: stagePath=/dev/vdb
	//blockmode: targetPath=/var/lib/kubelet/plugins/kubernetes.io/csi/volumeDevices/publish/zeccsi-pv-061c4b2e-afaa-4595-8973-662cf9edf677/ffd648f9-0be5-40b1-833e-c183aec7624a
	klog.Infof("%s Will mount src stagePath[%s] to dst targetPath[%s], blockmode[%v], fstype[%s], options[%v], vol[%s]", INFOLOG, stagePath, targetPath, isBlockMode, fsType, options, volId)
	if err := ns.mounter.Mount(stagePath, targetPath, fsType, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Success mount src stagePath[%s] to dst targetPath[%s], blockmode[%v], fstype[%s], options[%v], vol[%s]", INFOLOG, stagePath, targetPath, isBlockMode, fsType, options, volId)

	return &csi.NodePublishVolumeResponse{}, nil
}

/*
action:
This RPC MUST undo the work by the corresponding NodePublishVolume. This RPC SHALL be called by the CO at least once for each target_path that was successfully setup via NodePublishVolume.
If the corresponding Controller Plugin has PUBLISH_UNPUBLISH_VOLUME controller capability, the CO SHOULD issue all NodeUnpublishVolume (as specified above) before calling ControllerUnpublishVolume
for the given node and the given volume. The Plugin SHALL assume that this RPC will be executed on the node where the volume is being used.
*/
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	funcName := "NodeServer:NodeUnpublishVolume:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Target path missing in request.")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume Id missing in request.")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	targetPath := req.GetTargetPath()

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer ns.locks.Release(volId)

	klog.Infof("%s Will umount volume[%s], targetpath[%v]", INFOLOG, req.VolumeId, targetPath)
	err = mount.CleanupMountPoint(targetPath, ns.mounter.Interface, true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unmount targetpath[%s] error[%v]", targetPath, err)
	}

	klog.Infof("%s Success umount volume[%s], targetpath[%v]", INFOLOG, req.VolumeId, targetPath)
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

/*
action:
This RPC is called by the CO prior to the volume being consumed by any workloads on the node by NodePublishVolume.
The Plugin SHALL assume that this RPC will be executed on the node where the volume will be used.
This RPC SHOULD be called by the CO when a workload that wants to use the specified volume is placed (scheduled) on the specified node for
the first time or for the first time since a NodeUnstageVolume call for the specified volume was called and returned success on that node.
*/
func (ns *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	funcName := "NodeServer:NodeStageVolume:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if flag := ns.driver.ValidateNodeServiceRequest(csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME); !flag {
		return nil, status.Error(codes.Unimplemented, ERRORLOG+"Node has not stage capability")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Target path missing in request")
	}
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume capability missing in request")
	}

	volId, serial, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	targetPath := req.GetStagingTargetPath()

	// skip staging if volume is in block mode
	if req.GetVolumeCapability().GetBlock() != nil {
		klog.Infof("%s Skipping staging of volume[%s] on targetpath[%s] since it's in block mode", INFOLOG, req.GetVolumeId(), targetPath)
		return &csi.NodeStageVolumeResponse{}, nil
	}

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer ns.locks.Release(volId)

	fsType := "ext4"
	var options []string
	mnt := req.GetVolumeCapability().GetMount()
	if mnt == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"GetMount return nil")
	}
	if mnt.GetFsType() != "" {
		fsType = mnt.FsType
	}

	err = checkfsType(fsType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport fsType "+fsType)
	}

	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		return &csi.NodeStageVolumeResponse{}, nil
	}

	devicePath, err := ns.getDevPathBySerial(serial)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	for _, item := range common.FeatureGates {
		if item == common.FEAT_MOUNTOPT_ENABLE {
			options = append(options, mnt.MountFlags...)
		}
	}

	mkfsOptions := make([]string, 0)
	omitfsck := false

	diskMounter := &k8smount.SafeFormatAndMount{Interface: ns.k8smounter, Exec: utilexec.New()}

	klog.Infof("%s Will mount [%s] to [%s], vol[%s], fstype[%s], req options[%v], used options[%v], featureGates[%v]", INFOLOG, devicePath, targetPath, volId, fsType, mnt.MountFlags, options, common.FeatureGates)
	if err := common.FormatAndMount(diskMounter, devicePath, targetPath, fsType, mkfsOptions, options, omitfsck); err != nil {
		klog.Errorf("%s Mountdevice: FormatAndMount fail with mkfsOptions [%s], [%s], [%s], [%s], [%s] with error[%s]", ERRORLOG, devicePath, targetPath, fsType, mkfsOptions, options, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Success mount [%s] to [%s], vol[%s], fstype[%s], req options[%v], used options[%v], featureGates[%v]", INFOLOG, devicePath, targetPath, volId, fsType, mnt.MountFlags, options, common.FeatureGates)

	r := k8smount.NewResizeFs(diskMounter.Exec)
	needResize, err := r.NeedResize(devicePath, targetPath)
	if err != nil {
		klog.Errorf("%s Could not determine if volume[%s] need to be resized[%v], devicepath[%s], targetpath[%s]", ERRORLOG, req.VolumeId, err, devicePath, targetPath)
		return &csi.NodeStageVolumeResponse{}, nil
	}
	if needResize {
		klog.Infof("%s Will resizing volume[%s], devicepath[%s], targetpath[%s]", INFOLOG, req.VolumeId, devicePath, targetPath)
		if _, err := r.Resize(devicePath, targetPath); err != nil {
			return nil, status.Errorf(codes.Internal, "%s Could not resize volume[%s], error[%v]", ERRORLOG, req.VolumeId, err)
		}
	}

	return &csi.NodeStageVolumeResponse{}, nil
}

/*
action:
This RPC is a reverse operation of NodeStageVolume. This RPC MUST undo the work by the corresponding NodeStageVolume. \
This RPC SHALL be called by the CO once for each staging_target_path that was successfully setup via NodeStageVolume.
*/
func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	funcName := "NodeServer:NodeUnstageVolume:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if flag := ns.driver.ValidateNodeServiceRequest(csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME); !flag {
		return nil, status.Error(codes.Unimplemented, ERRORLOG+"Node has not unstage capability")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Target path missing in request")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	targetPath := req.GetStagingTargetPath()

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer ns.locks.Release(volId)

	notMnt, err := ns.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}
	if notMnt {
		return &csi.NodeUnstageVolumeResponse{}, nil
	}
	// count mount point
	_, cnt, err := mount.GetDeviceNameFromMount(ns.mounter, targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}

	klog.Infof("%s Will umount targetPath[%s]", INFOLOG, targetPath)
	err = ns.mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}
	klog.Infof("%s Success umount targetPath[%s]. before umount cnt[%d]", INFOLOG, targetPath, cnt)
	cnt--

	if cnt > 0 {
		klog.Errorf("%s Volume[%s] still mounted in instance[%s]", ERRORLOG, volId, ns.driver.GetNodeId())
		return nil, status.Error(codes.Internal, "unmount failed")
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

/*
Action: A Node Plugin MUST implement this RPC call if the plugin has PUBLISH_UNPUBLISH_VOLUME controller capability.

	The Plugin SHALL assume that this RPC will be executed on the node where the volume will be used.
	The CO SHOULD call this RPC for the node at which it wants to place the workload. The CO MAY call this RPC more than once for a given node.
	The SP SHALL NOT expect the CO to call this RPC more than once. The result of this call will be used by CO in ControllerPublishVolume.
*/
func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	funcName := "NodeServer:NodeGetInfo:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	vminfo := cloud.NewZecVm()
	vminfo.ZecVm_Type = driver.BasicVmType.Int()
	vminfo.ZecVm_Zone = ns.driver.GetNodeZone()

	//vm only has one type BasicVm
	vmType, ok := driver.VmTypeName[driver.VmType(vminfo.ZecVm_Type)]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "%s vm type error.", ERRORLOG)
	}

	topo := &csi.Topology{
		Segments: map[string]string{
			ns.driver.GetTopologyVmTypeKey(): vmType,
			ns.driver.GetTopologyZoneKey():   vminfo.ZecVm_Zone,
		},
	}

	_ = INFOLOG
	return &csi.NodeGetInfoResponse{
		NodeId:             ns.driver.GetNodeId(),
		MaxVolumesPerNode:  ns.driver.GetMaxVolumePerNode(),
		AccessibleTopology: topo,
	}, nil
}

/*
Action: This RPC call allows CO to expand volume on a node.
*/
func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {

	funcName := "NodeServer:NodeExpandVolume:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if req.VolumeCapability != nil && req.VolumeCapability.GetBlock() != nil {
		klog.Infof("%s skipping expand for block volume, req.GetVolumeId[%s]", INFOLOG, req.GetVolumeId())
		return &csi.NodeExpandVolumeResponse{}, nil
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume ID missing in request")
	}
	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume path missing in ad request")
	}

	reqSizeBytes, err := common.GetRequestSizeBytes(req.GetCapacityRange())
	if err != nil {
		return nil, status.Error(codes.OutOfRange, ERRORLOG+err.Error())
	}

	volId, serial, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	volPath := req.GetVolumePath()

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer ns.locks.Release(volId)

	devicePath, err := ns.getDevPathBySerial(serial)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	klog.Infof("%s Will Resize file system device[%s], mount path[%s], blockPath[%s], request size[%d]Byte", INFOLOG, devicePath, volPath, devicePath, reqSizeBytes)
	resizer := resizefs.NewResizeFs(ns.mounter)

	ok, err := resizer.Resize(devicePath, volPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}
	if !ok {
		return nil, status.Error(codes.Internal, ERRORLOG+"fail to expand vol fs")
	}

	blkSizeBytes, err := ns.getBlockSizeBytes(devicePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s get vol block size err, path=%s, err=%v.", ERRORLOG, devicePath, err)
	}
	klog.Infof("%s Success Resize file system device[%s], mount path[%s], blockPath[%s], Block size[%d]Byte, request size[%d]Byte", INFOLOG, devicePath, volPath, devicePath, blkSizeBytes, reqSizeBytes)

	if blkSizeBytes < reqSizeBytes {
		// It's possible that the somewhere the volume size was rounded up, getting more size than requested is a success
		return nil, status.Errorf(codes.Internal, "%s resize requested for [%v] but after resize volume was size [%v]",
			ERRORLOG, reqSizeBytes, blkSizeBytes)
	}

	return &csi.NodeExpandVolumeResponse{
		CapacityBytes: blkSizeBytes,
	}, nil
}

/*
Action: This RPC allows the CO to check the supported capabilities of node service provided by the Plugin.
*/
func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.
	NodeGetCapabilitiesResponse, error) {

	funcName := "NodeServer:NodeGetCapabilities:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: ns.driver.GetNodeCapability(),
	}, nil
}

/*
Action: NodeGetVolumeStats RPC call returns the volume capacity statistics available for the volume.
*/
func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {

	funcName := "NodeServer:NodeGetVolumeStats:"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	INFOLOG := "INFO:" + funcName + hash + " "
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume path missing in request")
	}

	volumePath := req.GetVolumePath()
	// block mode volume's stats can't be retrieved like those filesystem volumes
	pathType, err := ns.mounter.GetFileType(volumePath)
	if err != nil {
		klog.Errorf("%s GetFileType error, volumePath[%s], error[%v]", ERRORLOG, volumePath, err.Error())
		return nil, status.Errorf(codes.NotFound, "failed to GetFileType volumePath[%s], err[%v]", volumePath, err)
	}
	isBlockMode := pathType == mount.FileTypeBlockDev

	if isBlockMode {
		blockSize, err := ns.getBlockSizeBytes(volumePath)
		if err != nil {
			klog.Errorf("%s getBlockSizeBytes() error, volumePath[%s], error[%v]", ERRORLOG, volumePath, err.Error())
			return nil, status.Errorf(codes.Internal, "failed to get block capacity on path[%s], err[%v]", volumePath, err)
		}
		klog.Infof("%s getBlockSizeBytes() get path[%s], pathType[%s], BlockMode[%v], totalGB[%d]", INFOLOG, volumePath, pathType, isBlockMode, blockSize/common.Gib)
		return &csi.NodeGetVolumeStatsResponse{
			Usage: []*csi.VolumeUsage{
				{
					Unit:  csi.VolumeUsage_BYTES,
					Total: blockSize,
				},
			},
		}, nil
	}

	var stat unix.Statfs_t
	err = unix.Statfs(volumePath, &stat)
	if err != nil {
		klog.Errorf("%s Statfs() error, volumePath[%s], error[%v]", ERRORLOG, volumePath, err.Error())
		return nil, status.Errorf(codes.Internal, "fail to get volume[%s] status, err[%v]", volumePath, err)
	}

	totalB := int64(stat.Blocks * uint64(stat.Bsize))
	freeB := int64(stat.Bfree * uint64(stat.Bsize))
	availB := int64(stat.Bavail * uint64(stat.Bsize))
	klog.Infof("%s Statfs() get path[%s], pathType[%s], BlockMode[%v], totalGB[%.1f], availableGB[%.1f], usedGB[%.1f]", INFOLOG, volumePath, pathType, isBlockMode, float64(totalB)/(float64)(common.Gib), float64(availB)/(float64)(common.Gib), float64(totalB-freeB)/(float64)(common.Gib))
	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: availB,
				Total:     totalB,
				Used:      totalB - freeB,
				Unit:      csi.VolumeUsage_BYTES,
			},
		},
	}, nil
}

func checkfsType(fsType string) error {
	if fsType == "ext4" || fsType == "xfs" || fsType == "ext3" {
		return nil
	}
	return fmt.Errorf("ERROR: unsupport fsType %s", fsType)
}

func (ns *NodeServer) getDevPathBySerial(serial string) (devpath string, err error) {

	devpath = ""

	output, err := ns.mounter.Exec.Run("lsblk", "-d", "-o", "NAME", "-n", "-i")
	if err != nil {
		return "", fmt.Errorf("ERROR:getDevPathBySerial lsblk cmd error[%s], out[%s]"+err.Error(), output)
	}

	devlist := strings.Split(string(output[:]), "\n")

	for _, dev := range devlist {
		if dev != "" {
			devserial_path := "/sys/block/" + dev + "/device/block/" + dev + "/serial"
			devserial, err := os.ReadFile(devserial_path)
			if err != nil {
				klog.Errorf("ERROR:getDevPathBySerial read[%s] file error, dev[%s]", devserial_path, dev)
				continue
			}

			klog.Infof("INFO:getDevPathBySerial list dev[%s], serial[%s], target-serial[%s], list count[%d]", dev, devserial, serial, len(devlist)-1)
			if string(devserial[:]) == serial {
				if devpath != "" {
					return "", fmt.Errorf("ERROR:serial[%s/%s] dup", devserial, serial)
				}
				devpath = "/dev/" + dev
			}
		}
	}

	if devpath == "" {
		return "", fmt.Errorf("ERROR:can not find block device, serial=%s", serial)
	}

	klog.Infof("INFO:getDevPathBySerial return path[%s]", devpath)
	return devpath, nil
}

func (ns *NodeServer) getBlockSizeBytes(devicePath string) (int64, error) {

	output, err := ns.mounter.Exec.Run("blockdev", "--getsize64", devicePath)
	if err != nil {
		return -1, fmt.Errorf("ERROR:getBlockSizeBytes blockdev cmd error, path[%s], output[%s], err[%v]", devicePath, string(output), err)
	}
	strOut := strings.TrimSpace(string(output))
	gotSizeBytes, err := strconv.ParseInt(strOut, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse size %s into int a size", strOut)
	}
	return gotSizeBytes, nil
}

// createTargetMountPathIfNotExists creates the mountPath if it doesn't exist
// if in block volume mode, a file will be created
func (ns *NodeServer) createTargetMountPathIfNotExists(mountPath string, isBlockMode bool) error {
	exists, err := ns.mounter.ExistsPath(mountPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if isBlockMode {
		pathFile, err := os.OpenFile(mountPath, os.O_CREATE|os.O_RDWR, 0750)
		if err != nil {
			return err
		}
		if err = pathFile.Close(); err != nil {
			return err
		}
	} else {
		// Create a directory
		if err := os.MkdirAll(mountPath, 0750); err != nil {
			return err
		}
	}

	return nil
}
