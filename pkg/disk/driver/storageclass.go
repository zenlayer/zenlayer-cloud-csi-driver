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
	StorageClassMaxSizeName  = "maxSize"      //cloud disk min size (bytes)
	StorageClassMinSizeName  = "minSize"      //cloud disk max size (bytes)
	StorageClassFsTypeName   = "fsType"       //ext3 ext4 xfs
	StorageClassZoneId       = "zoneID"       //zone
	StorageClassPlaceGroupID = "placeGroupID" //groupid
)

type ZecStorageClass struct {
	diskType     VolumeType
	maxSize      int64
	minSize      int64
	fsType       string
	zoneID       string
	placeGroupID string
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
	}
}

func UpdateParmsZone(opt map[string]string, zoneID string) {
	if _, ok := opt[StorageClassZoneId]; ok {
		opt[StorageClassZoneId] = zoneID
	}
}

func NewZecStorageClassFromMap(opt map[string]string) (*ZecStorageClass, error) {
	volType := -1
	var maxSize int64 = ZEC_MAX_DISK_SIZE_BYTES
	var minSize int64 = ZEC_MIN_DISK_SIZE_BYTES
	fsType := "ext4"
	zoneID := ""
	placeGroupID := ""
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
		}
	}
	if zoneID == "" || placeGroupID == "" {
		klog.Infof("INFO:NewZecStorageClassFromMap storage-class not config zoneID and placeGroupID. use default %s/%s", common.DefaultZone, common.DefaultResourceGroup)

		//keep Consistent
		zoneID = ""
		placeGroupID = ""

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
	if maxSize > 0 && minSize > 0 {
		err = sc.setTypeSize(maxSize, minSize)
		if err != nil {
			return nil, fmt.Errorf("setTypeSize error")
		}
	}

	err = sc.setFsType(fsType) //just set, get no use
	if err != nil {
		return nil, fmt.Errorf("setFsType error")
	}

	sc.SetZone(zoneID)
	sc.SetPlaceGroupID(placeGroupID)

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
	if sc.maxSize < sc.minSize {
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

func (sc ZecStorageClass) FormatVolumeSizeByte(sizeByte int64) int64 {
	if sizeByte <= sc.GetMinSizeByte() {
		sizeByte = sc.GetMinSizeByte()
	}
	if sizeByte > sc.GetMaxSizeByte() {
		sizeByte = sc.GetMaxSizeByte()
	}
	return sizeByte
}

func (sc ZecStorageClass) GetRequiredVolumeSizeByte(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return int64(sc.minSize), nil
	}
	res := int64(0)
	if capRange.GetRequiredBytes() > 0 {
		res = capRange.GetRequiredBytes()
	}
	res = sc.FormatVolumeSizeByte(res)
	if capRange.GetLimitBytes() > 0 && res > capRange.GetLimitBytes() {
		return -1, fmt.Errorf("ERROR:GetRequiredVolumeSizeByte volume required bytes %d greater than limit bytes %d", res, capRange.GetLimitBytes())
	}
	if res < ZEC_MIN_DISK_SIZE_BYTES || res > ZEC_MAX_DISK_SIZE_BYTES {
		return -1, fmt.Errorf("ERROR:GetRequiredVolumeSizeByte return size error, size=%d", res)
	}
	return res, nil
}
