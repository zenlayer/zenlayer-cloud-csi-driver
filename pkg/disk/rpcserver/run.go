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
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"k8s.io/kubernetes/pkg/util/mount"
)

func Run(driver *driver.DiskDriver, cloud cloud.CloudManager, mounter *mount.SafeFormatAndMount,
	endpoint string, retryTimesMax int, driverType string) {

	var servers Servers

	servers.IdentityServer = NewIdentityServer(driver, cloud)
	if driverType == common.ControllerDriverType {
		servers.ControllerServer = NewControllerServer(driver, cloud, retryTimesMax)
		if common.MixDriver {
			servers.NodeServer = NewNodeServer(driver, cloud, mounter)
		}
	} else {
		servers.NodeServer = NewNodeServer(driver, cloud, mounter)
	}

	s := NewNonBlockingGRPCServer()
	s.Start(endpoint, servers)
	s.Wait()
}
