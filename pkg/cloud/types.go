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
Instance Status(v0.1.27)

	RUNNING			开机
	STOPPED			关机
	RECYCLE			已回收
	STOPPING		关机中
	BOOTING			开机中
	RELEASING		释放中
	REBOOT			重启中
	DEPLOYING		部署中
	REBUILDING		重装系统中
	RECYCLING		回收中
	RESIZING		变更规格中
	IMAGING 		制作镜像中
	CREATE_FAILED 	创建（重装）失败
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
Volume Status(v0.1.27)

	CREATING 	创建中
	IN_USE  	挂载实例的
	AVAILABLE 	可用，未挂载实例的
	CHANGING  	扩所容中
	ATTACHING 	挂载中
	DETACHING 	卸载中
	DELETING 	删除中
	RECYCLED 	处于回收状态
	FAILED 		创建失败
	DELETED		csi状态 已经删除
	STABLE		csi状态 稳定状态available或者in_use
*/
const (
	DiskStatusCreating  string = "CREATING"
	DiskStatusInUse     string = "IN_USE"
	DiskStatusAvaliable string = "AVAILABLE"
	DiskStatusChanging  string = "CHANGING"
	DiskStatusAttaching string = "ATTACHING"
	DiskStatusDetaching string = "DETACHING"
	DiskStatusDeleting  string = "DELETING"
	DiskStatusRecycled  string = "RECYCLED"
	DiskStatusFailed    string = "FAILED"
	DiskStatusDeleted   string = "DELETED"
	DiskStatusStable    string = "STABLE"
)

type ZecVolume struct {
	ZecVolume_Id         string
	ZecVolume_Name       string
	ZecVolume_Size       int
	ZecVolume_Type       int    //basic or standard
	ZecVolume_Status     string //DiskStatusAvaliable ...
	ZecVolume_Zone       string //zone
	ZecVolume_Serial     string //disk serial
	ZecVolume_InstanceId string //has been attached vmid
}

type ZecVm struct {
	ZecVm_Id     string
	ZecVm_Status string
	ZecVm_Type   int
	ZecVm_Zone   string
}

func (zv ZecVolume) Check() bool {
	return true
}

func NewZecVolume() (zecVolInfo *ZecVolume) {
	zecVolInfo = &ZecVolume{
		ZecVolume_Id:         "",
		ZecVolume_Name:       "",
		ZecVolume_Size:       -1,
		ZecVolume_Type:       -1,
		ZecVolume_Zone:       "",
		ZecVolume_Serial:     "",
		ZecVolume_InstanceId: "",
	}
	return
}

func NewZecVm() (zecVmInfo *ZecVm) {
	zecVmInfo = &ZecVm{
		ZecVm_Id:     "",
		ZecVm_Status: "",
	}
	return
}
