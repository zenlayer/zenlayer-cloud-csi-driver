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
	"reflect"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/klog"
)

type ControllerServer struct {
	driver        *driver.DiskDriver
	cloud         cloud.CloudManager
	locks         *common.ResourceLocks
	detachLimiter common.RetryLimiter
}

func NewControllerServer(d *driver.DiskDriver, c cloud.CloudManager, maxRetry int) *ControllerServer {
	return &ControllerServer{
		driver:        d,
		cloud:         c,
		locks:         common.NewResourceLocks(),
		detachLimiter: common.NewRetryLimiter(maxRetry),
	}
}

var _ csi.ControllerServer = &ControllerServer{}

/*
action: CSI operation create zec cloud disk

	This operation MAY create three types of volumes:
		1. Empty volumes: CREATE_DELETE_VOLUME
		2. Restore volume from snapshot: CREATE_DELETE_VOLUME and CREATE_DELETE_SNAPSHOT
		3. Clone volume: CREATE_DELETE_VOLUME and CLONE_VOLUME

args: ctx context.Context, req *csi.CreateVolumeRequest

return: *csi.CreateVolumeResponse, error
*/
func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	funcName := "ControllerServer:CreateVolume:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); !isValid {
		return nil, status.Error(codes.Unimplemented, ERRORLOG+"unsupported controller server capability "+",volname="+req.GetName())
	}

	if req.GetVolumeCapabilities() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"volume capabilities missing in request "+",volname="+req.GetName())
	} else if !cs.driver.ValidateVolumeCapabilities(req.GetVolumeCapabilities()) {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"volume capabilities not match "+",volname="+req.GetName())
	}

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"volume name missing in request ")
	}
	if len(req.GetName()) > 64 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"volume name is longer then 64")
	}

	volName := req.GetName()

	if acquired := cs.locks.TryAcquire(volName); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volName)
	}
	klog.Infof("%s succ lock resource [%s]", INFOLOG, volName)

	defer klog.Infof("%s succ unlock resource [%s]", INFOLOG, volName)
	defer cs.locks.Release(volName)

	//read conf and init storage-class
	sc, err := driver.NewZecStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+err.Error())
	}

	topo := &driver.Topology{}
	if req.GetAccessibilityRequirements() != nil && cs.driver.ValidatePluginCapabilityService(csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS) {
		klog.Infof("%s GetAccessibilityRequirements has val. volName[%s]", INFOLOG, volName)
		var err error
		topo, err = cs.PickTopology(req.GetAccessibilityRequirements())
		if err != nil || topo == nil {
			return nil, status.Error(codes.InvalidArgument, ERRORLOG+err.Error())
		}
		//这里需要将req.Parameters中的zoneID修改成topo.ZoneID,req是从storageclass.yaml中取的，如果是WaitForFirstConsumer模式走进这个分支storageclass中定义的zone不一定是这个vol选择的zone
		driver.UpdateParmsZone(req.GetParameters(), topo.GetZone())
	} else {
		klog.Infof("%s GetAccessibilityRequirements is nil, use storage-class config zone[%s], volname[%s]", INFOLOG, sc.GetZone(), volName)
		//only support one Vm type (BasicVm)
		topo = driver.NewTopology(sc.GetZone(), driver.BasicVmType)
	}

	// get request volume capacity range
	requiredSizeByte, err := sc.GetRequiredVolumeSizeByte(req.GetCapacityRange())
	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "%s unsupported capacity range, error[%s]. volname[%s], requiredSizeByte[%d]", ERRORLOG, err.Error(), volName, requiredSizeByte)
	}
	klog.Infof("%s Get required creating volume size in bytes[%d], storage-class[%v], topology[%v]", INFOLOG, requiredSizeByte, sc, topo)

	// should not fail when requesting to create a volume with already existing name and same capacity
	// should fail when requesting to create a volume with already existing name and different capacity.
	klog.Infof("%s Will findvolume by name, volname[%s], zone[%s], sizeGB[%d], type[%s]", INFOLOG, volName, topo.GetZone(), common.ByteCeilToGib(requiredSizeByte), sc.GetDiskType().String())
	exVolInfo, err := cs.cloud.FindVolumeByName(volName, topo.GetZone(), common.ByteCeilToGib(requiredSizeByte), sc.GetDiskType().String(), sc.GetPlaceGroupID())
	if err != nil {
		if exVolInfo != nil {
			return nil, status.Errorf(codes.AlreadyExists, "%s volumename[%s] exit, error[%s]", ERRORLOG, volName, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%s find volume by name[%s] error[%s]", ERRORLOG, volName, err.Error())
	}
	if exVolInfo != nil {
		exVolSizeByte := common.GibToByte(exVolInfo.ZecVolume_Size)
		if common.IsValidCapacityBytes(exVolSizeByte, req.GetCapacityRange()) && cs.IsValidTopology(exVolInfo, req.GetAccessibilityRequirements()) && exVolInfo.ZecVolume_Type == sc.GetDiskType().Int() {
			klog.Infof("%s Success findvolume by name, volname[%s], zone[%s], sizeGB[%d], type[%s]", INFOLOG, volName, topo.GetZone(), common.ByteCeilToGib(requiredSizeByte), sc.GetDiskType().String())
			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:           common.GenCsiVolId(exVolInfo.ZecVolume_Id, exVolInfo.ZecVolume_Serial),
					CapacityBytes:      exVolSizeByte,
					VolumeContext:      req.GetParameters(),
					ContentSource:      req.GetVolumeContentSource(),
					AccessibleTopology: cs.GetVolumeTopology(exVolInfo),
				},
			}, nil
		} else {
			// 用于capacitybytes的阈值外逻辑分支
			return nil, status.Errorf(codes.AlreadyExists, "%s volume[%s] already exist but is incompatible.", ERRORLOG, volName)
		}
	}

	volContSrc := req.GetVolumeContentSource()
	if volContSrc == nil {
		// create an empty volume
		requiredSizeGib := common.ByteCeilToGib(requiredSizeByte)
		klog.Infof("%s Will create empty volume[%s], size[%d]", INFOLOG, volName, requiredSizeGib)
		newVolId, err := cs.cloud.CreateVolume(volName, requiredSizeGib, sc.GetDiskType().String(), topo.GetZone(), sc.GetPlaceGroupID())
		if err != nil {
			klog.Errorf("%s Failed to create volume[%s], error[%v]", ERRORLOG, volName, err)
			return nil, status.Error(codes.Internal, err.Error()+volName)
		}

		newVolInfo, err := cs.cloud.FindVolume(newVolId)
		if err != nil {
			klog.Errorf("%s Failed to find volume[%s], error[%v]", ERRORLOG, newVolId, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		if newVolInfo == nil {
			klog.Errorf("%s Cannot find just created volume[%s/%s], please retrying later", ERRORLOG, volName, newVolId)
			return nil, status.Errorf(codes.Aborted, "cannot find volume[%s]", newVolId)
		}

		klog.Infof("%s Succeed create empty volume[%s/%s]", INFOLOG, volName, newVolId)
		//need set tag
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:           common.GenCsiVolId(newVolInfo.ZecVolume_Id, newVolInfo.ZecVolume_Serial),
				CapacityBytes:      requiredSizeByte,
				VolumeContext:      req.GetParameters(),
				AccessibleTopology: cs.GetVolumeTopology(newVolInfo),
			},
		}, nil
	} else {
		if volContSrc.GetSnapshot() != nil {
			//Create vol from exist snapshot

			//check capability
			if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); !isValid {
				return nil, status.Error(codes.Unimplemented, ERRORLOG+" unsupported controller server snapshot capability")
			}
			//get snapshot id
			if len(volContSrc.GetSnapshot().GetSnapshotId()) == 0 {
				return nil, status.Error(codes.InvalidArgument, ERRORLOG+" missing snapshotid")
			}

			snapId := volContSrc.GetSnapshot().GetSnapshotId()

			if acquired := cs.locks.TryAcquire(snapId); !acquired {
				return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, snapId)
			}
			klog.Infof("%s succ lock resource[%s]", INFOLOG, snapId)

			defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, snapId)
			defer cs.locks.Release(snapId)

			//get snapinfo
			klog.Infof("%s will create volume[%s] from snapid[%s] in zone[%s]", INFOLOG, volName, snapId, topo.GetZone())
			snapinfo, err := cs.cloud.FindSnapshot(snapId)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			if snapinfo == nil {
				return nil, status.Errorf(codes.NotFound, "%s cannot find content source snapshotid[%s], disk[%s]", ERRORLOG, snapId, volName)
			}

			//check snapshot create disk ability
			if !snapinfo.ZecVolumeSnap_DiskAbility {
				klog.Errorf("%s snap DiskAbility is false, snapshot not ready. disk[%s], snapid[%s]", ERRORLOG, volName, snapId)
				return nil, status.Errorf(codes.Internal, "%s snap DiskAbility is false, snapshot not ready. disk[%s], snapid[%s]", ERRORLOG, volName, snapId)
			}

			requiredSizeGib := common.ByteCeilToGib(requiredSizeByte)

			//restore vol from snap
			newVolId, err := cs.cloud.CreateVolumeFromSnapshot(volName, requiredSizeGib, sc.GetDiskType().String(), topo.GetZone(), sc.GetPlaceGroupID(), snapId)
			if err != nil {
				klog.Errorf("%s Failed to create volume[%s], snapid[%s], error[%v]", ERRORLOG, volName, snapId, err)
				return nil, status.Error(codes.Internal, err.Error()+volName)
			}

			newVolInfo, err := cs.cloud.FindVolume(newVolId)
			if err != nil {
				klog.Errorf("%s Failed to find volume[%s], error[%v]", ERRORLOG, newVolId, err)
				return nil, status.Error(codes.Internal, err.Error())
			}
			if newVolInfo == nil {
				klog.Errorf("%s Cannot find just created volume[%s/%s], please retrying later", ERRORLOG, volName, newVolId)
				return nil, status.Errorf(codes.Aborted, "cannot find volume[%s]", newVolId)
			}
			klog.Infof("%s Succeed create volume[%s/%s] from snapid[%s]", INFOLOG, volName, newVolId, snapId)

			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:      common.GenCsiVolId(newVolInfo.ZecVolume_Id, newVolInfo.ZecVolume_Serial),
					CapacityBytes: requiredSizeByte,
					VolumeContext: req.GetParameters(),
					ContentSource: &csi.VolumeContentSource{
						Type: &csi.VolumeContentSource_Snapshot{
							Snapshot: &csi.VolumeContentSource_SnapshotSource{
								SnapshotId: snapId,
							},
						},
					},
					AccessibleTopology: cs.GetVolumeTopology(newVolInfo),
				},
			}, nil

		} else if volContSrc.GetVolume() != nil {
			return nil, status.Error(codes.Unimplemented, ERRORLOG+" unsupported controller server clone capability")
		}
	}

	return nil, status.Error(codes.Internal, "Unpredictable error.")
}

