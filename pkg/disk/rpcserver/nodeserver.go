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
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
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

	if req.GetVolumeCapability().GetBlock() != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport Block Mode")
	}

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

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	stagePath := req.GetStagingTargetPath()

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer ns.locks.Release(volId)

	isBlockMode := false

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
		klog.Infof("%s TargetPath(%s) is mounted, not need mount again", INFOLOG, targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// set bind mount options
	options := []string{"bind"}
	if req.GetReadonly() {
		options = append(options, "ro")
	}
	for _, item := range common.FeatureGates {
		if item == common.FEAT_MOUNTOPT_ENABLE {
			options = append(options, mnt.MountFlags...)
		}
	}

	klog.Infof("%s Will mount src:%s to dst:%s", INFOLOG, stagePath, targetPath)
	if err := ns.mounter.Mount(stagePath, targetPath, fsType, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Bind mount %s at %s, fsType=%s, args options=%v, used options=%v, featureGates=%v", INFOLOG, stagePath, targetPath, fsType, mnt.MountFlags, options, common.FeatureGates)

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

	targetPath := req.GetTargetPath() ///var/lib/kubelet/pods/fc80f807-f73c-426e-85bf-e219c36c7973/volumes/kubernetes.io~csi/pvc-b6e33fdf-051a-4657-812f-e9967f3da0a2/mount

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer ns.locks.Release(volId)

	err = mount.CleanupMountPoint(targetPath, ns.mounter.Interface, true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unmount target path %s error: %v", targetPath, err)
	}

	klog.Infof("%s: Umount Successful for volume %s, target %v", INFOLOG, req.VolumeId, targetPath)
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

	if req.GetVolumeCapability().GetBlock() != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport Block Mode")
	}

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
	targetPath := req.GetStagingTargetPath() ///var/lib/kubelet/plugins/kubernetes.io/csi/disk.csi.zenlayer.com/92b9024f4b7dd04f9f46fcdec625bd19302dd7f286971330a5dec9fa170bbf39/globalmount

	klog.Infof("%s volumeid=%s, req.GetStagingTargetPath()=%s", INFOLOG, volId, req.GetStagingTargetPath())

	if acquired := ns.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer ns.locks.Release(volId)

	fsType := "ext4"
	var options []string
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
	if err := common.FormatAndMount(diskMounter, devicePath, targetPath, fsType, mkfsOptions, options, omitfsck); err != nil {
		klog.Errorf("%s Mountdevice: FormatAndMount fail with mkfsOptions %s, %s, %s, %s, %s with error: %s", ERRORLOG, devicePath, targetPath, fsType, mkfsOptions, options, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Mount %s to %s succeed, vol=%s, fstype=%s, args options=%v, used options=%v, featureGates=%v", INFOLOG, devicePath, targetPath, volId, fsType, mnt.MountFlags, options, common.FeatureGates)

	r := k8smount.NewResizeFs(diskMounter.Exec)
	needResize, err := r.NeedResize(devicePath, targetPath)
	if err != nil {
		klog.Infof("%s: Could not determine if volume %s need to be resized: %v", INFOLOG, req.VolumeId, err)
		return &csi.NodeStageVolumeResponse{}, nil
	}
	if needResize {
		klog.Infof("%s: Resizing volume %q created from a volume", INFOLOG, req.VolumeId)
		if _, err := r.Resize(devicePath, targetPath); err != nil {
			return nil, status.Errorf(codes.Internal, "Could not resize volume %s: %v", req.VolumeId, err)
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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
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

	err = ns.mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error())
	}

	klog.Infof("%s targetPath %s has been umount.before umount cnt=%d", INFOLOG, targetPath, cnt)
	cnt--

	if cnt > 0 {
		klog.Errorf("%s Volume %s still mounted in instance %s", ERRORLOG, volId, ns.driver.GetNodeId())
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

	if req.GetVolumeCapability().GetBlock() != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport Block Mode")
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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer ns.locks.Release(volId)

	devicePath, err := ns.getDevPathBySerial(serial)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

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
	klog.Infof("%s Resize file system device %s, mount path %s, blockPath=%s, Block size %d Byte, request size %d Byte.", INFOLOG, devicePath, volPath, devicePath, blkSizeBytes, reqSizeBytes)

	if blkSizeBytes < reqSizeBytes {
		// It's possible that the somewhere the volume size was rounded up, getting more size than requested is a success
		return nil, status.Errorf(codes.Internal, "%s resize requested for %v but after resize volume was size %v",
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
		klog.Errorf("%s GetFileType error, volumePath=%s, error=%v", ERRORLOG, volumePath, err.Error())
		return nil, status.Errorf(codes.NotFound, "failed to determine volume %s 's mode: %v", volumePath, err)
	}
	isBlockMode := pathType == mount.FileTypeBlockDev

	if isBlockMode {
		blockSize, err := ns.getBlockSizeBytes(volumePath)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get block capacity on path %s: %v", volumePath, err)
		}
		klog.Infof("%s path=%s, pathType=%s, BlockMode=%v, totalGB=%d", INFOLOG, volumePath, pathType, isBlockMode, blockSize/common.Gib)
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
		klog.Errorf("%s Statfs error, volumePath=%s, error=%v", ERRORLOG, volumePath, err.Error())
		return nil, status.Errorf(codes.Internal, "fail to get volume %s status, err=%v", volumePath, err)
	}

	totalB := int64(stat.Blocks * uint64(stat.Bsize))
	freeB := int64(stat.Bfree * uint64(stat.Bsize))
	klog.Infof("%s path=%s, pathType=%s, BlockMode=%v, totalGB=%d, availableGB=%d", INFOLOG, volumePath, pathType, isBlockMode, totalB/common.Gib, freeB/common.Gib)
	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: freeB,
				Total:     totalB,
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
		return "", fmt.Errorf("ERROR:getDevPathBySerial lsblk cmd error:%s, out:%s"+err.Error(), output)
	}

	devlist := strings.Split(string(output[:]), "\n")

	for _, dev := range devlist {
		if dev != "" {
			devserial_path := "/sys/block/" + dev + "/device/block/" + dev + "/serial"
			devserial, err := os.ReadFile(devserial_path)
			if err != nil {
				klog.Errorf("ERROR:getDevPathBySerial read %s file error, dev:%s", devserial_path, dev)
				return "", err
			}

			klog.Infof("INFO:getDevPathBySerial list dev:%s, serial:%s, target-serial:%s, list count=%d", dev, devserial, serial, len(devlist)-1)
			if string(devserial[:]) == serial {
				if len(devpath) != 0 {
					return "", fmt.Errorf("ERROR:serial dup")
				}
				devpath = "/dev/" + dev
			}
		}
	}

	if devpath == "" {
		return "", fmt.Errorf("ERROR:can not find block device, serial=%s", serial)
	}

	klog.Infof("INFO:getDevPathBySerial return path=%s", devpath)
	return devpath, nil
}

func (ns *NodeServer) getBlockSizeBytes(devicePath string) (int64, error) {

	output, err := ns.mounter.Exec.Run("blockdev", "--getsize64", devicePath)
	if err != nil {
		return -1, fmt.Errorf("ERROR:getBlockSizeBytes blockdev cmd error, path %s: output: %s, err: %v", devicePath, string(output), err)
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

func (ns *NodeServer) unmountDuplicateMountPoint(targetPath, volumeId string) error {
	klog.Infof("unmountDuplicateMountPoint: start to unmount remains data: %s", targetPath)
	pathParts := strings.Split(targetPath, "/")
	partsLen := len(pathParts)
	if partsLen > 2 && pathParts[partsLen-1] == "mount" {
		var err error
		// V1 used in Kubernetes 1.23 and earlier
		globalPathV1 := filepath.Join(common.KubeletRootDir, "plugins/kubernetes.io/csi/pv/", pathParts[partsLen-2], "/globalmount")
		// V2 used in Kubernetes 1.24 and later, see https://github.com/kubernetes/kubernetes/pull/107065
		// If a volume is mounted at globalPathV1 then kubelet is upgraded, kubelet will also mount the same volume at globalPathV2.
		volSha := fmt.Sprintf("%x", sha256.Sum256([]byte(volumeId)))
		globalPathV2 := filepath.Join(common.KubeletRootDir, "plugins/kubernetes.io/csi/", common.DefaultProvisionName, volSha, "/globalmount")

		v1Exists := common.IsFileExisting(globalPathV1)
		v2Exists := common.IsFileExisting(globalPathV2)
		// Community requires the node to be drained before upgrading, but we do not. So clean the V1 mountpoint here if both exists.
		if v1Exists && v2Exists {
			klog.Info("unmountDuplicateMountPoint: oldPath & newPath exists at same time")
			err = ns.forceUnmountPath(globalPathV1)
		}

		// Now we have either V1 or V2 mountpoint.
		// Unmount it if it is propagated to data disk, or kubelet with version < 1.26 will refuse to unstage the volume.
		// Unmounting may also be propagated back to KubeletRootDir, we will fix that in NodePublishVolume.
		globalPath2 := filepath.Join("/var/lib/container/kubelet/plugins/kubernetes.io/csi/pv/", pathParts[partsLen-2], "/globalmount")
		globalPath3 := filepath.Join("/var/lib/container/kubelet/plugins/kubernetes.io/csi/", common.DefaultProvisionName, volSha, "/globalmount")
		if common.IsFileExisting(globalPath2) {
			err = ns.unmountDuplicationPath(globalPath2)
		}
		if common.IsFileExisting(globalPath3) {
			err = ns.unmountDuplicationPath(globalPath3)
		}
		return err
	} else {
		klog.Warningf("Target Path is illegal format: %s", targetPath)
	}
	return nil
}

func (ns *NodeServer) unmountDuplicationPath(globalPath string) error {
	// check globalPath2 is mountpoint
	notmounted, err := ns.k8smounter.IsLikelyNotMountPoint(globalPath)
	if err != nil || notmounted {
		klog.Warningf("Global Path is not mounted: %s", globalPath)
		return nil
	}
	// check device is used by others
	refs, err := ns.k8smounter.GetMountRefs(globalPath)
	if err == nil && !ns.HasMountRefs(globalPath, refs) {
		klog.Infof("NodeUnpublishVolume: VolumeId Unmount global path %s for ack with kubelet data disk", globalPath)
		if err := ns.k8smounter.Unmount(globalPath); err != nil {
			klog.Errorf("NodeUnpublishVolume: volumeId: unmount global path %s failed with err: %v", globalPath, err)
			return status.Error(codes.Internal, err.Error())
		}
	} else {
		klog.Infof("Global Path %s is mounted by others: %v", globalPath, refs)
	}
	return nil
}

func (ns *NodeServer) forceUnmountPath(globalPath string) error {
	// check globalPath2 is mountpoint
	notmounted, err := ns.k8smounter.IsLikelyNotMountPoint(globalPath)
	if err != nil || notmounted {
		klog.Warningf("Global Path is not mounted: %s", globalPath)
		return nil
	}
	if err := ns.k8smounter.Unmount(globalPath); err != nil {
		klog.Errorf("NodeUnpublishVolume: volumeId: unmount global path %s failed with err: %v", globalPath, err)
		return status.Error(codes.Internal, err.Error())
	} else {
		klog.Infof("forceUmountPath: umount Global Path %s  successful", globalPath)
	}
	return nil
}

func (m *NodeServer) HasMountRefs(mountPath string, mountRefs []string) bool {
	// Copied from https://github.com/kubernetes/kubernetes/blob/53902ce5ede4/pkg/volume/util/util.go#L680-L706
	pathToFind := mountPath
	if i := strings.Index(mountPath, common.KubernetesPluginPathPrefix); i > -1 {
		pathToFind = mountPath[i:]
	}
	for _, ref := range mountRefs {
		if !strings.Contains(ref, pathToFind) {
			return true
		}
	}
	return false
}
