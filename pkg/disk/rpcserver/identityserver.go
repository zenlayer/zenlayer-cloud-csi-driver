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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/cloud"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

type IdentityServer struct {
	driver *driver.DiskDriver
	cloud  cloud.CloudManager
}

func NewIdentityServer(d *driver.DiskDriver, c cloud.CloudManager) *IdentityServer {
	return &IdentityServer{
		driver: d,
		cloud:  c,
	}
}

var _ csi.IdentityServer = &IdentityServer{}

func (is *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	funcName := "IdentityServer:Probe:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))

	var ready bool = true
	var err error = nil
	if !is.cloud.IsController() {
		klog.Infof("IdentityServer:Probe: Node driver Probe")
	} else {
		zones, err := is.cloud.GetZoneList()
		if err != nil {
			klog.Errorf("IdentityServer:Probe: Controller driver Probe error, lost console.zenlayer.com connect. error[%s]", err.Error())
			ready = false
		} else {
			klog.Infof("IdentityServer:Probe: Controller driver Probe, zonelist[%v]", zones)
		}
	}

	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: ready},
	}, err
}

/*
Action: Get plugin capabilities: CONTROLLER, ACCESSIBILITY, EXPANSION
*/
func (d *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	funcName := "IdentityServer:GetPluginCapabilities:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: d.driver.GetPluginCapability(),
	}, nil
}

/*
action: describe pv show
*/
func (d *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	funcName := "IdentityServer:GetPluginInfo:"
	info, hash := common.EntryFunction(funcName)
	klog.Info(info)
	defer klog.Info(common.ExitFunction(funcName, hash))

	if d.driver.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "ERROR: Missing Driver Name.")
	}

	if d.driver.GetVersion() == "" {
		return nil, status.Error(codes.InvalidArgument, "ERROR: Missing Driver Version.")
	}

	return &csi.GetPluginInfoResponse{
		Name:          d.driver.GetName(),
		VendorVersion: d.driver.GetVersion(),
	}, nil
}
