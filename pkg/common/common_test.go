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

func TestRetryLimiter(t *testing.T) {
	maxRetry := 4
	r := NewRetryLimiter(maxRetry)
	r.Add("a")
	r.Add("a")
	r.Add("a")
	r.Add("b")
	r.Add("a")
	if r.Try("a") != true {
		t.Errorf("expect true but actually false")
	}
	r.Add("a")
	if r.Try("a") != false {
		t.Errorf("expect false but actually true")
	}
}

func TestUnlimitedRetryLimiter(t *testing.T) {
	r := NewRetryLimiter(0)
	key := "key"
	for i := 1; i <= 5; i++ {
		r.Add(key)
		current := r.GetCurrentRetryTimes(key)
		if current != i {
			t.Errorf("expect %d but actually %d", i, current)
		}
		if r.Try(key) != true {
			t.Errorf("expect true but actually false")
		}
	}
}
