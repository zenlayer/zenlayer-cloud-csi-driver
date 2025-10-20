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

package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
)

/*
VolumeCapability_AccessMode_UNKNOWN 							VolumeCapability_AccessMode_Mode = 0
VolumeCapability_AccessMode_SINGLE_NODE_WRITER 					VolumeCapability_AccessMode_Mode = 1
VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY 			VolumeCapability_AccessMode_Mode = 2
VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY 				VolumeCapability_AccessMode_Mode = 3
VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER 			VolumeCapability_AccessMode_Mode = 4
VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER 			VolumeCapability_AccessMode_Mode = 5
VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER 			VolumeCapability_AccessMode_Mode = 6
VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER 			VolumeCapability_AccessMode_Mode = 7
*/
var DefaultVolumeAccessModeType = []csi.VolumeCapability_AccessMode_Mode{
	csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
}

/*
	 controller type support capability/feature:

		ControllerServiceCapability_RPC_UNKNOWN                  		ControllerServiceCapability_RPC_Type = 0
		ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME     		ControllerServiceCapability_RPC_Type = 1
		ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME 		ControllerServiceCapability_RPC_Type = 2
		ControllerServiceCapability_RPC_LIST_VOLUMES             		ControllerServiceCapability_RPC_Type = 3
		ControllerServiceCapability_RPC_GET_CAPACITY             		ControllerServiceCapability_RPC_Type = 4
		ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT 			ControllerServiceCapability_RPC_Type = 5
		ControllerServiceCapability_RPC_LIST_SNAPSHOTS         			ControllerServiceCapability_RPC_Type = 6
		ControllerServiceCapability_RPC_CLONE_VOLUME 					ControllerServiceCapability_RPC_Type = 7
		ControllerServiceCapability_RPC_PUBLISH_READONLY 				ControllerServiceCapability_RPC_Type = 8
		ControllerServiceCapability_RPC_EXPAND_VOLUME 					ControllerServiceCapability_RPC_Type = 9
		ControllerServiceCapability_RPC_LIST_VOLUMES_PUBLISHED_NODES 	ControllerServiceCapability_RPC_Type = 10
		ControllerServiceCapability_RPC_VOLUME_CONDITION 				ControllerServiceCapability_RPC_Type = 11
		ControllerServiceCapability_RPC_GET_VOLUME 						ControllerServiceCapability_RPC_Type = 12
		ControllerServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER 		ControllerServiceCapability_RPC_Type = 13
*/
var DefaultControllerServiceCapability = []csi.ControllerServiceCapability_RPC_Type{
	csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
	csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
}

/*
	 node type support capability/feature:

		NodeServiceCapability_RPC_UNKNOWN              					NodeServiceCapability_RPC_Type = 0
		NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME 					NodeServiceCapability_RPC_Type = 1
		NodeServiceCapability_RPC_GET_VOLUME_STATS 						NodeServiceCapability_RPC_Type = 2
		NodeServiceCapability_RPC_EXPAND_VOLUME 						NodeServiceCapability_RPC_Type = 3
		NodeServiceCapability_RPC_VOLUME_CONDITION 						NodeServiceCapability_RPC_Type = 4
		NodeServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER 				NodeServiceCapability_RPC_Type = 5
		NodeServiceCapability_RPC_VOLUME_MOUNT_GROUP 					NodeServiceCapability_RPC_Type = 6
*/
var DefaultNodeServiceCapability = []csi.NodeServiceCapability_RPC_Type{
	csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
	csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
	csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
}

var DefaultPluginCapability = []*csi.PluginCapability{
	{
		Type: &csi.PluginCapability_Service_{
			Service: &csi.PluginCapability_Service{
				Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
			},
		},
	},
	{
		Type: &csi.PluginCapability_VolumeExpansion_{
			VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
				Type: csi.PluginCapability_VolumeExpansion_OFFLINE,
			},
		},
	},
	{
		Type: &csi.PluginCapability_Service_{
			Service: &csi.PluginCapability_Service{
				Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
			},
		},
	},
}

const (
	BasicVolumeType    VolumeType = 1
	StandardVolumeType VolumeType = 2
	DefaultVolumeType             = StandardVolumeType
	BasicNvmeName      string     = "Basic NVMe SSD"
	StandardNvmeName   string     = "Standard NVMe SSD"
)

var VolumeTypeName = map[VolumeType]string{
	BasicVolumeType:    BasicNvmeName,
	StandardVolumeType: StandardNvmeName,
}

type VolumeType int

func (v VolumeType) Int() int {
	return int(v)
}

func (v VolumeType) String() string {
	return VolumeTypeName[v]
}

func StringToType(s string) VolumeType {
	if BasicNvmeName == s {
		return BasicVolumeType
	} else if s == StandardNvmeName {
		return StandardVolumeType
	} else {
		panic("VolumeType ERROR")
	}
}

func (v VolumeType) IsValid() bool {
	if _, ok := VolumeTypeName[v]; !ok {
		return false
	} else {
		return true
	}
}

const (
	BasicVmType   VmType = 1
	DefaultVmType VmType = BasicVmType
	BasicVmName   string = "BasicVm"
)

var VmTypeName = map[VmType]string{
	BasicVmType: BasicVmName,
}

var VmTypeValue = map[string]VmType{
	BasicVmName: BasicVmType,
}

type VmType int

func (v VmType) Int() int {
	return int(v)
}

func (v VmType) IsVaild() bool {
	if _, ok := VmTypeName[v]; !ok {
		return false
	} else {
		return true
	}
}

// no use, reserve
var VmTypeAttachPreferred = map[VmType]VolumeType{
	BasicVmType: BasicVolumeType,
}

// topology disk mapping to vm
var VolumeTypeAttachConstraint = map[VolumeType][]VmType{
	BasicVolumeType: {
		BasicVmType,
	},
	StandardVolumeType: {
		BasicVmType,
	},
}

const (
	ZEC_MAX_DISK_SIZE_BYTES int64 = 32768 * 1024 * 1024 * 1024
	ZEC_MIN_DISK_SIZE_BYTES int64 = 20 * 1024 * 1024 * 1024
)

/*min zec cloud disk size*/
var VolumeTypeToMinSize = map[VolumeType]int64{
	BasicVolumeType:    ZEC_MIN_DISK_SIZE_BYTES,
	StandardVolumeType: ZEC_MIN_DISK_SIZE_BYTES,
}

/*max zec cloud disk size*/
var VolumeTypeToMaxSize = map[VolumeType]int64{
	BasicVolumeType:    ZEC_MAX_DISK_SIZE_BYTES,
	StandardVolumeType: ZEC_MAX_DISK_SIZE_BYTES,
}

/*disk rep*/
const (
	DiskSingleReplicaType  int = 1
	DiskMultiReplicaType   int = 2
	DiskThreeReplicaType   int = 3
	DefaultDiskReplicaType int = DiskMultiReplicaType
)
