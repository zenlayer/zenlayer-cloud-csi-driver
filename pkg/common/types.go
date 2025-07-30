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

package common

var FeatureGates []string
var DefaultZone string
var DefaultResourceGroup string
var MixDriver bool

const FEAT_MOUNTOPT_ENABLE = "enable_mount_opt"

const Int64Max = int64(^uint64(0) >> 1)

const DefaultTimeFormat = "2006-01-02 15:04:05"

const KubeletRootDir = "/var/lib/kubelet"

const KubernetesPluginPathPrefix = "/plugins/kubernetes.io/"

const (
	Kib int64 = 1024
	Mib int64 = Kib * 1024
	Gib int64 = Mib * 1024
	Tib int64 = Gib * 1024
)

const ZECVOLID_LEN = 19     //zec DiskId len 		example:1440908808376556310
const ZECVOLSERIAL_LEN = 20 //zec Disk Serial len 	example:d100om84ggf2oqdh05eg
const ZECVMID_LEN = 19      //zec VmID len			example:1455441132925494374

const (
	Version              string = "v1.0.0"
	DefaultProvisionName string = "disk.csi.zenlayer.com"
	NodeDriverType       string = "node"
	ControllerDriverType string = "controller"
)

const (
	FileSystemExt3    string = "ext3"
	FileSystemExt4    string = "ext4"
	FileSystemXfs     string = "xfs"
	DefaultFileSystem string = FileSystemExt4
)

const (
	ZEC_PLATFORM            string = "Zenlayer Elastic Compute"
	ZEC_SYS_VENDOR_PATH     string = "/sys/class/dmi/id/sys_vendor"     //Vm platform
	ZEC_PRODUCT_FAMILY_PATH string = "/sys/class/dmi/id/product_family" //Vm zone
	ZEC_PRODUCT_SERIAL_PATH string = "/sys/class/dmi/id/product_serial" //Vm console ID
)
