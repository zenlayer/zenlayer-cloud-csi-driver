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
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog"
)

func EntryFunction(functionName string) (info string, hash string) {
	current := time.Now()
	hash = GenerateHashInEightBytes(current.UTC().String())
	return fmt.Sprintf("Enter %s at %s hash %s ", functionName,
		current.Format(DefaultTimeFormat), hash), hash
}

func ExitFunction(functionName, hash string) (info string) {
	current := time.Now()
	return fmt.Sprintf("Exit %s at %s hash %s ", functionName,
		current.Format(DefaultTimeFormat), hash)
}

func GenerateHashInEightBytes(input string) string {
	h := fnv.New32a()
	h.Write([]byte(input))
	return fmt.Sprintf("%.8x", h.Sum32())
}

func RetryOnError(backoff wait.Backoff, fn func() error) error {
	return retry.OnError(backoff, func(e error) bool {
		return true
	}, fn)
}

type retryLimiter struct {
	record   map[string]int
	maxRetry int
	mux      sync.RWMutex
}

type RetryLimiter interface {
	Add(id string)
	Try(id string) bool
	GetMaxRetryTimes() int
	GetCurrentRetryTimes(id string) int
}

func NewRetryLimiter(maxRetry int) RetryLimiter {
	return &retryLimiter{
		record:   map[string]int{},
		maxRetry: maxRetry,
		mux:      sync.RWMutex{},
	}
}

func (r *retryLimiter) Add(id string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.record[id]++
}

func (r *retryLimiter) Try(id string) bool {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.maxRetry == 0 || r.record[id] <= r.maxRetry
}

func (r *retryLimiter) GetMaxRetryTimes() int {
	return r.maxRetry
}

func (r *retryLimiter) GetCurrentRetryTimes(id string) int {
	return r.record[id]
}

func GenCsiVolId(volid, serial string) (csivolId string) {
	return volid + "-" + serial
}

func ParseCsiVolId(csivolId string) (volid string, serial string, err error) {
	s := strings.Split(csivolId, "-")
	volid = s[0]
	serial = s[len(s)-1]
	if len(volid) != ZECVOLID_LEN || len(serial) != ZECVOLSERIAL_LEN {
		klog.Errorf("ERROR:ParseCsiVolId volume id len or serial len error. %d/%d", len(volid), len(serial)) //not return
	}
	return volid, serial, nil
}

func VerifyEnv() (platform string, err error) {

	zec_platform, err := os.ReadFile(ZEC_SYS_VENDOR_PATH)
	if err != nil {
		klog.Errorf("ERROR: VerifyEnv Read conf err.path=%s", ZEC_SYS_VENDOR_PATH)
		return "", fmt.Errorf("ERROR: VerifyEnv Read conf err.path=%s", ZEC_SYS_VENDOR_PATH)
	}
	platform = strings.Replace(string(zec_platform[:]), "\n", "", -1)
	if platform != ZEC_PLATFORM {
		klog.Errorf("ERROR: Unsupported virtual machine platforms or unsupported virtual machine versions")
		return "", fmt.Errorf("ERROR: Unsupported virtual machine platforms or unsupported virtual machine versions")
	}

	//check blockdev cmd
	cmd_blockdev := exec.Command("blockdev", "-V")
	out, err := cmd_blockdev.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ERROR: blockdev cmd error:%s, out:%s"+err.Error(), out)
	}

	//check lsblk cmd
	cmd_lsblk := exec.Command("lsblk", "-V")
	out, err = cmd_lsblk.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ERROR: lsblk cmd error:%s, out:%s"+err.Error(), out)
	}

	return platform, nil
}

func GetZecEnv() (vmid string, vm_zone string, err error) {
	zec_vm_id, err := os.ReadFile(ZEC_PRODUCT_SERIAL_PATH)
	if err != nil {
		klog.Errorf("ERROR: GetZecEnv Read conf err.path=%s.", ZEC_PRODUCT_SERIAL_PATH)
		return "", "", fmt.Errorf("ERROR: GetZecEnv Read conf err.path=%s", ZEC_PRODUCT_SERIAL_PATH)
	}
	vmid = strings.Replace(string(zec_vm_id[:]), "\n", "", -1)
	if len(vmid) != ZECVMID_LEN {
		klog.Errorf("ERROR: GetZecEnv product_serial is invalid.vmid=%s", vmid)
		return "", "", fmt.Errorf("ERROR: GetZecEnv product_serial is invalid.vmid=%s", vmid)
	}

	zec_vm_zone, err := os.ReadFile(ZEC_PRODUCT_FAMILY_PATH)
	if err != nil {
		klog.Errorf("ERROR: GetZecEnv Read conf err.path=%s.", ZEC_PRODUCT_FAMILY_PATH)
		return "", "", fmt.Errorf("ERROR: GetZecEnv Read conf err.path=%s", ZEC_PRODUCT_FAMILY_PATH)
	}
	vm_zone = strings.Replace(string(zec_vm_zone[:]), "\n", "", -1)
	if len(vm_zone) == 0 {
		klog.Errorf("ERROR: GetZecEnv product_family is nil.")
		return "", "", fmt.Errorf("ERROR: GetZecEnv product_family is nil")
	}

	return vmid, vm_zone, nil
}

func GetZecSecret(akpath string, skpath string) (ak string, pw string, err error) {
	//get access key id and access key password
	zec_ak, err := os.ReadFile(akpath)
	if err != nil {
		klog.Errorf("ERROR: GetZecSecret Read secretAk err.path=%s.", akpath)
		return "", "", fmt.Errorf("ERROR: GetZecSecret Read secretAk err.path=%s", akpath)
	}

	zec_pw, err := os.ReadFile(skpath)
	if err != nil {
		klog.Errorf("ERROR: GetZecSecret Read secretPw err.path=%s.", skpath)
		return "", "", fmt.Errorf("ERROR: GetZecSecret Read secretPw err.path=%s", skpath)
	}

	return string(zec_ak[:]), string(zec_pw[:]), nil
}
