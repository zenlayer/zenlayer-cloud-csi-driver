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

type VolumeManager interface {
	/*
		action: Create a zec cloud disk
	*/
	CreateVolume(volName string, volSize int, volCategory string, zoneId string, resouceGroupID string) (volId string, err error)

	/*
		action: Delete a zec cloud disk
	*/
	DeleteVolume(volId string) (err error)

	/*
		action: Attach a zec cloud disk to Vm
	*/
	AttachVolume(volId string, vmId string) (err error)

	/*
		action: Detach a zec cloud disk from Vm, remove block device
	*/
	DetachVolume(volId string) (err error)

	/*
		action: Resize a zec cloud disk, expand
	*/
	ResizeVolume(volId string, requestSize int) (err error)

	/*
		action: Get disk info by volid
	*/
	FindVolume(volId string) (zecVolInfo *ZecVolume, err error)

	/*
		action: Get disk info by volname
	*/
	FindVolumeByName(volName string, zoneID string, sizeGB int, diskType string) (zecVolInfo *ZecVolume, err error)

	/*
		action: not support
	*/
	CloneVolume() (err error)
}

type UtilManager interface {
	/*
		action: Get Vm info by vmid
	*/
	FindInstance(vmId string) (zecVmInfo *ZecVm, err error)

	/*
		action: Get Zone Info
	*/
	GetZone(zoneId string) error

	/*
		action: List support cloud disk zones
	*/
	GetZoneList() (zoneNameList []string, err error)

	/*
		action: Judgment driver type
	*/
	IsController() bool

	/*
		action: probe check zenlayer cloud connection
	*/
	Probe() error

	/*
		action: Check Vm status by vmid
	*/
	GetVmStatus(vmId string) (exist bool, status string, err error)

	/*
		action: not support
	*/
	FindTag() (err error)

	/*
		action: not support
	*/
	IsValidTags() bool

	/*
		action: not support
	*/
	AttachTags() (err error)
}

type SnapshotManager interface {
	/*
		action: not support
	*/
	FindSnapshot() (err error)

	/*
		action: not support
	*/
	FindSnapshotByName() (err error)

	/*
		action: not support
	*/
	CreateSnapshot() (err error)

	/*
		action: not support
	*/
	DeleteSnapshot() (err error)

	/*
		action: not support
	*/
	CreateVolumeFromSnapshot() (err error)
}

type CloudManager interface {
	VolumeManager
	UtilManager
	SnapshotManager
}
