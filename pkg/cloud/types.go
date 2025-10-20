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

package cloud

/*
Instance Status(v0.2.0)

	DEPLOYING 创建部署中。
	REBUILDING 重装系统中。
	REBOOT 系统重启中。
	RUNNING 运行开机。
	STOPPED 已关机状态。
	BOOTING 系统启动中。
	RELEASING 实例释放中。
	STOPPING 系统关机中。
	RECYCLE 实例处于回收站。
	RECYCLING 实例回收中。
	CREATE_FAILED 创建失败。
	IMAGING 镜像制作中。
	RESIZING 变更规格大小中。
*/
const (
	VmStatusRunning      string = "RUNNING"
	VmStatusStopped      string = "STOPPED"
	VmStatusRecycle      string = "RECYCLE"
	VmStatusStopping     string = "STOPPING"
	VmStatusBooting      string = "BOOTING"
	VmStatusReleasing    string = "RELEASING"
	VmStatusReboot       string = "REBOOT"
	VmStatusDeploying    string = "DEPLOYING"
	VmStatusRebuilding   string = "REBUILDING"
	VmStatusRecycling    string = "RECYCLING"
	VmStatusResizing     string = "RESIZING"
	VmStatusImaging      string = "IMAGING"
	VmStatusCreateFailed string = "CREATE_FAILED"
)

/*
Volume Status(v0.2.0)

	CREATING 			创建中
	IN_USE  			挂载实例的
	AVAILABLE 			可用，未挂载实例的
	CHANGING  			扩容中
	ATTACHING 			挂载中
	DETACHING 			卸载中
	DELETING 			销毁中
	RECYCLING			回收中
	RECYCLED 			处于回收状态等待销毁
	FAILED 				创建失败
	SNAPSHOT_CREATING 	快照创建中
	ROLLING_BACK 		快照回滚中
	DELETED				csi定义状态 已经删除
	STABLE				csi定义状态 稳定状态available或者in_use
*/
const (
	DiskStatusCreating         string = "CREATING"
	DiskStatusInUse            string = "IN_USE"
	DiskStatusAvailable        string = "AVAILABLE"
	DiskStatusChanging         string = "CHANGING"
	DiskStatusAttaching        string = "ATTACHING"
	DiskStatusDetaching        string = "DETACHING"
	DiskStatusDeleting         string = "DELETING"
	DiskStatusRecycling        string = "RECYCLING"
	DiskStatusRecycled         string = "RECYCLED"
	DiskStatusFailed           string = "FAILED"
	DiskStatusSnapshotCreating string = "SNAPSHOT_CREATING"
	DiskStatusRollingBack      string = "ROLLING_BACK"
	DiskStatusDeleted          string = "DELETED"
	DiskStatusStable           string = "STABLE"
)

/*
Snapshot Status(v0.2.0)

	CREATING		创建中
	AVAILABLE		创建成功
	FAILED			创建失败
	ROLLING_BACK	回滚到此快照中
	DELETING		释放中
	DELETED			csi定义状态 已经删除
*/
const (
	SnapStatusCreating    string = "CREATING"
	SnapStatusAvailable   string = "AVAILABLE"
	SnapStatusFailed      string = "FAILED"
	SnapStatusRollingBack string = "ROLLING_BACK"
	SnapStatusDeleting    string = "DELETING"
	SnapStatusDeleted     string = "DELETED"
)

/*disk info*/
type ZecVolume struct {
	ZecVolume_Id              string
	ZecVolume_Name            string
	ZecVolume_Size            int
	ZecVolume_Type            int    //basic or standard
	ZecVolume_Status          string //DiskStatusAvaliable ...
	ZecVolume_Zone            string //zone
	ZecVolume_Serial          string //disk serial
	ZecVolume_InstanceId      string //has been attached vmid
	ZecVolume_Portable        bool   //support detach
	ZecVolume_SnapshotAbility bool   //support create snap
	ZecVolume_ResourceGroupId string
}

/*vm info*/
type ZecVm struct {
	ZecVm_Id     string
	ZecVm_Status string
	ZecVm_Type   int
	ZecVm_Zone   string
}

/*snap info*/
type ZecVolumeSnap struct {
	ZecVolumeSnap_Id              string
	ZecVolumeSnap_Name            string
	ZecVolumeSnap_status          string //快照状态
	ZecVolumeSnap_SrcDiskId       string //快照属于哪个disk
	ZecVolumeSnap_DiskAbility     bool   //是否具有创建disk的能力
	ZecVolumeSnap_ZoneId          string
	ZecVolumeSnap_ResourceGroupId string
	ZecVolumeSnap_CreateTime      string
}

func (zv ZecVolume) Check() bool {
	return true
}

func NewZecVolume() *ZecVolume {
	zecVolInfo := &ZecVolume{
		ZecVolume_Id:              "",
		ZecVolume_Name:            "",
		ZecVolume_Size:            -1,
		ZecVolume_Type:            -1,
		ZecVolume_Status:          "",
		ZecVolume_Zone:            "",
		ZecVolume_Serial:          "",
		ZecVolume_InstanceId:      "",
		ZecVolume_Portable:        false,
		ZecVolume_SnapshotAbility: false,
		ZecVolume_ResourceGroupId: "",
	}
	return zecVolInfo
}

func NewZecVm() *ZecVm {
	zecVmInfo := &ZecVm{
		ZecVm_Id:     "",
		ZecVm_Status: "",
		ZecVm_Type:   -1,
		ZecVm_Zone:   "",
	}
	return zecVmInfo
}

func NewZecVolumeSnap() *ZecVolumeSnap {
	zecVolSnapInfo := &ZecVolumeSnap{
		ZecVolumeSnap_Id:              "",
		ZecVolumeSnap_Name:            "",
		ZecVolumeSnap_status:          "",
		ZecVolumeSnap_SrcDiskId:       "",
		ZecVolumeSnap_DiskAbility:     false,
		ZecVolumeSnap_ZoneId:          "",
		ZecVolumeSnap_ResourceGroupId: "",
		ZecVolumeSnap_CreateTime:      "",
	}
	return zecVolSnapInfo
}
