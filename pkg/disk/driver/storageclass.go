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
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"k8s.io/klog"
)

const (
	StorageClassTypeName     = "type"         //1=basic, 2=standard
	StorageClassFsTypeName   = "fsType"       //ext3 ext4 xfs
	StorageClassZoneId       = "zoneID"       //zone
	StorageClassPlaceGroupID = "placeGroupID" //groupid
	StorageClassBurstEnable  = "burstEnable"  //burst
)

type ZecStorageClass struct {
	diskType     VolumeType
	maxSize      int64
	minSize      int64
	fsType       string
	zoneID       string
	placeGroupID string
	burstEnable  bool
}

func NewDefaultZecStorageClassFromType(diskType VolumeType) *ZecStorageClass {
	if !diskType.IsValid() {
		return nil
	}
	return &ZecStorageClass{
		diskType:     diskType,
		maxSize:      VolumeTypeToMaxSize[diskType],
		minSize:      VolumeTypeToMinSize[diskType],
		fsType:       common.DefaultFileSystem,
		zoneID:       "",
		placeGroupID: "",
		burstEnable:  false,
	}
}

// UpdateParmsZone 把 storage-class 参数里已有的 zoneID 改写成本次实际选中的 zone。
func UpdateParmsZone(opt map[string]string, zoneID string) {
	for k := range opt {
		if strings.EqualFold(k, StorageClassZoneId) {
			opt[k] = zoneID
		}
	}
}

func NewZecStorageClassFromMap(opt map[string]string) (*ZecStorageClass, error) {
	volType := -1
	fsType := "ext4"
	zoneID := ""
	placeGroupID := ""
	burstenable := "false"
	var err error

	for k, v := range opt {
		switch strings.ToLower(k) {
		case strings.ToLower(StorageClassTypeName):
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			volType = iv
		case strings.ToLower(StorageClassZoneId):
			zoneID = v
		case strings.ToLower(StorageClassPlaceGroupID):
			placeGroupID = v
		case strings.ToLower(StorageClassBurstEnable):
			burstenable = v
		case strings.ToLower(StorageClassFsTypeName):
			if v != "" {
				fsType = v
			}
		}
	}
	if zoneID == "" || placeGroupID == "" {
		klog.Infof("INFO:NewZecStorageClassFromMap storage-class not config zoneID and placeGroupID. use default-zone[%s], use default-resourceGroup[%s]", common.DefaultZone, common.DefaultResourceGroup)

		zoneID = common.DefaultZone
		placeGroupID = common.DefaultResourceGroup
		if zoneID == "" || placeGroupID == "" {
			return nil, fmt.Errorf("storageclass missing zoneID or placeGroupID")
		}
	}

	var t VolumeType //Basic 1 or Standard 2
	if volType == -1 {
		t = DefaultVolumeType
	} else {
		t = VolumeType(volType)
	}

	if !t.IsValid() {
		return nil, fmt.Errorf("unsupported volume type %d", volType)
	}
	sc := NewDefaultZecStorageClassFromType(t)

	err = sc.setTypeSize(ZEC_MAX_DISK_SIZE_BYTES, ZEC_MIN_DISK_SIZE_BYTES)
	if err != nil {
		return nil, fmt.Errorf("setTypeSize error")
	}
	err = sc.setFsType(fsType) //just set, get no use
	if err != nil {
		return nil, fmt.Errorf("setFsType error")
	}

	sc.SetZone(zoneID)
	sc.SetPlaceGroupID(placeGroupID)
	if burstenable == "true" {
		sc.SetBurstEnable(true)
	} else {
		sc.SetBurstEnable(false)
	}

	return sc, nil
}

func (sc ZecStorageClass) GetDiskType() VolumeType {
	return sc.diskType
}

func (sc ZecStorageClass) ConvertToDiskCategory(vt VolumeType) string {
	return vt.String()
}

func (sc ZecStorageClass) GetMinSizeByte() int64 {
	return int64(sc.minSize)
}

func (sc ZecStorageClass) GetMaxSizeByte() int64 {
	return int64(sc.maxSize)
}

