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

import (
	"testing"
)

func TestEntryFun(t *testing.T) {
	tests := []struct {
		name     string
		funcname string
	}{
		{
			name:     "normal",
			funcname: "CreateVolume",
		},
		{
			name:     "no fun name",
			funcname: "",
		},
	}
	for _, v := range tests {
		info, hash := EntryFunction(v.funcname)
		t.Logf("name %s: info %s, hash %s", v.name, info, hash)
	}
}

func TestExitFun(t *testing.T) {
	tests := []struct {
		name     string
		funcname string
		hash     string
	}{
		{
			name:     "normal",
			funcname: "CreateVolume",
			hash:     "hqp29hw2",
		},
		{
			name:     "no fun name",
			funcname: "",
			hash:     "hqp29hw2",
		},
	}
	for _, v := range tests {
		info := ExitFunction(v.funcname, v.hash)
		t.Logf("name %s: info %s", v.name, info)
	}
}

// TestParseCsiVolId 锁定 volume id 的解析语义。
//
// 重点是"非法 id 必须返回错误": 调用方依赖这个错误把请求映射成各自 RPC 规范要求的
// 返回码(卸载/删除类幂等成功, 其余 NotFound)。之前的实现只打日志、返回 nil error,
// 会把截断出来的垃圾 volid 直接发给云 API。
func TestParseCsiVolId(t *testing.T) {
	const (
		validVolId  = "1440908808376556310"  // 19 位
		validSerial = "d100om84ggf2oqdh05eg" // 20 位
	)

	tests := []struct {
		name       string
		csiVolId   string
		wantVolId  string
		wantSerial string
		wantErr    bool
	}{
		{
			name:       "normal",
			csiVolId:   validVolId + "-" + validSerial,
			wantVolId:  validVolId,
			wantSerial: validSerial,
		},
		{
			// serial 里含 "-" 时必须在第一个 "-" 处切分。老实现用 Split 取首尾元素,
			// 这里会把 serial 错切成 "5eg"。
			name:       "serial contains a dash",
			csiVolId:   validVolId + "-" + "d100om84ggf2oqdh-5eg",
			wantVolId:  validVolId,
			wantSerial: "d100om84ggf2oqdh-5eg",
		},
		{
			name:     "empty id",
			csiVolId: "",
			wantErr:  true,
		},
		{
			name:     "no separator",
			csiVolId: validVolId + validSerial,
			wantErr:  true,
		},
		{
			// csi-sanity DefaultIDGenerator.GenerateInvalidVolumeID()
			name:     "csi-sanity invalid volume id",
			csiVolId: "fake-vol-id",
			wantErr:  true,
		},
		{
			// csi-sanity DefaultIDGenerator.GenerateUniqueValidVolumeID(): 一个 uuid。
			// 对 CO 而言格式合法, 但不是本驱动发出的 id, 必须报错好让调用方回 NotFound。
			name:     "csi-sanity unique valid volume id (uuid)",
			csiVolId: "8e5f5ec3-4a1f-4a8f-9d0e-1b57cab76750",
			wantErr:  true,
		},
		{
			name:     "volid too short",
			csiVolId: "144090880837655" + "-" + validSerial,
			wantErr:  true,
		},
		{
			name:     "volid too long",
			csiVolId: validVolId + "0" + "-" + validSerial,
			wantErr:  true,
		},
		{
			name:     "serial too short",
			csiVolId: validVolId + "-" + "d100om84ggf2oqdh",
			wantErr:  true,
		},
		{
			name:     "empty serial",
			csiVolId: validVolId + "-",
			wantErr:  true,
		},
		{
			name:     "empty volid",
			csiVolId: "-" + validSerial,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		volId, serial, err := ParseCsiVolId(test.csiVolId)
		if test.wantErr {
			if err == nil {
				t.Errorf("name %s: expect error, but actually got volid[%s] serial[%s]", test.name, volId, serial)
			}
			continue
		}
		if err != nil {
			t.Errorf("name %s: expect no error, but actually %v", test.name, err)
			continue
		}
		if volId != test.wantVolId || serial != test.wantSerial {
			t.Errorf("name %s: expect volid[%s] serial[%s], but actually volid[%s] serial[%s]",
				test.name, test.wantVolId, test.wantSerial, volId, serial)
		}
	}
}

// TestGenCsiVolIdRoundTrip 保证 GenCsiVolId 生成的 id 一定能被 ParseCsiVolId 还原,
// 这是存量 PV 升级后仍然可用的前提(volume id 是持久化在 PV 上的)。
func TestGenCsiVolIdRoundTrip(t *testing.T) {
	volId := "1440908808376556310"
	serial := "d100om84ggf2oqdh05eg"

	gotVolId, gotSerial, err := ParseCsiVolId(GenCsiVolId(volId, serial))
	if err != nil {
		t.Fatalf("round trip failed: %v", err)
	}
	if gotVolId != volId || gotSerial != serial {
		t.Errorf("round trip mismatch: expect volid[%s] serial[%s], but actually volid[%s] serial[%s]",
			volId, serial, gotVolId, gotSerial)
	}
}

func TestGenerateHashInEightBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		hash  string
	}{
		{
			name:  "normal",
			input: "snapshot",
			hash:  "2aa38b8d",
		},
		{
			name:  "empty input",
			input: "",
			hash:  "811c9dc5",
		},
	}
	for _, v := range tests {
		res := GenerateHashInEightBytes(v.input)
		if v.hash != res {
			t.Errorf("name %s: expect %s but actually %s", v.name, v.hash, res)
		}
	}
}
