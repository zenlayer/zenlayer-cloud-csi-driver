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
	"reflect"
	"testing"

	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
)

func TestNewDefaultStorageClassFromType(t *testing.T) {
	tests := []struct {
		name     string
		diskType VolumeType
		sc       *ZecStorageClass
	}{
		{
			name:     "normal",
			diskType: DefaultVolumeType,
			sc: &ZecStorageClass{
				diskType:     DefaultVolumeType,
				maxSize:      VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:      VolumeTypeToMinSize[DefaultVolumeType],
				fsType:       common.DefaultFileSystem,
				zoneID:       "",
				placeGroupID: "",
			},
		},
		{
			name:     "normal2",
			diskType: 1,
			sc: &ZecStorageClass{
				diskType:     1,
				maxSize:      VolumeTypeToMaxSize[1],
				minSize:      VolumeTypeToMinSize[1],
				fsType:       common.DefaultFileSystem,
				zoneID:       "",
				placeGroupID: "",
			},
		},
		{
			name:     "invalid volume type",
			diskType: 10,
			sc:       nil,
		},
	}

	for _, test := range tests {
		res := NewDefaultZecStorageClassFromType(test.diskType)
		if !reflect.DeepEqual(test.sc, res) {
			t.Errorf("name %s: expect %v, but actually %v", test.name, test.sc, res)
		}
	}
}
