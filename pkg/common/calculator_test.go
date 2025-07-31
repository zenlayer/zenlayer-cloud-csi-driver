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

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func TestGibToByte(t *testing.T) {
	tests := []struct {
		name  string
		gib   int
		bytes int64
	}{
		{
			name:  "normal",
			gib:   23,
			bytes: 23 * Gib,
		},
		{
			name:  "large number",
			gib:   65536000,
			bytes: 65536000 * Gib,
		},
		{
			name:  "zero Gib",
			gib:   0,
			bytes: 0,
		},
		{
			name:  "minus Gib",
			gib:   -24,
			bytes: -24 * Gib,
		},
	}

	for _, test := range tests {
		res := GibToByte(test.gib)
		if test.bytes != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.bytes, res)
		}
	}
}

func TestByteCeilToGib(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		gib   int
	}{
		{
			name:  "normal",
			bytes: 1 * Gib,
			gib:   1,
		},
		{
			name:  "ceil to gib",
			bytes: 1*Gib + 1,
			gib:   2,
		},
		{
			name:  "zero bytes",
			bytes: 0,
			gib:   0,
		},
		{
			name:  "min value",
			bytes: -1,
			gib:   0,
		},
	}
	for _, test := range tests {
		res := ByteCeilToGib(test.bytes)
		if test.gib != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.gib, res)
		}
	}
}

func TestIsValidCapacityBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		capRange *csi.CapacityRange
		isValid  bool
	}{
		{
			name:  "normal",
			bytes: 1 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 1 * Gib,
				LimitBytes:    1 * Gib,
			},
			isValid: true,
		},
		{
			name:  "invalid range",
			bytes: 10 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 11 * Gib,
				LimitBytes:    10 * Gib,
			},
			isValid: false,
		},
		{
			name:     "empty range",
			bytes:    10 * Gib,
			capRange: &csi.CapacityRange{},
			isValid:  true,
		},
		{
			name:     "nil range",
			bytes:    1 * Gib,
			capRange: nil,
			isValid:  true,
		},
		{
			name:  "without floor",
			bytes: 10 * Gib,
			capRange: &csi.CapacityRange{
				LimitBytes: 10*Gib + 1,
			},
			isValid: true,
		},
		{
			name:  "invalid floor",
			bytes: 11 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 11*Gib + 1,
			},
			isValid: false,
		},
		{
			name:  "without ceil",
			bytes: 1 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 1 * Gib,
			},
			isValid: true,
		},
		{
			name:  "invalid ceil",
			bytes: 1 * Gib,
			capRange: &csi.CapacityRange{
				LimitBytes: 1*Gib - 1,
			},
			isValid: false,
		},
	}

	for _, test := range tests {
		res := IsValidCapacityBytes(test.bytes, test.capRange)
		if test.isValid != res {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isValid, res)
		}
	}
}

func TestGetRequestSizeBytes(t *testing.T) {
	tests := []struct {
		name     string
		capRange *csi.CapacityRange
		bytes    int64
	}{
		{
			name: "normal",
			capRange: &csi.CapacityRange{
				RequiredBytes: 20 * Gib,
				LimitBytes:    20 * Gib,
			},
			bytes: 20 * Gib,
		},
		{
			name:     "empty range",
			capRange: &csi.CapacityRange{},
			bytes:    0,
		},
		{
			name:     "nil range",
			capRange: nil,
			bytes:    0,
		},
		{
			name: "normal range",
			capRange: &csi.CapacityRange{
				RequiredBytes: 22 * Gib,
				LimitBytes:    23 * Gib,
			},
			bytes: 22 * Gib,
		},
		{
			name: "invalid range",
			capRange: &csi.CapacityRange{
				RequiredBytes: 11 * Gib,
				LimitBytes:    10 * Gib,
			},
			bytes: -1,
		},
		{
			name: "less then zero",
			capRange: &csi.CapacityRange{
				RequiredBytes: -11 * Gib,
				LimitBytes:    -10 * Gib,
			},
			bytes: -1,
		},
	}

	for _, test := range tests {
		res, _ := GetRequestSizeBytes(test.capRange)
		if test.bytes != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.bytes, res)
		}
	}
}
