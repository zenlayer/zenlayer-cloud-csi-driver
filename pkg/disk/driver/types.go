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

// cloud disk only support ReadWriteOnce
var DefaultVolumeAccessModeType = []csi.VolumeCapability_AccessMode_Mode{
	csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
}

// controller type support capability/feature
var DefaultControllerServiceCapability = []csi.ControllerServiceCapability_RPC_Type{
	csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
	//csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
	//csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
}

// node type support capability/feature
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
	BasicVolumeType    VolumeType = 1 //经济型
	StandardVolumeType VolumeType = 2 //标准型
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

const (
	DiskSingleReplicaType  int = 1
	DiskMultiReplicaType   int = 2
	DiskThreeReplicaType   int = 3
	DefaultDiskReplicaType int = DiskMultiReplicaType
)
