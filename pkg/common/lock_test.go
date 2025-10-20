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

import "testing"

func TestResourceLocks_TryAcquire(t *testing.T) {
	lock := NewResourceLocks()
	tests := []struct {
		name       string
		resourceId string
		isAcquire  bool
	}{
		{
			name:       "first acquire",
			resourceId: "vol-1",
			isAcquire:  true,
		},
		{
			name:       "acquire fail",
			resourceId: "vol-1",
			isAcquire:  false,
		},
		{
			name:       "acquire another",
			resourceId: "vol-2",
			isAcquire:  true,
		},
	}

	for _, test := range tests {
		res := lock.TryAcquire(test.resourceId)
		if test.isAcquire != res {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isAcquire, res)
		}
	}
}

func TestResourceLocks_Release(t *testing.T) {
	lock := NewResourceLocks()
	resourceA := "vol-1"
	resourceB := "vol-2"

	lock.Release(resourceA)

	if res := lock.TryAcquire(resourceB); !res {
		t.Errorf("try acquire %s failed", resourceB)
	}

	lock.Release(resourceB)

	lock.Release(resourceB)
}