/*
action: CSI operation delete zec cloud disk

args: ctx context.Context, req *csi.DeleteVolumeRequest

return: *csi.DeleteVolumeResponse, error
*/
func (cs *ControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	funcName := "ControllerServer:DeleteVolume:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.
		ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); !isValid {
		klog.Errorf("%s invalid delete volume req[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "invalid delete volume req")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume id missing in request")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if acquired := cs.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer cs.locks.Release(volId)

	klog.Infof("%s Will delete volumeid[%s]", INFOLOG, volId)
	if err := cs.cloud.DeleteVolume(volId); err != nil {
		klog.Errorf("%s Failed to delete volumeid[%s], error[%v]", ERRORLOG, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Succeed delete volumeid[%s]", INFOLOG, volId)

	return &csi.DeleteVolumeResponse{}, nil
}

/*
action: CSI operation attach zec cloud disk to VM

args: ctx context.Context, req *csi.ControllerPublishVolumeRequest

return: *csi.ControllerPublishVolumeResponse, error
*/
func (cs *ControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {

	funcName := "ControllerServer:ControllerPublishVolume:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); !isValid {
		klog.Errorf("%s Invalid publish volume req[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "Invalid publish volume req")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume ID missing")
	}

	if len(req.GetNodeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Node ID missing")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"No volume capability")
	}

	if req.GetReadonly() {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport ReadOnly cloud disk")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if acquired := cs.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer cs.locks.Release(volId)

	exVolInfo, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if exVolInfo == nil {
		return nil, status.Errorf(codes.NotFound, "%s Volume: %s does not exist", ERRORLOG, volId)
	}

	vmId := req.GetNodeId()

	vmexist, vmstatus, err := cs.cloud.GetVmStatus(vmId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+"GetVmstatus return err"+err.Error()+vmId)
	}
	if !vmexist {
		return nil, status.Error(codes.NotFound, ERRORLOG+"Vm not exist"+vmId)
	}
	if vmstatus != cloud.VmStatusRunning {
		return nil, status.Error(codes.Internal, ERRORLOG+"Vm status not running"+vmId)
	}

	// Volume published to another node
	if len(exVolInfo.ZecVolume_InstanceId) != 0 {
		if exVolInfo.ZecVolume_InstanceId == vmId {
			klog.Warningf("%s Volumeid[%s] has been already attached on vm[%s]", INFOLOG, volId, vmId)
			return &csi.ControllerPublishVolumeResponse{}, nil
		} else {
			klog.Errorf("%s Volumeid[%s] expected attached on vm[%s], but actually vm[%s]", ERRORLOG, volId, vmId, exVolInfo.ZecVolume_InstanceId)
			return nil, status.Error(codes.FailedPrecondition, "Volume published to another node")
		}
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume capability missing in request")
	}

	klog.Infof("%s Will to Publish volumeid[%s], vmid[%s]", INFOLOG, volId, vmId)
	err = cs.cloud.AttachVolume(volId, vmId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId+vmId)
	}

	newVolInfo, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if newVolInfo == nil {
		return nil, status.Errorf(codes.NotFound, "%s Volume: %s does not exist", ERRORLOG, volId)
	}
	if newVolInfo.ZecVolume_InstanceId != vmId {
		klog.Errorf("%s after attach volume, volume info vmid error, need vmid[%s], volinfo.vmid[%s], will detach vol[%s]", ERRORLOG, vmId, newVolInfo.ZecVolume_InstanceId, volId)
		err = cs.cloud.DetachVolume(volId)
		if err != nil {
			klog.Errorf("%s revert attach action error[%v], volid[%s], need vmid[%s]", ERRORLOG, err.Error(), volId, vmId)
		}
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId+vmId)
	}
	klog.Infof("%s Succeed to Publish volumeid[%s], vmid[%s]", INFOLOG, volId, vmId)

	return &csi.ControllerPublishVolumeResponse{}, nil
}

