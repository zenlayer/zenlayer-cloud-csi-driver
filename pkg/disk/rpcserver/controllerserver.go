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

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	volName := req.GetName()

	if acquired := cs.locks.TryAcquire(volName); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volName)
	}
	klog.Infof("%s succ lock resource %s", INFOLOG, volName)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volName)
	defer cs.locks.Release(volName)

	//read conf and init storage-class
	sc, err := driver.NewZecStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+err.Error())
	}

	topo := &driver.Topology{}
	if req.GetAccessibilityRequirements() != nil && cs.driver.ValidatePluginCapabilityService(csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS) {
		klog.Info(INFOLOG + " GetAccessibilityRequirements has val. volName=" + volName)
		var err error
		topo, err = cs.PickTopology(req.GetAccessibilityRequirements())
		if err != nil || topo == nil {
			return nil, status.Error(codes.InvalidArgument, ERRORLOG+err.Error())
		}
		//这里需要将req.Parameters中的zoneID修改成topo.ZoneID,req是从storageclass.yaml中取的，如果是WaitForFirstConsumer模式走进这个分支storageclass中定义的zone不一定是这个vol选择的zone
		driver.UpdateParmsZone(req.GetParameters(), topo.GetZone())
	} else {
		klog.Info(INFOLOG + "GetAccessibilityRequirements is nil, use storage-class config zone " + sc.GetZone() + " volname=" + volName)
		//only support one Vm type (BasicVm)
		topo = driver.NewTopology(sc.GetZone(), driver.BasicVmType)
	}

	// get request volume capacity range
	requiredSizeByte, err := sc.GetRequiredVolumeSizeByte(req.GetCapacityRange())
	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "%s unsupported capacity range, error: %s. volname=%s.", ERRORLOG, err.Error(), volName)
	}
	klog.Infof("%s: Get required creating volume size in bytes %d, storage-class %v, topology %v", INFOLOG, requiredSizeByte, sc, topo)

	// should not fail when requesting to create a volume with already existing name and same capacity
	// should fail when requesting to create a volume with already existing name and different capacity.
	klog.Infof("%s: Will findvolume by name, volname=%s, zone=%s, sizeGB=%d, type=%s", INFOLOG, volName, topo.GetZone(), common.ByteCeilToGib(requiredSizeByte), sc.GetDiskType().String())
	exVolInfo, err := cs.cloud.FindVolumeByName(volName, topo.GetZone(), common.ByteCeilToGib(requiredSizeByte), sc.GetDiskType().String())
	if err != nil {
		if exVolInfo != nil {
			return nil, status.Errorf(codes.AlreadyExists, "%s volumename exit %s, error=%s", ERRORLOG, volName, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%s find volume by name error: %s, %s", ERRORLOG, volName, err.Error())
	}
	if exVolInfo != nil {
		exVolSizeByte := common.GibToByte(exVolInfo.ZecVolume_Size)
		if common.IsValidCapacityBytes(exVolSizeByte, req.GetCapacityRange()) && cs.IsValidTopology(exVolInfo, req.GetAccessibilityRequirements()) && exVolInfo.ZecVolume_Type == sc.GetDiskType().Int() {
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
			return nil, status.Errorf(codes.AlreadyExists, "%s volume %s already exist but is incompatible.", ERRORLOG, volName)
		}
	}

	volContSrc := req.GetVolumeContentSource()
	if volContSrc == nil {
		// create an empty volume
		requiredSizeGib := common.ByteCeilToGib(requiredSizeByte)
		newVolId, err := cs.cloud.CreateVolume(volName, requiredSizeGib, sc.GetDiskType().String(), topo.GetZone(), sc.GetPlaceGroupID())
		if err != nil {
			klog.Errorf("%s Failed to create volume %s, error: %v", ERRORLOG, volName, err)
			return nil, status.Error(codes.Internal, err.Error()+volName)
		}

		newVolInfo, err := cs.cloud.FindVolume(newVolId)
		if err != nil {
			klog.Errorf("%s Failed to find volume %s, error: %v", ERRORLOG, newVolId, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		if newVolInfo == nil {
			klog.Infof("%s Cannot find just created volume [%s/%s], please retrying later.", ERRORLOG, volName, newVolId)
			return nil, status.Errorf(codes.Aborted, "cannot find volume %s", newVolId)
		}

		klog.Infof("%s Succeed to create empty volume [%s/%s].", INFOLOG, volName, newVolId)
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
			return nil, status.Error(codes.Internal, "Not support snap.")
		} else if volContSrc.GetVolume() != nil {
			return nil, status.Error(codes.Internal, "Not support clone.")
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
		klog.Errorf("%s invalid delete volume req: %v", ERRORLOG, req)
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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer cs.locks.Release(volId)

	if err := cs.cloud.DeleteVolume(volId); err != nil {
		klog.Errorf("%s Failed to delete volume %s, error: %v", ERRORLOG, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Succeed to delete volume [%s].", INFOLOG, volId)

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

	if req.GetVolumeCapability().GetBlock() != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport Block Mode")
	}

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); !isValid {
		klog.Errorf("%s: Invalid publish volume req: %v", ERRORLOG, req)
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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
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
			klog.Warningf("%s: Volume %s has been already attached on instance %s", INFOLOG, volId, vmId)
			return &csi.ControllerPublishVolumeResponse{}, nil
		} else {
			klog.Errorf("%s: Volume %s expected attached on instance %s, but actually %s.", ERRORLOG, volId, vmId, exVolInfo.ZecVolume_InstanceId)
			return nil, status.Error(codes.FailedPrecondition, "Volume published to another node")
		}
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Volume capability missing in request")
	}

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
		klog.Errorf("%s after attach volume, volume info vmid error, need vmid=%s, volinfo.vmid:%s, will detach vol:%s", ERRORLOG, vmId, newVolInfo.ZecVolume_InstanceId, volId)
		err = cs.cloud.DetachVolume(volId)
		if err != nil {
			klog.Errorf("%s revert attach action error:%v, volid=%s, need vmid=%s", ERRORLOG, err.Error(), volId, vmId)
		}
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId+vmId)
	}
	klog.Infof("%s Succeed to Publish volume [%s], vm [%s].", INFOLOG, volId, vmId)

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
		klog.Errorf("%s Invalid unpublish volume req: %v", ERRORLOG, req)
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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer cs.locks.Release(volId)

	exVol, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if exVol == nil {
		//can not pass csi-sanity, if disk not exist do not return error
		klog.Warningf("%s Volume %s is not exist, req vmid %s", INFOLOG, volId, vmId)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	} else if exVol.ZecVolume_InstanceId == "" {
		klog.Warningf("%s Volume %s is not attached to any instance, req vmid %s", INFOLOG, volId, vmId)
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
	klog.Infof("%s Volume id %s retry times is %d", volId, INFOLOG, cs.detachLimiter.GetCurrentRetryTimes(volId))
	if !cs.detachLimiter.Try(volId) {
		return nil, status.Errorf(codes.Internal, "%s volume %s exceeds max retry times %d.", ERRORLOG, volId, cs.detachLimiter.GetMaxRetryTimes())
	}

	err = cs.cloud.DetachVolume(volId)
	if err != nil {
		klog.Errorf("%s Failed to detach disk image: %s from instance %s with error: %s",
			ERRORLOG, volId, vmId, err.Error())
		cs.detachLimiter.Add(volId)
		return nil, status.Error(codes.Internal, err.Error()+volId)
	}
	klog.Infof("%s Succeed to UnPublish volume [%s], vm [%s].", INFOLOG, volId, vmId)

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

	if req.GetVolumeCapability().GetBlock() != nil {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport Block Mode")
	}

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
	klog.Infof("%s succ lock resource %s", INFOLOG, volId)

	defer klog.Infof("%s succ unlock resource %s", INFOLOG, volId)
	defer cs.locks.Release(volId)

	exVol, err := cs.cloud.FindVolume(volId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERRORLOG+err.Error()+volId)
	}
	if exVol == nil {
		return nil, status.Errorf(codes.NotFound, "%s Volume: %s does not exist", ERRORLOG, volId)
	}

	// Get capacity
	voltype := driver.VolumeType(exVol.ZecVolume_Type)
	if !voltype.IsValid() {
		klog.Errorf("%s unsupport voltype %d, volid=%s", ERRORLOG, voltype, volId)
		return nil, status.Errorf(codes.Internal, "%s unsupport voltype %d, volid=%s", ERRORLOG, voltype, volId)
	}

	sc := driver.NewDefaultZecStorageClassFromType(voltype)
	requiredSizeBytes, err := sc.GetRequiredVolumeSizeByte(req.GetCapacityRange())
	if err != nil {
		return nil, status.Error(codes.OutOfRange, ERRORLOG+err.Error()+volId)
	}

	nodeExpansionRequired := true

	exVolSizeBytes := common.GibToByte(exVol.ZecVolume_Size) //disk current size bytes
	if exVolSizeBytes >= requiredSizeBytes {
		klog.Infof("%s: Volume current %s size %d >= request expand size %d", hash, volId, exVolSizeBytes, requiredSizeBytes)

		return &csi.ControllerExpandVolumeResponse{
			CapacityBytes:         exVolSizeBytes,
			NodeExpansionRequired: nodeExpansionRequired,
		}, nil
	}

	klog.Infof("%s ExpandVolume get args requireSize=%d, currentSize=%d, volId=%s.", INFOLOG, requiredSizeBytes, exVolSizeBytes, volId)

	if requiredSizeBytes%common.Gib != 0 {
		return nil, status.Errorf(codes.OutOfRange, "%s required size bytes %d cannot be divided into Gib %d, volId=%s.",
			ERRORLOG, requiredSizeBytes, common.Gib, volId)
	}

	requiredSizeGib := int(requiredSizeBytes / common.Gib)

	if err = cs.cloud.ResizeVolume(volId, requiredSizeGib); err != nil {
		klog.Errorf("%s Failed to resize volume %s, error: %v", ERRORLOG, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Succeed to Resize volume [%s].", INFOLOG, volId)

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         requiredSizeBytes,
		NodeExpansionRequired: nodeExpansionRequired,
	}, nil
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.
	ValidateVolumeCapabilitiesResponse, error) {
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
	klog.Info("ControllerGetCapabilities")
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
			klog.Infof("INFO:PickTopology() requirement.GetPreferred() k=%s, v=%s", k, v)
			switch k {
			case cs.driver.GetTopologyZoneKey():
				topo.SetZone(v)
			case cs.driver.GetTopologyVmTypeKey():
				t, ok := driver.VmTypeValue[v]
				if !ok {
					return nil, fmt.Errorf("unsuport instance type %s", v)
				}
				topo.SetVmType(t)
			default:
				return nil, fmt.Errorf("invalid topology key %s", k)
			}

		}
		return topo, nil
	}

	for _, topology := range requirement.GetRequisite() {
		for k, v := range topology.GetSegments() {
			klog.Infof("INFO:PickTopology() requirement.GetRequisite() k=%s, v=%s", k, v)
			switch k {
			case cs.driver.GetTopologyZoneKey():
				topo.SetZone(v)
			case cs.driver.GetTopologyVmTypeKey():
				t, ok := driver.VmTypeValue[v]
				if !ok {
					return nil, fmt.Errorf("unsuport instance type %s", v)
				}
				topo.SetVmType(t)
			default:
				return nil, fmt.Errorf("invalid topology key %s", k)
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

func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{},
	}, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return &csi.DeleteSnapshotResponse{}, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
