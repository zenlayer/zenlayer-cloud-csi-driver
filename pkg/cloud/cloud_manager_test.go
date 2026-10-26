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

import (
	"testing"

	csicommon "github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
)

func TestNewZecCloudManagerOpenAPIDomain(t *testing.T) {
	for _, domain := range []string{"", "portal-dev.example.com"} {
		cm, err := NewZecCloudManager("ak", "sk", csicommon.ControllerDriverType, domain)
		if err != nil || cm == nil || cm.zecClient == nil || !cm.Iscontroller {
			t.Fatalf("domain %q: manager = %+v, err = %v; want a controller client", domain, cm, err)
		}
	}
	cm, err := NewZecCloudManager("", "", csicommon.NodeDriverType, "portal-dev.example.com")
	if err != nil || cm.zecClient != nil || cm.Iscontroller {
		t.Fatalf("node manager = %+v, err = %v; want no client", cm, err)
	}
}