/*
action: CSI operation detach zec cloud disk from VM

args: ctx context.Context, req *csi.ControllerUnpublishVolumeRequest

return: *csi.ControllerUnpublishVolumeResponse, error
*/
func (cs *ControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	funcName := "ControllerServer:ControllerUnpublishVolume:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.
		ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME); !isValid {
		klog.Errorf("%s Invalid unpublish volume req[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "Invalid unpublish volume req")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume ID missing in request")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	vmId := req.GetNodeId()

	if acquired := cs.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer cs.locks.Release(volId)

	exVol, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if exVol == nil {
		//can not pass csi-sanity, if disk not exist do not return error
		klog.Warningf("%s Volume[%s] is not exist, req vmid[%s]", INFOLOG, volId, vmId)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	} else if exVol.ZecVolume_InstanceId == "" {
		klog.Warningf("%s Volume[%s] is not attached to any instance, req vmid[%s]", INFOLOG, volId, vmId)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	vmexist, vmstatus, err := cs.cloud.GetVmStatus(vmId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+"GetVmstatus return err"+err.Error()+vmId)
	}
	if !vmexist {
		return nil, status.Error(codes.NotFound, ERRORLOG+"Vm not exist"+vmId)
	}
	if vmstatus != cloud.VmStatusRunning {
		return nil, status.Error(codes.Internal, ERRORLOG+"Vm status not running"+vmId)
	}

	// do detach
	if !cs.detachLimiter.Try(volId) {
		return nil, status.Errorf(codes.Internal, "%s volume %s exceeds max retry times %d", ERRORLOG, volId, cs.detachLimiter.GetMaxRetryTimes())
	}

	klog.Infof("%s Will to UnPublish volume[%s], vm[%s]", INFOLOG, volId, vmId)
	err = cs.cloud.DetachVolume(volId)
	if err != nil {
		klog.Errorf("%s Failed to detach volume[%s] from vm[%s] with error[%s]", ERRORLOG, volId, vmId, err.Error())
		cs.detachLimiter.Add(volId)
		return nil, status.Error(codes.Internal, err.Error()+volId)
	}
	klog.Infof("%s Succeed to UnPublish volume[%s], vm[%s]", INFOLOG, volId, vmId)

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

/*
action: CSI operation resize zec cloud disk

args: ctx context.Context, req *csi.ControllerExpandVolumeRequest

return: *csi.ControllerExpandVolumeResponse, error
*/
func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	funcName := "ControllerServer:ControllerExpandVolume:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"No volume id is provided.")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if acquired := cs.locks.TryAcquire(volId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, volId)
	defer cs.locks.Release(volId)

	exVol, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if exVol == nil {
		return nil, status.Errorf(codes.NotFound, "%s Volume[%s] does not exist", ERRORLOG, volId)
	}

	// Get capacity
	voltype := driver.VolumeType(exVol.ZecVolume_Type)
	if !voltype.IsValid() {
		klog.Errorf("%s unsupport voltype[%d], volid[%s]", ERRORLOG, voltype, volId)
		return nil, status.Errorf(codes.Internal, "%s unsupport voltype[%d], volid[%s]", ERRORLOG, voltype, volId)
	}

	sc := driver.NewDefaultZecStorageClassFromType(voltype)
	requiredSizeBytes, err := sc.GetRequiredVolumeSizeByte(req.GetCapacityRange())
	if err != nil {
		return nil, status.Error(codes.OutOfRange, ERRORLOG+err.Error()+volId)
	}

	nodeExpansionRequired := req.GetVolumeCapability().GetBlock() == nil

	exVolSizeBytes := common.GibToByte(exVol.ZecVolume_Size) //disk current size bytes
	if exVolSizeBytes >= requiredSizeBytes {
		klog.Infof("%s: Volume[%s] current size[%d] >= request expand size[%d]", hash, volId, exVolSizeBytes, requiredSizeBytes)

		return &csi.ControllerExpandVolumeResponse{
			CapacityBytes:         exVolSizeBytes,
			NodeExpansionRequired: nodeExpansionRequired,
		}, nil
	}

	klog.Infof("%s Will to Resize volume[%s], ExpandVolume get args requireSize[%d], currentSize[%d]", INFOLOG, volId, requiredSizeBytes, exVolSizeBytes)

	if requiredSizeBytes%common.Gib != 0 {
		return nil, status.Errorf(codes.OutOfRange, "%s required size bytes[%d] cannot be divided into Gib[%d], volId[%s]", ERRORLOG, requiredSizeBytes, common.Gib, volId)
	}

	requiredSizeGib := int(requiredSizeBytes / common.Gib)

	if err = cs.cloud.ResizeVolume(volId, requiredSizeGib); err != nil {
		klog.Errorf("%s Failed to resize volume[%s], error[%v]", ERRORLOG, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Succeed to Resize volume[%s] to size[%d]", INFOLOG, volId, requiredSizeGib)

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         requiredSizeBytes,
		NodeExpansionRequired: nodeExpansionRequired,
	}, nil
}

/*
action: CreateSnapshot allows the CO to create a snapshot.This operation MUST be idempotent.
should fail when requesting to create a snapshot with already existing name and different source volume ID

args: ctx context.Context, req *csi.CreateSnapshotRequest

return: *csi.CreateSnapshotResponse, error
*/
func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	funcName := "ControllerServer:CreateSnapshot:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); !isValid {
		klog.Errorf("%s Invalid create snapshot request[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "")
	}

	if len(req.GetSourceVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"missing volumeID")
	}

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"missing snapshot name")
	}
	if len(req.GetName()) > 64 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"snapshot name is longer then 64")
	}

	srcVolId, _, err := common.ParseCsiVolId(req.GetSourceVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	snapName := req.GetName()

	//lock srcvol
	if acquired := cs.locks.TryAcquire(srcVolId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, srcVolId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, srcVolId)
	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, srcVolId)
	defer cs.locks.Release(srcVolId)

	exVolInfo, err := cs.cloud.FindVolume(srcVolId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+srcVolId)
	}
	if exVolInfo == nil {
		return nil, status.Errorf(codes.NotFound, "%s Volume[%s] does not exist", ERRORLOG, srcVolId)
	}
	if !exVolInfo.ZecVolume_SnapshotAbility {
		return nil, status.Errorf(codes.Internal, "%s Volume[%s] do not support SnapshotAbility", ERRORLOG, srcVolId)
	}

	snapSize := exVolInfo.ZecVolume_Size

	var ready_to_use bool
	klog.Infof("%s Will Find exist Snapshot name[%s], src volid[%s]", INFOLOG, snapName, srcVolId)
	existsnap, err := cs.cloud.FindSnapshotByName(snapName, srcVolId, exVolInfo.ZecVolume_Zone, exVolInfo.ZecVolume_ResourceGroupId)
	if err != nil {
		if existsnap != nil {
			return nil, status.Errorf(codes.AlreadyExists, "%s Find exist snapshot, err %v, name=%s", ERRORLOG, err.Error(), snapName)
		} else {
			return nil, status.Errorf(codes.Internal, "%s Find snap by name return error %v, name=%s", ERRORLOG, err.Error(), snapName)
		}
	}

	if existsnap != nil {
		if existsnap.ZecVolumeSnap_SrcDiskId == srcVolId {
			klog.Infof("%s Success Find exist snapshot name[%s], snapshotid[%s], source volumeid[%s], req source volumeid[%s]", INFOLOG, existsnap.ZecVolumeSnap_Name, existsnap.ZecVolumeSnap_Id, existsnap.ZecVolumeSnap_SrcDiskId, srcVolId)
			if existsnap.ZecVolumeSnap_status == cloud.SnapStatusAvailable {
				ready_to_use = true
			} else {
				ready_to_use = false
			}

			t, err := time.Parse(time.RFC3339, existsnap.ZecVolumeSnap_CreateTime)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "%s Parse Create error, err[%s], snapid[%s], createtime[%s]", ERRORLOG, err.Error(), existsnap.ZecVolumeSnap_Id, existsnap.ZecVolumeSnap_CreateTime)
			}
			ts := timestamppb.New(t)
			klog.Infof("%s Success Find Snapshot name[%s], src volid[%s]", INFOLOG, snapName, srcVolId)
			return &csi.CreateSnapshotResponse{
				Snapshot: &csi.Snapshot{
					SnapshotId:     existsnap.ZecVolumeSnap_Id,
					SourceVolumeId: existsnap.ZecVolumeSnap_SrcDiskId,
					ReadyToUse:     ready_to_use,
					CreationTime:   ts,
					SizeBytes:      int64(snapSize) * common.Gib,
				},
			}, nil
		} else {
			klog.Errorf("%s snapshot name[%s] already exist, but below volume[%s], different to req source volume[%s]", ERRORLOG, snapName, existsnap.ZecVolumeSnap_SrcDiskId, srcVolId)
			return nil, status.Errorf(codes.AlreadyExists, "%s snapshot name[%s] already exist, but below volume[%s], different to req source volume[%s]", ERRORLOG, snapName, existsnap.ZecVolumeSnap_SrcDiskId, srcVolId)
		}
	}

	sc, err := driver.NewZecSnapshotClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	_ = sc

	klog.Infof("%s Will Create a New snapshot name[%s], srcVolId[%s]", INFOLOG, snapName, srcVolId)
	newSnapId, err := cs.cloud.CreateSnapshot(snapName, srcVolId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s create snapshot[%s] from source volume[%s] error[%s]", ERRORLOG, snapName, srcVolId, err.Error())
	}

	snapInfo, err := cs.cloud.FindSnapshot(newSnapId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s Find snapshot[%s] error[%s]", ERRORLOG, newSnapId, err.Error())
	}
	if snapInfo == nil {
		return nil, status.Errorf(codes.Internal, "%s cannot find just created snapshot id[%s]", ERRORLOG, newSnapId)
	}

	klog.Infof("%s Success Create snapshot name[%s], snapshotid[%s], source volumeid[%s], createtime[%s]", INFOLOG, snapInfo.ZecVolumeSnap_Name, snapInfo.ZecVolumeSnap_Id, snapInfo.ZecVolumeSnap_SrcDiskId, snapInfo.ZecVolumeSnap_CreateTime)
	if snapInfo.ZecVolumeSnap_status == cloud.SnapStatusAvailable {
		ready_to_use = true
	} else {
		ready_to_use = false
	}

	// to *timestamppb.Timestamp
	t, err := time.Parse(time.RFC3339, snapInfo.ZecVolumeSnap_CreateTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s Parse Create error, err[%s], snapid[%s], createtime[%s]", ERRORLOG, err.Error(), newSnapId, snapInfo.ZecVolumeSnap_CreateTime)
	}
	ts := timestamppb.New(t)

	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     snapInfo.ZecVolumeSnap_Id,
			SourceVolumeId: snapInfo.ZecVolumeSnap_SrcDiskId,
			ReadyToUse:     ready_to_use,
			CreationTime:   ts,
			SizeBytes:      int64(snapSize) * common.Gib,
		},
	}, nil
}

