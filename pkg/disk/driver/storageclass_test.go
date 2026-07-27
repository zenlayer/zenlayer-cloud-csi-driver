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

	"github.com/container-storage-interface/spec/lib/go/csi"
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

// TestUpdateParmsZone 锁定两条语义:
//  1. key 匹配必须大小写无关, 与 NewZecStorageClassFromMap 的解析行为保持一致;
//  2. 参数里原本没有配 zoneID 时不新增 key(WaitForFirstConsumer 下 zone 由 topology
//     决定, UpdateParmsZone 的职责只是把已有的旧值改写正确)。
func TestUpdateParmsZone(t *testing.T) {
	const newZone = "asia-north-1a"

	tests := []struct {
		name   string
		opt    map[string]string
		expect map[string]string
	}{
		{
			name:   "exact key",
			opt:    map[string]string{"zoneID": "na-west-1a"},
			expect: map[string]string{"zoneID": newZone},
		},
		{
			name:   "lower case key",
			opt:    map[string]string{"zoneid": "na-west-1a"},
			expect: map[string]string{"zoneid": newZone},
		},
		{
			name:   "upper case key",
			opt:    map[string]string{"ZONEID": "na-west-1a"},
			expect: map[string]string{"ZONEID": newZone},
		},
		{
			name:   "mixed case key",
			opt:    map[string]string{"ZoneId": "na-west-1a"},
			expect: map[string]string{"ZoneId": newZone},
		},
		{
			name:   "key absent is not added",
			opt:    map[string]string{"type": "1", "placeGroupID": "xxx"},
			expect: map[string]string{"type": "1", "placeGroupID": "xxx"},
		},
		{
			name:   "other keys untouched",
			opt:    map[string]string{"zoneID": "na-west-1a", "type": "2", "burstEnable": "true"},
			expect: map[string]string{"zoneID": newZone, "type": "2", "burstEnable": "true"},
		},
		{
			name:   "empty map",
			opt:    map[string]string{},
			expect: map[string]string{},
		},
		{
			name:   "nil map does not panic",
			opt:    nil,
			expect: nil,
		},
	}

	for _, test := range tests {
		UpdateParmsZone(test.opt, newZone)
		if !reflect.DeepEqual(test.expect, test.opt) {
			t.Errorf("name %s: expect %v, but actually %v", test.name, test.expect, test.opt)
		}
	}
}

// TestGetRequiredVolumeSizeByte 锁定容量取整语义: 向上抬到最小规格 -> 向上取整到
// 整数 GiB -> 校验 limitBytes 和最大规格。返回值必须是云盘实际容量, 以便
// CreateVolume / ControllerExpandVolume 直接回填给 CO。
func TestGetRequiredVolumeSizeByte(t *testing.T) {
	sc := NewDefaultZecStorageClassFromType(DefaultVolumeType)
	minByte := sc.GetMinSizeByte()
	maxByte := sc.GetMaxSizeByte()

	tests := []struct {
		name     string
		capRange *csi.CapacityRange
		expect   int64
		wantErr  bool
	}{
		{
			name:     "nil capacity range returns min size",
			capRange: nil,
			expect:   minByte,
		},
		{
			name:     "empty capacity range returns min size",
			capRange: &csi.CapacityRange{},
			expect:   minByte,
		},
		{
			name:     "exactly min size",
			capRange: &csi.CapacityRange{RequiredBytes: minByte},
			expect:   minByte,
		},
		{
			name:     "below min size is raised to min size",
			capRange: &csi.CapacityRange{RequiredBytes: 10 * common.Gib},
			expect:   minByte,
		},
		{
			name:     "below min size but limit forbids raising",
			capRange: &csi.CapacityRange{RequiredBytes: 10 * common.Gib, LimitBytes: 15 * common.Gib},
			wantErr:  true,
		},
		{
			name:     "below min size and limit exactly allows min size",
			capRange: &csi.CapacityRange{RequiredBytes: 10 * common.Gib, LimitBytes: minByte},
			expect:   minByte,
		},
		{
			name:     "non whole gib is rounded up",
			capRange: &csi.CapacityRange{RequiredBytes: 20*common.Gib + 512*common.Mib},
			expect:   21 * common.Gib,
		},
		{
			name:     "rounded up size exceeds limit",
			capRange: &csi.CapacityRange{RequiredBytes: 20*common.Gib + 512*common.Mib, LimitBytes: 20*common.Gib + 512*common.Mib},
			wantErr:  true,
		},
		{
			name:     "whole gib within limit",
			capRange: &csi.CapacityRange{RequiredBytes: 21 * common.Gib, LimitBytes: 21 * common.Gib},
			expect:   21 * common.Gib,
		},
		{
			name:     "only limit bytes provided",
			capRange: &csi.CapacityRange{LimitBytes: 30 * common.Gib},
			expect:   minByte,
		},
		{
			name:     "required greater than limit",
			capRange: &csi.CapacityRange{RequiredBytes: 100 * common.Gib, LimitBytes: 50 * common.Gib},
			wantErr:  true,
		},
		{
			name:     "exactly max size",
			capRange: &csi.CapacityRange{RequiredBytes: maxByte},
			expect:   maxByte,
		},
		{
			name:     "exceed max size",
			capRange: &csi.CapacityRange{RequiredBytes: maxByte + common.Gib},
			wantErr:  true,
		},
		{
			name:     "negative required bytes",
			capRange: &csi.CapacityRange{RequiredBytes: -1},
			wantErr:  true,
		},
		{
			name:     "negative limit bytes",
			capRange: &csi.CapacityRange{RequiredBytes: minByte, LimitBytes: -1},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		res, err := sc.GetRequiredVolumeSizeByte(test.capRange)
		if test.wantErr {
			if err == nil {
				t.Errorf("name %s: expect error, but actually got size %d", test.name, res)
			}
			continue
		}
		if err != nil {
			t.Errorf("name %s: expect no error, but actually %v", test.name, err)
			continue
		}
		if res != test.expect {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.expect, res)
		}
		// 返回值必须是整数 GiB 对齐的, 否则 CreateVolume 回填给 CO 的容量会和云上不一致
		if common.GibToByte(common.ByteCeilToGib(res)) != res {
			t.Errorf("name %s: size %d is not whole GiB aligned", test.name, res)
		}
	}
}
