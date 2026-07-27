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
	"fmt"
	"time"
)

type TimeoutError struct {
	timeout time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Wait timeout [%s] ", e.timeout)
}

func (e *TimeoutError) Timeout() time.Duration { return e.timeout }

func NewTimeoutError(timeout time.Duration) *TimeoutError {
	return &TimeoutError{timeout: timeout}
}

// WaitForSpecificOrError wait a function return true or error.
//
// f is evaluated once up front and only then on every waitInterval tick, so an
// operation that already reached its target state returns immediately instead of
// paying a full interval of latency on every single call.
func WaitForSpecificOrError(f func() (bool, error), timeout time.Duration, waitInterval time.Duration) error {
	stop, err := f()
	if err != nil {
		return err
	}
	if stop {
		return nil
	}

	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			stop, err := f()
			if err != nil {
				return err
			}
			if stop {
				return nil
			}
		case <-timer.C:
			return NewTimeoutError(timeout)
		}
	}
}

/*
每次等待一个云盘/快照状态迁移的上限。

	驱动内部的等待时长必须小于 sidecar 的 --timeout, 否则 sidecar 会先判超时并重试,
	而驱动这一侧仍在持有 volId 锁继续等待, 重试请求立刻拿到 Aborted(operation
	pending), 白白浪费一轮 backoff。改动这两个常量时必须同步 chart 里的 sidecar
	--timeout(见 chart/templates/provisioner.yaml 的注释):

	  CreateVolume / AttachVolume / DetachVolume / ResizeVolume
	      各等待一次状态迁移            → 上限 WaitStatusTimeout(180s)
	  CreateSnapshot / DeleteSnapshot
	      各等待一次状态迁移            → 上限 WaitStatusTimeout(180s)
	  DeleteVolume
	      回收站机制需要释放 maxReleaseDiskTimes(2) 次, 每次等一轮
	                                    → 上限 2 * WaitStatusTimeout = 360s
*/
const (
	WaitStatusTimeout  = 180 * time.Second
	WaitStatusInterval = 3 * time.Second
)

// WaitFor wait a function return true.
func WaitFor(f func() (bool, error)) error {
	return WaitForSpecificOrError(f, WaitStatusTimeout, WaitStatusInterval)
}