/*
action: DeleteSnapshot allows the CO to delete a snapshot.
This operation MUST be idempotent.

args: ctx context.Context, req *csi.DeleteSnapshotRequest

return: *csi.DeleteSnapshotResponse, error
*/
func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	funcName := "ControllerServer:DeleteSnapshot:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); !isValid {
		klog.Errorf("%s Invalid delete snapshot request[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "")
	}

	if len(req.GetSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing snapshot id.")
	}

	snapId := req.GetSnapshotId()
	//lock
	if acquired := cs.locks.TryAcquire(snapId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, snapId)
	}
	klog.Infof("%s succ lock resource[%s]", INFOLOG, snapId)
	defer klog.Infof("%s succ unlock resource[%s]", INFOLOG, snapId)
	defer cs.locks.Release(snapId)

	exsnap, err := cs.cloud.FindSnapshot(snapId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exsnap == nil {
		return &csi.DeleteSnapshotResponse{}, nil
	}

	klog.Infof("%s Will to delete snapshot[%s]", INFOLOG, snapId)
	if err = cs.cloud.DeleteSnapshot(snapId); err != nil {
		klog.Errorf("%s Failed to delete snapshot[%s], error[%v]", ERRORLOG, snapId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	klog.Infof("%s Succeed to delete snapshot[%s]", INFOLOG, snapId)
	return &csi.DeleteSnapshotResponse{}, nil
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	funcName := "ControllerServer:ValidateVolumeCapabilities:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"No volume id is provided")
	}

	if len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"No volume capabilities are provided")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	vol, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if vol == nil {
		return nil, status.Errorf(codes.NotFound, "%s volume %s does not exist", ERRORLOG, volId)
	}

	// check capability
	for _, c := range req.GetVolumeCapabilities() {
		found := false
		for _, c1 := range cs.driver.GetVolumeCapability() {
			if c1.GetMode() == c.GetAccessMode().GetMode() {
				found = true
			}
		}
		if !found {
			return &csi.ValidateVolumeCapabilitiesResponse{
				Message: "Driver does not support mode:" + c.GetAccessMode().GetMode().String(),
			}, status.Error(codes.InvalidArgument, ERRORLOG+"Driver does not support mode:"+c.GetAccessMode().GetMode().String())
		}
	}
	_ = INFOLOG
	return &csi.ValidateVolumeCapabilitiesResponse{}, nil
}

