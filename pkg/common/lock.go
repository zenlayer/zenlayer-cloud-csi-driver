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
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	OperationPendingFmt = "already an operation for the specified resource %s"
)

type ResourceLocks struct {
	locks sets.String
	mux   sync.Mutex
}

func NewResourceLocks() *ResourceLocks {
	return &ResourceLocks{
		locks: sets.NewString(),
	}
}

func (lock *ResourceLocks) TryAcquire(resourceID string) bool {
	lock.mux.Lock()
	defer lock.mux.Unlock()
	if lock.locks.Has(resourceID) {
		return false
	}
	lock.locks.Insert(resourceID)
	return true
}

func (lock *ResourceLocks) Release(resourceID string) {
	lock.mux.Lock()
	defer lock.mux.Unlock()
	lock.locks.Delete(resourceID)
}
