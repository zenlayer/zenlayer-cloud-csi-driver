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
	"strconv"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/klog"
)

type ControllerServer struct {
	csi.UnimplementedControllerServer

	driver *driver.DiskDriver
	cloud  cloud.CloudManager
	locks  *common.ResourceLocks
}

func NewControllerServer(d *driver.DiskDriver, c cloud.CloudManager) *ControllerServer {
	return &ControllerServer{
		driver: d,
		cloud:  c,
		locks:  common.NewResourceLocks(),
	}
}

var _ csi.ControllerServer = &ControllerServer{}

func genCsiVolId(volInfo *cloud.ZecVolume) (string, error) {
	if volInfo == nil {
		return "", fmt.Errorf("nil volume info")
	}
	if volInfo.ZecVolume_Id == "" || volInfo.ZecVolume_Serial == "" {
		return "", fmt.Errorf("volume name[%s] id[%s] serial[%s]: id or serial is empty, cannot build volume id",
			volInfo.ZecVolume_Name, volInfo.ZecVolume_Id, volInfo.ZecVolume_Serial)
	}
	return common.GenCsiVolId(volInfo.ZecVolume_Id, volInfo.ZecVolume_Serial), nil
}

// snapCreateTimeLayouts 列出 DescribeSnapshots 的 createTime 可能采用的时间格式。
// SDK 只把这个字段描述成"创建时间", 没有承诺具体格式, 所以按从严到宽依次尝试。
var snapCreateTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02",
}