func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	funcName := "ControllerServer:ControllerGetCapabilities:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.driver.GetControllerCapability(),
	}, nil
}

func (cs *ControllerServer) PickTopology(requirement *csi.TopologyRequirement) (*driver.Topology, error) {

	topo := &driver.Topology{}
	if requirement == nil {
		return nil, nil
	}

	for _, topology := range requirement.GetPreferred() {
		for k, v := range topology.GetSegments() {
			klog.Infof("INFO:PickTopology() requirement.GetPreferred() k[%s], v[%s]", k, v)
			switch k {
			case cs.driver.GetTopologyZoneKey():
				topo.SetZone(v)
			case cs.driver.GetTopologyVmTypeKey():
				t, ok := driver.VmTypeValue[v]
				if !ok {
					return nil, fmt.Errorf("unsuport instance type[%s]", v)
				}
				topo.SetVmType(t)
			default:
				return nil, fmt.Errorf("invalid topology key[%s]", k)
			}

		}
		return topo, nil
	}

	for _, topology := range requirement.GetRequisite() {
		for k, v := range topology.GetSegments() {
			klog.Infof("INFO:PickTopology() requirement.GetRequisite() k[%s], v[%s]", k, v)
			switch k {
			case cs.driver.GetTopologyZoneKey():
				topo.SetZone(v)
			case cs.driver.GetTopologyVmTypeKey():
				t, ok := driver.VmTypeValue[v]
				if !ok {
					return nil, fmt.Errorf("unsuport instance type[%s]", v)
				}
				topo.SetVmType(t)
			default:
				return nil, fmt.Errorf("invalid topology key[%s]", k)
			}

		}
		return topo, nil
	}

	return nil, nil
}