func (sc ZecStorageClass) GetFsType() string {
	return sc.fsType
}

func (sc ZecStorageClass) GetZone() string {
	return sc.zoneID
}

func (sc ZecStorageClass) GetPlaceGroupID() string {
	return sc.placeGroupID
}

func (sc ZecStorageClass) GetBurstEnable() bool {
	return sc.burstEnable
}

func (sc *ZecStorageClass) setFsType(fs string) error {
	if !IsValidFileSystemType(fs) {
		return fmt.Errorf("unsupported filesystem type %s", fs)
	}
	sc.fsType = fs
	return nil
}

func (sc *ZecStorageClass) setTypeSize(maxSize, minSize int64) error {
	if maxSize < 0 || minSize <= 0 {
		return nil
	}
	if maxSize < minSize {
		return fmt.Errorf("max size must greater than or equal to min size")
	}
	sc.maxSize, sc.minSize = maxSize, minSize
	return nil
}

func (sc *ZecStorageClass) SetZone(zone string) {
	sc.zoneID = zone
}

func (sc *ZecStorageClass) SetPlaceGroupID(placeGroupID string) {
	sc.placeGroupID = placeGroupID
}

func (sc *ZecStorageClass) SetBurstEnable(burstEnable bool) {
	sc.burstEnable = burstEnable
}

func (sc ZecStorageClass) FormatVolumeSizeByte(sizeByte int64) int64 {
	if sizeByte <= sc.GetMinSizeByte() {
		sizeByte = sc.GetMinSizeByte()
	}
	if sizeByte > sc.GetMaxSizeByte() {
		sizeByte = sc.GetMaxSizeByte()
	}
	return sizeByte
}

/*
GetRequiredVolumeSizeByte 把 CO 给的容量区间换算成云盘实际能创建的字节数。

	返回值保证是"整数 GiB 对齐"的实际容量, 调用方可以直接把它回填到
	CreateVolumeResponse.Volume.CapacityBytes, 不需要再做一次取整,
	也不会出现"上报值与云上实际容量不一致"的问题(见 ControllerExpandVolume
	的同样处理)。

	取整语义:
	  1. requiredBytes 小于云盘最小规格时向上抬到最小规格(而不是直接报错
	     OutOfRange —— 否则用户写 10Gi 的 PVC 会永久 Pending)。
	  2. 再向上取整到整数 GiB, 因为云 API 的 DiskSize 单位就是 GiB。
	  3. 只有在取整后的容量超出 limitBytes 或超出云盘最大规格时才返回错误,
	     此时确实无法满足请求。

	limitBytes 必须参与判断: 之前的实现完全忽略它, 会返回一个大于用户上限的容量。
*/
func (sc ZecStorageClass) GetRequiredVolumeSizeByte(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return sc.GetMinSizeByte(), nil
	}

	requiredBytes := capRange.GetRequiredBytes()
	limitBytes := capRange.GetLimitBytes()
	if requiredBytes < 0 || limitBytes < 0 {
		return -1, fmt.Errorf("capacity range [%d,%d] should not be less than zero", requiredBytes, limitBytes)
	}
	if limitBytes > 0 && requiredBytes > limitBytes {
		return -1, fmt.Errorf("volume required bytes %d greater than limit bytes %d", requiredBytes, limitBytes)
	}

	// 向上抬到云盘最小规格, 再向上取整到整数 GiB
	if requiredBytes < sc.GetMinSizeByte() {
		requiredBytes = sc.GetMinSizeByte()
	}
	actualBytes := common.GibToByte(common.ByteCeilToGib(requiredBytes))

	if actualBytes > sc.GetMaxSizeByte() {
		return -1, fmt.Errorf("required size %d bytes exceeds max volume size %d bytes", actualBytes, sc.GetMaxSizeByte())
	}
	if limitBytes > 0 && actualBytes > limitBytes {
		return -1, fmt.Errorf("volume size %d bytes (rounded up to whole GiB, min volume size is %d bytes) exceeds limit bytes %d",
			actualBytes, sc.GetMinSizeByte(), limitBytes)
	}

	return actualBytes, nil
}