/*
parseSnapCreationTime 把云端返回的快照创建时间转成 CSI 需要的 timestamp。

	这里刻意不返回 error: 走到这个函数时快照在云上已经创建成功了, 若因为解析不了一个
	时间戳就让 CreateSnapshot 返回 Internal, 幂等重试会在同一处反复失败 —— 快照永远
	不会 ReadyToUse, 云上资源也就永久泄漏。所以全部格式都不匹配时退化成"当前时间"并打
	一条 error 日志: creation_time 在规范里只是描述信息, 用它换取整条链路能继续推进。
*/
func parseSnapCreationTime(raw string, snapId string, errorLog string) *timestamppb.Timestamp {
	trimmed := strings.TrimSpace(raw)
	if trimmed != "" {
		for _, layout := range snapCreateTimeLayouts {
			if t, err := time.Parse(layout, trimmed); err == nil {
				return timestamppb.New(t)
			}
		}
		// 有些接口把时间返回成 epoch 数字字符串(秒或毫秒)
		if epoch, err := strconv.ParseInt(trimmed, 10, 64); err == nil && epoch > 0 {
			if len(trimmed) >= 13 {
				return timestamppb.New(time.UnixMilli(epoch))
			}
			return timestamppb.New(time.Unix(epoch, 0))
		}
	}

	klog.Errorf("%s cannot parse snapshot createtime[%s], snapid[%s], fall back to current time", errorLog, raw, snapId)
	return timestamppb.Now()
}

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
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, ERRORLOG+err.Error())
		}
		if topo == nil {
			return nil, status.Errorf(codes.InvalidArgument, "%s cannot pick topology from accessibility requirements, volname[%s]", ERRORLOG, volName)
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
			csiVolId, err := genCsiVolId(exVolInfo)
			if err != nil {
				klog.Errorf("%s %v, please retry later", ERRORLOG, err)
				return nil, status.Errorf(codes.Aborted, "%s %v, please retry later", ERRORLOG, err)
			}
			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:           csiVolId,
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
		newVolId, err := cs.cloud.CreateVolume(volName, requiredSizeGib, sc.GetDiskType().String(), topo.GetZone(), sc.GetPlaceGroupID(), sc.GetBurstEnable())
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

		csiVolId, err := genCsiVolId(newVolInfo)
		if err != nil {
			klog.Errorf("%s %v, please retry later", ERRORLOG, err)
			return nil, status.Errorf(codes.Aborted, "%s %v, please retry later", ERRORLOG, err)
		}

		klog.Infof("%s Succeed create empty volume[%s/%s]", INFOLOG, volName, newVolId)
		//need set tag
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:           csiVolId,
				CapacityBytes:      common.GibToByte(requiredSizeGib),
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
			newVolId, err := cs.cloud.CreateVolumeFromSnapshot(volName, requiredSizeGib, sc.GetDiskType().String(), topo.GetZone(), sc.GetPlaceGroupID(), snapId, sc.GetBurstEnable())
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
			csiVolId, err := genCsiVolId(newVolInfo)
			if err != nil {
				klog.Errorf("%s %v, please retry later", ERRORLOG, err)
				return nil, status.Errorf(codes.Aborted, "%s %v, please retry later", ERRORLOG, err)
			}

			klog.Infof("%s Succeed create volume[%s/%s] from snapid[%s]", INFOLOG, volName, newVolId, snapId)

			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:      csiVolId,
					CapacityBytes: common.GibToByte(requiredSizeGib),
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
		klog.Warningf("%s %v, treat as already deleted", INFOLOG, err)
		return &csi.DeleteVolumeResponse{}, nil
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

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME); !isValid {
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
	// 只判非 nil 是不够的: 一块 ZEC 云盘同一时刻只能挂在一台实例上, 所以 MULTI_NODE_*
	// 之类本驱动没有上报的访问模式必须在这里就拒掉, 否则会真的把卷 attach 上去, 由 CO
	// 以为自己拿到了多节点共享语义。按规范这是参数非法, 回 InvalidArgument, 与
	// ValidateVolumeCapabilities / NodePublishVolume 共用同一套访问模式判断。
	if !cs.driver.ValidateVolumeCapability(req.GetVolumeCapability()) {
		klog.Errorf("%s unsupported volume capability access mode[%s], volumeid[%s]", ERRORLOG, req.GetVolumeCapability().GetAccessMode().GetMode().String(), req.GetVolumeId())
		return nil, status.Errorf(codes.InvalidArgument, "%s unsupported volume capability access mode[%s]", ERRORLOG, req.GetVolumeCapability().GetAccessMode().GetMode().String())
	}

	if req.GetReadonly() {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+" unsupport ReadOnly cloud disk")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		// id 格式对 CO 合法但不是本驱动发出的, 等价于"该卷在本驱动不存在", 按规范回
		// NotFound(csi-sanity: ControllerPublishVolume "should fail when the volume
		// does not exist" 断言的就是 NotFound)。
		return nil, status.Errorf(codes.NotFound, "%s %v", ERRORLOG, err)
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
		return nil, status.Errorf(codes.NotFound, "%s node[%s] does not exist", ERRORLOG, vmId)
	}
	// 实例不在 RUNNING 状态不是驱动内部错误, 而是需要外部先把实例拉起来的前置条件,
	// 所以用 FailedPrecondition 而不是 Internal —— 后者会让运维误以为是驱动故障。
	if vmstatus != cloud.VmStatusRunning {
		return nil, status.Errorf(codes.FailedPrecondition, "%s node[%s] status is [%s], attach requires status [%s]", ERRORLOG, vmId, vmstatus, cloud.VmStatusRunning)
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
		klog.Errorf("%s after attach volume, volume is attached to an unexpected vm, need vmid[%s], volinfo.vmid[%s], vol[%s]. will NOT detach: the volume may be in use by that vm", ERRORLOG, vmId, newVolInfo.ZecVolume_InstanceId, volId)
		return nil, status.Errorf(codes.Internal, "%s after attach, volume[%s] not attached to expected vm[%s], actual vm[%s]", ERRORLOG, volId, vmId, newVolInfo.ZecVolume_InstanceId)
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

	// node_id 在 CSI 规范中是必填字段
	if len(req.GetNodeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"Node ID missing in request")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		// 与 DeleteVolume 同理, unpublish 必须幂等: 解析不出来的 volume id 不对应任何
		// 云盘, 也就不存在需要解除的挂载关系, 返回成功而不是报错阻塞 Pod 重新调度。
		klog.Warningf("%s %v, treat detach as done (idempotent)", INFOLOG, err)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
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

	// 卷已挂载在某实例上。若其实际挂载的实例与请求卸载的 node 不一致，
	// 说明卷并未挂在该 node 上，按 CSI 幂等语义直接返回成功（含 node 已被删除的场景）。
	// 这里不再查询 VM 是否存在/是否 running：
	//   - 底层 DetachVolume 使用 InstanceCheckFlag=false，允许对关机/非 running 实例卸载；
	//   - node 不存在时按幂等应返回成功，而非 NotFound，否则会阻塞卷卸载与 Pod 重新调度。
	if exVol.ZecVolume_InstanceId != vmId {
		klog.Warningf("%s Volume[%s] attached to instance[%s], not the req node[%s], treat detach as done (idempotent)", INFOLOG, volId, exVol.ZecVolume_InstanceId, vmId)
		return &csi.ControllerUnpublishVolumeResponse{}, nil
	}

	// do detach
	klog.Infof("%s Will to UnPublish volume[%s], vm[%s]", INFOLOG, volId, vmId)
	err = cs.cloud.DetachVolume(volId)
	if err != nil {
		klog.Errorf("%s Failed to detach volume[%s] from vm[%s] with error[%s]", ERRORLOG, volId, vmId, err.Error())
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

	if isValid := cs.driver.ValidateControllerServiceRequest(csi.
		ControllerServiceCapability_RPC_EXPAND_VOLUME); !isValid {
		klog.Errorf("%s Invalid expand volume req[%v]", ERRORLOG, req)
		return nil, status.Error(codes.Unimplemented, "Invalid expand volume req")
	}

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ERRORLOG+"No volume id is provided.")
	}

	volId, _, err := common.ParseCsiVolId(req.GetVolumeId())
	if err != nil {
		//该卷在本驱动不存在, 按规范回 NotFound
		return nil, status.Errorf(codes.NotFound, "%s %v", ERRORLOG, err)
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

	requiredSizeGib := common.ByteCeilToGib(requiredSizeBytes)
	actualSizeBytes := common.GibToByte(requiredSizeGib)

	exVolSizeBytes := common.GibToByte(exVol.ZecVolume_Size) //disk current size bytes
	if exVolSizeBytes >= actualSizeBytes {
		klog.Infof("%s: Volume[%s] current size[%d] >= request expand size[%d]", hash, volId, exVolSizeBytes, actualSizeBytes)

		return &csi.ControllerExpandVolumeResponse{
			CapacityBytes:         exVolSizeBytes,
			NodeExpansionRequired: nodeExpansionRequired,
		}, nil
	}

	klog.Infof("%s Will to Resize volume[%s], ExpandVolume get args requireSize[%d], currentSize[%d]", INFOLOG, volId, actualSizeBytes, exVolSizeBytes)

	if err = cs.cloud.ResizeVolume(volId, requiredSizeGib); err != nil {
		klog.Errorf("%s Failed to resize volume[%s], error[%v]", ERRORLOG, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("%s Succeed to Resize volume[%s] to size[%d]", INFOLOG, volId, requiredSizeGib)

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         actualSizeBytes,
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
		//源卷在本驱动不存在, 按规范回 NotFound(csi-sanity: CreateSnapshot
		//"should fail when the volume does not exist" 断言的就是 NotFound)
		return nil, status.Errorf(codes.NotFound, "%s %v", ERRORLOG, err)
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

			ts := parseSnapCreationTime(existsnap.ZecVolumeSnap_CreateTime, existsnap.ZecVolumeSnap_Id, ERRORLOG)
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

	// Parse (and validate) the snapshot class parameters. The resulting config
	// is not consumed yet, so the value is intentionally discarded.
	if _, err = driver.NewZecSnapshotClassFromMap(req.GetParameters()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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
	ts := parseSnapCreationTime(snapInfo.ZecVolumeSnap_CreateTime, newSnapId, ERRORLOG)

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

/*
action: ValidateVolumeCapabilities tells the CO whether the given volume can be
used with the requested capabilities.

	Per the CSI spec an unsupported capability is NOT an RPC error: the plugin
	must answer successfully with `confirmed` unset and an optional `message`.
	Only a successful validation may set `confirmed`, and leaving it unset on
	success would make the CO treat the volume as incompatible.

args: ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest

return: *csi.ValidateVolumeCapabilitiesResponse, error
*/
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
		//该卷在本驱动不存在, 按规范回 NotFound(csi-sanity: ValidateVolumeCapabilities
		//"should fail when the requested volume does not exist" 断言的就是 NotFound)
		return nil, status.Errorf(codes.NotFound, "%s %v", ERRORLOG, err)
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
			//不支持的访问模式不是rpc错误,按csi规范返回成功响应且不带confirmed
			klog.Infof("%s volume[%s] does not support mode[%s]", INFOLOG, volId, c.GetAccessMode().GetMode().String())
			return &csi.ValidateVolumeCapabilitiesResponse{
				Message: ERRORLOG + "Driver does not support mode:" + c.GetAccessMode().GetMode().String(),
			}, nil
		}
	}

	//校验通过必须回填confirmed,否则CO会认为该卷不满足要求
	klog.Infof("%s volume[%s] supports all requested capabilities", INFOLOG, volId)
	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeContext:      req.GetVolumeContext(),
			VolumeCapabilities: req.GetVolumeCapabilities(),
			Parameters:         req.GetParameters(),
		},
	}, nil
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
	// The existing volume is compatible when its topology matches any one of the
	// requisite topologies. If nothing matches (including an empty volTops), the
	// volume is not in an acceptable topology.
	for _, reqTop := range requirement.GetRequisite() {
		for _, volTop := range volTops {
			// 用 proto.Equal 做语义比较，而非 reflect.DeepEqual：后者会比较 protobuf
			// 消息里 Unmarshal 后才填充的内部簿记字段（sizeCache/state 等），导致
			// 线上收到的 reqTop 与现场构造的 volTop 即使 Segments 相同也判不等，
			// 进而破坏 CreateVolume 的幂等重试。
			if proto.Equal(reqTop, volTop) {
				return true
			}
		}
	}
	return false
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
