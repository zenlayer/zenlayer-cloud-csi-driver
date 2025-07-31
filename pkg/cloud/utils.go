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
func WaitForSpecificOrError(f func() (bool, error), timeout time.Duration, waitInterval time.Duration) error {
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

// WaitForSpecific wait a function return true.
func WaitForSpecific(f func() bool, timeout time.Duration, waitInterval time.Duration) error {
	return WaitForSpecificOrError(func() (bool, error) {
		return f(), nil
	}, timeout, waitInterval)
}

// WaitFor wait a function return true.
func WaitFor(f func() (bool, error)) error {
	return WaitForSpecificOrError(f, 180*time.Second, 3*time.Second)
}
