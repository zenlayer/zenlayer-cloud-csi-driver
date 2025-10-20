// +-------------------------------------------------------------------------
// | Copyright (C) 2025 Zenlayer, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/rpcserver"
	"k8s.io/klog"
)

var (
	driverName           = flag.String("drivername", common.DefaultProvisionName, "name of csi driver.")
	endpoint             = flag.String("endpoint", "unix://csi/csi.sock", "CSI UnixSocket endpoint.")
	maxVolume            = flag.Int64("maxvolume", 9, "Maximum number of volumes that controller can publish to the node.")
	nodeNmae             = flag.String("nodename", "", "localhost hostname.")
	retryDetachTimesMax  = flag.Int("retry-detach-times-max", 10, "Maximum retry times of failed detach volume. Set to 0 to disable the limit.")
	secretpathAK         = flag.String("secretpathAK", "/etc/zec_secret/AccessKeyID", "zec console accesskeyId file.")
	secretpathPW         = flag.String("secretpathPW", "/etc/zec_secret/AccessKeyPassword", "zec console accesskeyPassword file.")
	driverType           = flag.String("drivertype", "", "zec csi driver type, node or controller.")
	mixDriver            = flag.Bool("mixDriver", false, "")
	featureGates         = flag.String("featureGates", "", "Alpha feature")
	defaultZone          = flag.String("defaultZone", "", "")
	defaultResourceGroup = flag.String("defaultGroup", "", "")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	//verify args
	if *driverType != common.NodeDriverType && *driverType != common.ControllerDriverType {
		klog.Fatalf("ERROR: main() missing drivertype")
		return
	}
	common.MixDriver = *mixDriver
	common.DefaultZone = *defaultZone
	common.DefaultResourceGroup = *defaultResourceGroup
	if common.DefaultZone != "" {
		if common.DefaultResourceGroup == "" {
			klog.Fatalf("ERROR: main() miss defaultzone and defaultResourceGroup")
			return
		}
	} else {
		if common.DefaultResourceGroup != "" {
			klog.Fatalf("ERROR: main() miss defaultzone and defaultResourceGroup")
			return
		}
	}
	if *featureGates != "" {
		common.FeatureGates = strings.Split(*featureGates, ",")
	}

	//verify Vm platform and check local-env
	zec_platform, err := common.VerifyEnv()
	if err != nil {
		klog.Fatal(err.Error())
		return
	}

	//get local vmid and the zone to which the vm belongs
	zec_vm_id, zec_vm_zone, err := common.GetZecEnv()
	if err != nil {
		klog.Fatal(err.Error())
		return
	}

	//if csi driver type is controller get zenlayer console accesskey and password
	var ak, pw string
	if *driverType == common.ControllerDriverType {
		ak, pw, err = common.GetZecSecret(*secretpathAK, *secretpathPW)
		if err != nil {
			klog.Fatal(err.Error())
			return
		}
	}

	//init zec cloud manager
	cm, err := cloud.NewZecCloudManager(string(ak[:]), string(pw[:]), *driverType)
	if err != nil {
		klog.Fatal(err.Error())
		return
	}
	if *driverType == common.ControllerDriverType {
		trycount := 0
		for {
			err = cm.Probe()
			if err != nil {
				klog.Errorf("ERROR: main() controller probe error return [%s].after 3s try again.", err.Error())
				trycount++
				time.Sleep(3 * time.Second)
			} else {
				break
			}

			if trycount > 5 {
				klog.Fatal(err)
				return
			}
		}
	}

	klog.Infof("Starting ZenlayerCSI Version[%s], AccessKeyID[%s], AccessKeyPassword[%s], localhostname[%s], Platform[%s], VmID[%s], VmZone[%s], featureGates[%v], defaultzone[%s], defaultResourceGroup[%s], mixDriver[%v], drivetype[%s], maxvol[%d]",
		common.Version, ak, pw, *nodeNmae, zec_platform, zec_vm_id, zec_vm_zone, common.FeatureGates, common.DefaultZone, common.DefaultResourceGroup, common.MixDriver, *driverType, *maxVolume)

	diskDriverInput := &driver.InitDiskDriverInput{
		Name:          *driverName,
		Version:       common.Version,
		NodeId:        zec_vm_id,
		ZoneId:        zec_vm_zone,
		MaxVolume:     *maxVolume,
		VolumeCap:     driver.DefaultVolumeAccessModeType,
		ControllerCap: driver.DefaultControllerServiceCapability,
		NodeCap:       driver.DefaultNodeServiceCapability,
		PluginCap:     driver.DefaultPluginCapability,
	}

	mounter := common.NewSafeMounter()
	driver := driver.GetDiskDriver()
	driver.InitDiskDriver(diskDriverInput)

	rpcserver.Run(driver, cm, mounter, *endpoint, *retryDetachTimesMax, *driverType)
	os.Exit(0)
}