func (cs *ControllerServer) IsValidTopology(zecVolInfo *cloud.ZecVolume, requirement *csi.TopologyRequirement) bool {
	if zecVolInfo == nil {
		return false
	}
	if requirement == nil || len(requirement.GetRequisite()) == 0 {
		return true
	}
	volTops := cs.GetVolumeTopology(zecVolInfo)
	res := true
	for _, reqTop := range requirement.GetRequisite() {
		for _, volTop := range volTops {
			if reflect.DeepEqual(reqTop, volTop) {
				return true
			} else {
				res = false
			}
		}
	}
	return res
}

func (cs *ControllerServer) GetVolumeTopology(zecVolInfo *cloud.ZecVolume) []*csi.Topology {
	if zecVolInfo == nil {
		return nil
	}
	volType := driver.VolumeType(zecVolInfo.ZecVolume_Type)
	if !volType.IsValid() {
		return nil
	}

	var topo []*csi.Topology

	for _, vmType := range driver.VolumeTypeAttachConstraint[volType] {
		topo = append(topo, &csi.Topology{
			Segments: map[string]string{
				cs.driver.GetTopologyVmTypeKey(): driver.VmTypeName[vmType],
				cs.driver.GetTopologyZoneKey():   zecVolInfo.ZecVolume_Zone,
			},
		})
	}
	return topo
}

func (cs *ControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
