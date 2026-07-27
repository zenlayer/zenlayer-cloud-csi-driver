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

// entryTimes correlates a function invocation's entry time with the unique hash
// returned by EntryFunction, so ExitFunction can report the elapsed duration
// without changing the (funcName, hash) call signature used across the codebase.
var (
	entryTimes   = make(map[string]time.Time)
	entryTimesMu sync.Mutex
)

func EntryFunction(functionName string) (info string, hash string) {
	current := time.Now().UTC()
	hash = GenerateHashInEightBytes(current.UTC().String())
	entryTimesMu.Lock()
	entryTimes[hash] = current
	entryTimesMu.Unlock()
	return fmt.Sprintf("Enter[%s], UTC Time[%s], Hash[%s]", functionName, current.Format(DefaultTimeFormat), hash), hash
}

func ExitFunction(functionName, hash string) (info string) {
	current := time.Now().UTC()

	entryTimesMu.Lock()
	start, ok := entryTimes[hash]
	if ok {
		delete(entryTimes, hash)
	}
	entryTimesMu.Unlock()

	if ok {
		return fmt.Sprintf("Exit[%s], UTC Time[%s], Hash[%s], Elapsed[%dms]", functionName, current.Format(DefaultTimeFormat), hash, current.Sub(start).Milliseconds())
	}
	return fmt.Sprintf("Exit[%s], UTC Time[%s], Hash[%s]", functionName, current.Format(DefaultTimeFormat), hash)
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

func GenCsiVolId(volid, serial string) (csivolId string) {
	return volid + "-" + serial
}

// ParseCsiVolId 把 GenCsiVolId 生成的 csi volume id 还原成云盘 id 和 serial。
//
// id 只会以 "<diskId>-<serial>" 的形式生成, 且 diskId 是纯数字不含 "-", 所以在第一个
// "-" 处切分总能正确还原两段(用 strings.Cut 而不是 Split 取首尾, 后者在 serial 含 "-"
// 时会错位)。
//
// 形状不符的 id 不可能对应任何 ZEC 云盘, 因此这里返回错误, 而不是像之前那样只打一条
// 日志、把截断出来的垃圾 volid 继续发给云 API(那样云端只会回一个语义无关的错误)。
// 调用方必须按各自 RPC 的规范要求映射这个错误:
//   - 卸载/删除类 RPC 要求幂等, 应视作"目标已不存在"直接返回成功;
//   - 其余 RPC 应返回 NotFound —— 这个 id 对 CO 而言格式合法, 只是不由本驱动发出,
//     即该卷在本驱动这里不存在。
func ParseCsiVolId(csivolId string) (volid string, serial string, err error) {
	volid, serial, found := strings.Cut(csivolId, "-")
	if !found {
		return "", "", fmt.Errorf("invalid csi volume id[%s]: want format <diskId>-<serial>", csivolId)
	}
	if len(volid) != ZECVOLID_LEN || len(serial) != ZECVOLSERIAL_LEN {
		return "", "", fmt.Errorf("invalid csi volume id[%s]: disk id[%s] len[%d] want[%d], serial[%s] len[%d] want[%d]",
			csivolId, volid, len(volid), ZECVOLID_LEN, serial, len(serial), ZECVOLSERIAL_LEN)
	}
	return volid, serial, nil
}

func VerifyEnv() (platform string, err error) {

	zec_platform, err := os.ReadFile(ZEC_SYS_VENDOR_PATH)
	if err != nil {
		klog.Errorf("ERROR:VerifyEnv() Read conf err. path[%s]", ZEC_SYS_VENDOR_PATH)
		return "", fmt.Errorf("ERROR:VerifyEnv() Read conf err. path[%s]", ZEC_SYS_VENDOR_PATH)
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
		klog.Errorf("ERROR:VerifyEnv() missing cmd blockdev")
		return "", fmt.Errorf("ERROR: blockdev cmd error: %v, out: %s", err, out)
	}

	//check lsblk cmd
	cmd_lsblk := exec.Command("lsblk", "-V")
	out, err = cmd_lsblk.CombinedOutput()
	if err != nil {
		klog.Errorf("ERROR:VerifyEnv() missing cmd lsblk")
		return "", fmt.Errorf("ERROR: lsblk cmd error: %v, out: %s", err, out)
	}

	return platform, nil
}

func GetZecEnv() (vmid string, vm_zone string, err error) {
	zec_vm_id, err := os.ReadFile(ZEC_PRODUCT_SERIAL_PATH)
	if err != nil {
		klog.Errorf("ERROR: GetZecEnv() Read conf err. path[%s]", ZEC_PRODUCT_SERIAL_PATH)
		return "", "", fmt.Errorf("ERROR: GetZecEnv Read conf err. path[%s]", ZEC_PRODUCT_SERIAL_PATH)
	}
	vmid = strings.Replace(string(zec_vm_id[:]), "\n", "", -1)
	if len(vmid) != ZECVMID_LEN {
		klog.Errorf("ERROR: GetZecEnv() product_serial is invalid. vmid[%s]", vmid)
		return "", "", fmt.Errorf("ERROR: GetZecEnv product_serial is invalid. vmid[%s]", vmid)
	}

	zec_vm_zone, err := os.ReadFile(ZEC_PRODUCT_FAMILY_PATH)
	if err != nil {
		klog.Errorf("ERROR: GetZecEnv() Read conf err. path[%s]", ZEC_PRODUCT_FAMILY_PATH)
		return "", "", fmt.Errorf("ERROR: GetZecEnv Read conf err. path[%s]", ZEC_PRODUCT_FAMILY_PATH)
	}
	vm_zone = strings.Replace(string(zec_vm_zone[:]), "\n", "", -1)
	if len(vm_zone) == 0 {
		klog.Errorf("ERROR: GetZecEnv() product_family is nil.")
		return "", "", fmt.Errorf("ERROR: GetZecEnv product_family is nil")
	}

	return vmid, vm_zone, nil
}

func GetZecSecret(akpath string, skpath string) (ak string, pw string, err error) {
	//get access key id and access key password
	zec_ak, err := os.ReadFile(akpath)
	if err != nil {
		klog.Errorf("ERROR: GetZecSecret() Read secretAk err. path[%s]", akpath)
		return "", "", fmt.Errorf("ERROR: GetZecSecret Read secretAk err. path[%s]", akpath)
	}

	zec_pw, err := os.ReadFile(skpath)
	if err != nil {
		klog.Errorf("ERROR: GetZecSecret() Read secretPw err. path[%s]", skpath)
		return "", "", fmt.Errorf("ERROR: GetZecSecret Read secretPw err. path[%s]", skpath)
	}

	// 裁剪首尾空白/换行
	return strings.TrimSpace(string(zec_ak)), strings.TrimSpace(string(zec_pw)), nil
}
