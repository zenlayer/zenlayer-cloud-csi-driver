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
	"errors"
	"fmt"

	csicommon "github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/common"
	"github.com/zenlayer/zenlayer-cloud-csi-driver/pkg/disk/driver"
	common "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/common"
	zec "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/zec20250901"
	"k8s.io/klog"
)

var _ CloudManager = &zecCloudManager{}

type zecCloudManager struct {
	zecClient    *zec.Client
	Iscontroller bool
	//...
}

/*
Action: Base request conf
zenlayer sdk describe(v0.2.0):

	type Config struct {
		AutoRetry     bool `default:"false"`
		MaxRetryTime  int  `default:"3"`
		RetryDuration time.Duration
		HttpTransport *http.Transport   `default:""`
		Transport     http.RoundTripper `default:""`
		Proxy         string            `default:""`
		Scheme        string            `default:"HTTPS"`
		Domain        string            `default:""`
		// timeout in seconds
		Timeout int `default:"300"`
		Debug   *bool
	}

	type BaseRequest struct {
		httpMethod  string
		scheme      string
		domain      string
		params      map[string]string
		formParams  map[string]string
		headers     map[string]string
		contentType string
		body        []byte
		path        string
		action      string
		serviceName string
		apiVersion  string
		autoRetries bool `default:"false"`
		maxAttempts int  `default:"0"`
	}

	BaseRequest: SetMaxAttempts(int)/SetAutoRetries(bool)/GetApiVersion()(string)/GetMaxAttempts()(int)/GetAutoRetries()(bool)
*/
func NewZecCloudManager(CLOUDAK string, CLOUDSK string, drivertype string) (*zecCloudManager, error) {
	cm := &zecCloudManager{
		zecClient:    nil,
		Iscontroller: true,
	}
	if drivertype == csicommon.NodeDriverType {
		cm.Iscontroller = false
		return cm, nil
	}
	var err error

	config := common.NewConfig().WithTimeout(300).WithAutoRetry(false).WithMaxRetryTime(3)

	cloud_secretKeyId := CLOUDAK
	cloud_secretKeyPassword := CLOUDSK

	cm.zecClient, err = zec.NewClient(config, cloud_secretKeyId, cloud_secretKeyPassword)
	if err != nil || cm.zecClient == nil {
		klog.Errorf("ERROR:Init zec.NewClient error[%v]", err.Error())
		return nil, err
	}
	return cm, nil
}

/*
Action: check connect to console.zenlayer.com, return Cloud Disk support Region list

Args:

Return:

zenlayer sdk describe(v0.2.0):

	type DescribeDiskRegionsRequest struct {
	    *common.BaseRequest
	}

	type DescribeDiskRegionsResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		RegionIds []string `json:"regionIds,omitempty"`
	}

	type DescribeDiskRegionsResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *DescribeDiskRegionsResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) Probe() error {
	funcName := "zecCloudManager:Probe:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	request := zec.NewDescribeDiskRegionsRequest()

	response, err := cm.zecClient.DescribeDiskRegions(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		klog.Errorf("%s API ERROR, error[%v]", ERRORLOG, err.Error())
		return err
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v]", ERRORLOG, err.Error())
		return err
	}
	_ = response
	return err
}

func (cm *zecCloudManager) GetZoneList() ([]string, error) {
	funcName := "zecCloudManager:GetZoneList:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	request := zec.NewDescribeDiskRegionsRequest()

	response, err := cm.zecClient.DescribeDiskRegions(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v]", ERRORLOG, err.Error())
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v]", ERRORLOG, err.Error())
		return nil, err
	}
	if response.Response == nil {
		return nil, errors.New(ERRORLOG + " nil response")
	}

	return response.Response.RegionIds, nil
}

func (cm zecCloudManager) IsController() bool {
	return cm.Iscontroller
}

/*
Action: create cloud disk.

Args: volName string, volSize int, volCategory string, zoneId string, resouceGroupID string

Return: volId string, err error

zenlayer sdk describe(v0.2.0):

	type CreateDisksRequest struct {
	    *common.BaseRequest
	    ZoneId *string `json:"zoneId,omitempty"`							// ZoneId 云硬盘所属的可用区ID。
	    DiskName *string `json:"diskName,omitempty"`						// DiskName 云盘名称。范围1到64个字符。仅支持输入字母、数字、-/_和英文句点(.)。且必须以数字或字母开头和结尾。
	    DiskSize *int `json:"diskSize,omitempty"`							// DiskSize 云硬盘大小，单位GiB。
	    DiskAmount *int `json:"diskAmount,omitempty"`						// DiskAmount 需要创建的云硬盘的数量。
	    InstanceId *string `json:"instanceId,omitempty"`					// InstanceId 云硬盘挂在的实例ID。
	    ResourceGroupId *string `json:"resourceGroupId,omitempty"`			// ResourceGroupId 云硬盘所在的资源组ID。如不指定则放入默认资源组。
	    DiskCategory *string `json:"diskCategory,omitempty"`				// DiskCategory 云硬盘种类。Basic NVMe SSD: 经济型 NVMe SSD。Standard NVMe SSD: 标准型 NVMe SSD。默认为Standard NVMe SSD。
	    SnapshotId *string `json:"snapshotId,omitempty"`					// SnapshotId 使用快照ID进行创建。如果传入则根据此快照创建云硬盘，快照的云盘类型必须为数据盘快照。
	    MarketingOptions *MarketingInfo `json:"marketingOptions,omitempty"`	// MarketingOptions 市场营销的相关选项。
	}

	type MarketingInfo struct {
		DiscountCode *string `json:"discountCode,omitempty"`				// DiscountCode 使用市场发放的折扣码。如果折扣码不存在，最终折扣将不会生效。
		UsePocVoucher *bool `json:"usePocVoucher,omitempty"`				// UsePocVoucher 是否使用POC代金券。 如果系统不存在POC代金券，相关创建流程会失败。
	}

	type CreateDisksResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		DiskIds []string `json:"diskIds,omitempty"`							// DiskIds 创建的云硬盘ID列表。
		OrderNumber *string `json:"orderNumber,omitempty"`					// OrderNumber 本次创建对应的订单编号。
	}

	type CreateDisksResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *CreateDisksResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) CreateVolume(volName string, volSize int, volCategory string, zoneId string, resouceGroupID string) (volId string, err error) {

	funcName := "zecCloudManager:CreateVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volName == "" || volSize == 0 || volCategory == "" || zoneId == "" || resouceGroupID == "" {
		return "", errors.New(ERRORLOG + "args error")
	}

	var diskcount int = 1
	request := zec.NewCreateDisksRequest()
	request.DiskName = &volName
	request.DiskSize = &volSize
	request.ZoneId = &zoneId
	request.DiskAmount = &diskcount
	request.ResourceGroupId = &resouceGroupID
	request.DiskCategory = &volCategory

	response, err := cm.zecClient.CreateDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
			return "", err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
		return "", err
	}
	if response.Response == nil {
		return "", errors.New(ERRORLOG + " nil response. volName=" + volName)
	}

	volId = response.Response.DiskIds[0]

	if err = cm.waitDiskStatus(DiskStatusAvailable, volId); err != nil {
		return "", err
	}

	return volId, nil
}

func (cm *zecCloudManager) CreateVolumeFromSnapshot(volName string, volSize int, volCategory string, zoneId string, resourceGroupID string, snapshotId string) (volId string, err error) {
	funcName := "zecCloudManager:CreateVolumeFromSnapshot:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if volName == "" || volSize == 0 || volCategory == "" || zoneId == "" || resourceGroupID == "" || snapshotId == "" {
		return "", errors.New(ERRORLOG + "args error")
	}

	var diskcount int = 1
	request := zec.NewCreateDisksRequest()
	request.DiskName = &volName
	request.DiskSize = &volSize
	request.ZoneId = &zoneId
	request.DiskAmount = &diskcount
	request.ResourceGroupId = &resourceGroupID
	request.DiskCategory = &volCategory
	request.SnapshotId = &snapshotId

	response, err := cm.zecClient.CreateDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
			return "", err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
		return "", err
	}
	if response.Response == nil {
		return "", errors.New(ERRORLOG + " nil response. volName=" + volName)
	}

	volId = response.Response.DiskIds[0]

	if err = cm.waitDiskStatus(DiskStatusAvailable, volId); err != nil {
		return "", err
	}
	_ = INFOLOG
	return volId, nil
}

/*
Action: delete cloud disk.

Args: volId

Return: err error

zenlayer sdk describe(v0.2.0):

	type ReleaseDiskRequest struct {
	    *common.BaseRequest
	    DiskId *string `json:"diskId,omitempty"`
	}

	type ReleaseDiskResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) DeleteVolume(volId string) (err error) {

	funcName := "zecCloudManager:DeleteVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if volId == "" {
		return errors.New(ERRORLOG + "args error")
	}

	request := zec.NewReleaseDiskRequest()
	request.DiskId = &volId

	response, err := cm.zecClient.ReleaseDisk(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil && err.(*common.ZenlayerCloudSdkError).Code == OPERATION_FAILED_RESOURCE_NOT_FOUND {
			klog.Infof("%s API INFO, volId[%s] Not Found", INFOLOG, volId)
			return nil
		}
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusDeleted, volId); err != nil {
		return err
	}
	_ = response
	return nil
}

/*
Action: attach cloud disk to Vm.

Args: volId, vmid

Return: err error

zenlayer sdk describe(v0.2.0):

	type AttachDisksRequest struct {
	    *common.BaseRequest
	    DiskIds []string `json:"diskIds,omitempty"`
	    InstanceId *string `json:"instanceId,omitempty"`
	}

	type AttachDisksResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) AttachVolume(volId string, vmId string) (err error) {

	funcName := "zecCloudManager:AttachVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volId == "" || vmId == "" {
		return errors.New(ERRORLOG + "args error" + ".VolId=" + volId + ".Vmid=" + vmId)
	}

	request := zec.NewAttachDisksRequest()
	request.DiskIds = append(request.DiskIds, volId)
	request.InstanceId = &vmId

	response, err := cm.zecClient.AttachDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volId[%s], vmId[%s]", ERRORLOG, err.Error(), volId, vmId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volId[%s], vmId[%s]", ERRORLOG, err.Error(), volId, vmId)
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusInUse, volId); err != nil {
		return err
	}
	_ = response
	return nil
}

/*
Action: detach cloud disk.

Args: volId

Return: err error

zenlayer sdk describe(v0.2.0):

	type DetachDisksRequest struct {
	    *common.BaseRequest
	    DiskIds []string `json:"diskIds,omitempty"`
	    InstanceCheckFlag *bool `json:"instanceCheckFlag,omitempty"`	// InstanceCheckFlag 是否检测实例的运行状态。 默认为true，即实例关机才允许被卸载。否则必须实例关机才能调用本接口。
	}

	type DetachDisksResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) DetachVolume(volId string) (err error) {

	funcName := "zecCloudManager:DetachVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volId == "" {
		return errors.New(ERRORLOG + "args error" + ".volId=" + volId)
	}

	request := zec.NewDetachDisksRequest()
	request.DiskIds = append(request.DiskIds, volId)
	var checkVmRunning bool = false
	request.InstanceCheckFlag = &checkVmRunning

	response, err := cm.zecClient.DetachDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusAvailable, volId); err != nil {
		return err
	}
	_ = response
	return nil
}

/*
Action: Resize cloud disk.

Args: volId, newSize

Return: err error

zenlayer sdk describe(v0.2.0):

	type ResizeDiskRequest struct {
	    *common.BaseRequest
	    DiskId *string `json:"diskId,omitempty"`	// DiskId 云硬盘ID。通过DescribeDisks接口查询。
	    DiskSize *int `json:"diskSize,omitempty"`	// DiskSize 云硬盘扩容后的大小。单位GiB。必须大于当前云硬盘大小。云盘最大限制为32768GB(32TB)。
	}

	type ResizeDiskResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) ResizeVolume(volId string, requestSize int) error {

	funcName := "zecCloudManager:ResizeVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volId == "" || requestSize == 0 {
		return errors.New(ERRORLOG + "args error" + ".volId=" + volId)
	}

	request := zec.NewResizeDiskRequest()
	request.DiskId = &volId
	request.DiskSize = &requestSize

	response, err := cm.zecClient.ResizeDisk(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusStable, volId); err != nil {
		return err
	}
	_ = response
	return nil
}

/*
Action: get disk info by volid.

Args: volId

Return: zecVolInfo *ZecVolume, err error

zenlayer sdk describe(v0.2.0):

	type DescribeDisksRequest struct {
	    *common.BaseRequest
	    DiskIds []string `json:"diskIds,omitempty"`									// DiskIds 根据云盘ID列表筛选。
	    DiskName *string `json:"diskName,omitempty"`								// DiskName 根据云盘名称筛选，该字段支持模糊搜索。
	    DiskStatus *string `json:"diskStatus,omitempty"`							// DiskStatus 根据云盘的状态进行筛选。
	    DiskType *string `json:"diskType,omitempty"`								// DiskType 根据云盘的类型进行筛选。
	    DiskCategory *string `json:"diskCategory,omitempty"`	 					// DiskCategory 根据云盘的分类进行筛选。
	    InstanceId *string `json:"instanceId,omitempty"`							// InstanceId 根据云盘挂载的实例ID进行筛选。
	    ZoneId *string `json:"zoneId,omitempty"`									// ZoneId 根据云盘所在的可用区进行筛选。
	    PageNum *int `json:"pageNum,omitempty"`										// PageNum 返回的分页大小。默认为20，最大为1000。
	    PageSize *int `json:"pageSize,omitempty"`	 								// PageSize 返回的分页数。默认为1。
	    RegionId *string `json:"regionId,omitempty"`								// RegionId 根据云盘所在的节点ID进行筛选。
	    SnapshotAbility *bool `json:"snapshotAbility,omitempty"`					// SnapshotAbility 根据云盘是否有快照能力进行筛选。
	    ResourceGroupId *string `json:"resourceGroupId,omitempty"`					// ResourceGroupId 根据快照所属的资源组进行筛选。
	}

	type DescribeDisksResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		TotalCount *int64 `json:"totalCount,omitempty"`								// TotalCount 符合条件的数据总数。
		DataSet []*DiskInfo `json:"dataSet,omitempty"`								// DataSet 云盘的结果集。
	}

	type DiskInfo struct {
		DiskId *string `json:"diskId,omitempty"`									// DiskId 云盘的 ID。
		DiskName *string `json:"diskName,omitempty"`								// DiskName 云盘的名称。
		RegionId *string `json:"regionId,omitempty"`								// RegionId 云盘所在的节点ID。
		ZoneId *string `json:"zoneId,omitempty"`									// ZoneId 云盘所在节点的可用区ID。
		DiskType *string `json:"diskType,omitempty"`								// DiskType 云盘的类型。
		Portable *bool `json:"portable,omitempty"`									// Portable 是否可卸载。
		DiskCategory *string `json:"diskCategory,omitempty"`						// DiskCategory 云盘的类别。
		DiskSize *int `json:"diskSize,omitempty"`									// DiskSize 云盘的大小。单位：GiB。
		DiskStatus *string `json:"diskStatus,omitempty"`							// DiskStatus 云盘的状态。
		InstanceId *string `json:"instanceId,omitempty"`							// InstanceId 云盘绑定实例的ID。
		InstanceName *string `json:"instanceName,omitempty"`						// InstanceName 云盘绑定实例的名称。
		CreateTime *string `json:"createTime,omitempty"`							// CreateTime 创建时间。
		ExpiredTime *string `json:"expiredTime,omitempty"`							// ExpiredTime 到期时间。
		Period *int `json:"period,omitempty"`										// Period 周期。
		ResourceGroupId *string `json:"resourceGroupId,omitempty"`					// ResourceGroupId 云盘所属的资源组ID。
		ResourceGroupName *string `json:"resourceGroupName,omitempty"`				// ResourceGroupName 云盘所属的资源组名称。
		Serial *string `json:"serial,omitempty"`									// Serial 云盘序号。可能为null，表示取不到值。
		SnapshotAbility *bool `json:"snapshotAbility,omitempty"`					// SnapshotAbility 是否具体快照能力。
		AutoSnapshotPolicyId *string `json:"autoSnapshotPolicyId,omitempty"`		// AutoSnapshotPolicyId 云盘关联的自动快照策略ID。
	}

	type DescribeDisksResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *DescribeDisksResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) FindVolume(volId string) (*ZecVolume, error) {

	funcName := "zecCloudManager:FindVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volId == "" {
		return nil, errors.New(ERRORLOG + "args error" + ".volid=" + volId)
	}

	request := zec.NewDescribeDisksRequest()
	request.DiskIds = append(request.DiskIds, volId)

	response, err := cm.zecClient.DescribeDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
		return nil, err
	}

	if response.Response.TotalCount == nil || response.Response == nil {
		return nil, errors.New(ERRORLOG + " nil response. volid=" + volId)
	}
	if *response.Response.TotalCount == 0 {
		return nil, nil
	}
	if *response.Response.TotalCount != 1 {
		return nil, errors.New(ERRORLOG + "Volid Repeat. volid=" + volId)
	}
	if response.Response.DataSet[0] == nil {
		return nil, errors.New(ERRORLOG + " nil response. volid=" + volId)
	}

	zecVolInfo := NewZecVolume()
	if response.Response.DataSet[0].DiskId != nil {
		zecVolInfo.ZecVolume_Id = *response.Response.DataSet[0].DiskId
	}
	if response.Response.DataSet[0].DiskName != nil {
		zecVolInfo.ZecVolume_Name = *response.Response.DataSet[0].DiskName
	}
	if response.Response.DataSet[0].DiskSize != nil {
		zecVolInfo.ZecVolume_Size = *response.Response.DataSet[0].DiskSize
	}
	if response.Response.DataSet[0].DiskCategory != nil {
		zecVolInfo.ZecVolume_Type = int(driver.StringToType(*response.Response.DataSet[0].DiskCategory))
	}
	if response.Response.DataSet[0].DiskStatus != nil {
		zecVolInfo.ZecVolume_Status = *response.Response.DataSet[0].DiskStatus
	}
	if response.Response.DataSet[0].ZoneId != nil {
		zecVolInfo.ZecVolume_Zone = *response.Response.DataSet[0].ZoneId
	}
	if response.Response.DataSet[0].InstanceId != nil {
		zecVolInfo.ZecVolume_InstanceId = *response.Response.DataSet[0].InstanceId
	}
	if response.Response.DataSet[0].Serial != nil {
		zecVolInfo.ZecVolume_Serial = *response.Response.DataSet[0].Serial
	}
	if response.Response.DataSet[0].Portable != nil {
		zecVolInfo.ZecVolume_Portable = *response.Response.DataSet[0].Portable
	}
	if response.Response.DataSet[0].SnapshotAbility != nil {
		zecVolInfo.ZecVolume_SnapshotAbility = *response.Response.DataSet[0].SnapshotAbility
	}
	if response.Response.DataSet[0].ResourceGroupId != nil {
		zecVolInfo.ZecVolume_ResourceGroupId = *response.Response.DataSet[0].ResourceGroupId
	}

	return zecVolInfo, nil
}

func (cm *zecCloudManager) FindVolumeByName(volName string, zoneID string, sizeGB int, diskType string, ResourceGroupId string) (*ZecVolume, error) {

	funcName := "zecCloudManager:FindVolumeByName:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if volName == "" {
		return nil, errors.New(ERRORLOG + "args error")
	}

	var diskname string = "^" + volName + "$"
	var pagesize int = 20
	var pagenum int = 1
	request := zec.NewDescribeDisksRequest()
	request.DiskName = &diskname
	request.ZoneId = &zoneID
	request.ResourceGroupId = &ResourceGroupId
	request.PageSize = &pagesize
	request.PageNum = &pagenum

	response, err := cm.zecClient.DescribeDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
		return nil, err
	}
	if response.Response.TotalCount == nil || response.Response == nil {
		return nil, errors.New(ERRORLOG + " nil response. volName=" + volName)
	}

	if *response.Response.TotalCount == 0 {
		return nil, nil
	}

	//用于保存所有数据
	count := int(*response.Response.TotalCount)
	var DiskListAll []*zec.DiskInfo = make([]*zec.DiskInfo, count)

	//如果需要分页
	if count > *request.PageSize {

		apicallnum := count / *request.PageSize
		if count%*request.PageSize != 0 {
			apicallnum++
		}

		//从第一页开始遍历
		var firstpage int = 1
		request.PageNum = &firstpage
		var end int = 0
		for ra := 0; ra < apicallnum; ra++ {
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
					return nil, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], volName[%s]", ERRORLOG, err.Error(), volName)
				return nil, err
			}
			klog.Infof("%s Traversal once,sum count[%d], pagenum[%d], pagesize[%d], rsp count[%d]", INFOLOG, count, request.PageNum, request.PageSize, len(response.Response.DataSet))
			*request.PageNum++
			for i := 0; i < len(response.Response.DataSet); i++ {
				DiskListAll[i+end] = response.Response.DataSet[i]
			}
			end += len(response.Response.DataSet)
		}

	} else { //如果小于等于pagesize不需要分页,直接用第一个请求赋值
		for i := 0; i < count; i++ {
			DiskListAll[i] = response.Response.DataSet[i]
		}
	}

	if len(DiskListAll) != count {
		klog.Errorf("%s API ERROR, Find DiskListAll len[%d] not equal to response.TotalCount[%d]", ERRORLOG, len(DiskListAll), count)
		return nil, fmt.Errorf("%s API ERROR, Find DiskListAll len[%d] not equal to response.TotalCount[%d]", ERRORLOG, len(DiskListAll), count)
	}

	zecVolInfos := make([]*ZecVolume, count)
	for c := 0; c < count; c++ {
		zecVolInfo := NewZecVolume()
		if DiskListAll[c].DiskId != nil {
			zecVolInfo.ZecVolume_Id = *DiskListAll[c].DiskId
		}
		if DiskListAll[c].DiskName != nil {
			zecVolInfo.ZecVolume_Name = *DiskListAll[c].DiskName
		}
		if DiskListAll[c].DiskSize != nil {
			zecVolInfo.ZecVolume_Size = *DiskListAll[c].DiskSize
		}
		if DiskListAll[c].DiskCategory != nil {
			zecVolInfo.ZecVolume_Type = int(driver.StringToType(*DiskListAll[c].DiskCategory))
		}
		if DiskListAll[c].DiskStatus != nil {
			zecVolInfo.ZecVolume_Status = *DiskListAll[c].DiskStatus
		}
		if DiskListAll[c].ZoneId != nil {
			zecVolInfo.ZecVolume_Zone = *DiskListAll[c].ZoneId
		}
		if DiskListAll[c].InstanceId != nil {
			zecVolInfo.ZecVolume_InstanceId = *DiskListAll[c].InstanceId
		}
		if DiskListAll[c].Serial != nil {
			zecVolInfo.ZecVolume_Serial = *DiskListAll[c].Serial
		}
		if DiskListAll[c].Portable != nil {
			zecVolInfo.ZecVolume_Portable = *DiskListAll[c].Portable
		}
		if DiskListAll[c].SnapshotAbility != nil {
			zecVolInfo.ZecVolume_SnapshotAbility = *DiskListAll[c].SnapshotAbility
		}
		if DiskListAll[c].ResourceGroupId != nil {
			zecVolInfo.ZecVolume_ResourceGroupId = *DiskListAll[c].ResourceGroupId
		}

		zecVolInfos[c] = zecVolInfo
		if zoneID == zecVolInfos[c].ZecVolume_Zone {
			if zecVolInfos[c].ZecVolume_Type == int(driver.StringToType(diskType)) {
				if zecVolInfos[c].ZecVolume_Size != sizeGB {
					//csi规定创建不允许创建同名不同大小的卷
					klog.Errorf("%s find different size but same volname volumeid[%s], volumename[%s], existvol size[%d], request new size[%d]", ERRORLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name, zecVolInfos[c].ZecVolume_Size, sizeGB)
					return zecVolInfos[c], fmt.Errorf("%s find different size but same volname volumeid[%s], volumename[%s], existvol size[%d], request new size[%d]", ERRORLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name, zecVolInfos[c].ZecVolume_Size, sizeGB)
				}
				if zecVolInfos[c].ZecVolume_Size == sizeGB {
					klog.Infof("%s find volumeid[%s], volumename[%s], serial[%s], size[%d], zone[%s].", INFOLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name,
						zecVolInfos[c].ZecVolume_Serial, zecVolInfos[c].ZecVolume_Size, zecVolInfos[c].ZecVolume_Zone)
					return zecVolInfos[c], nil
				}
			}
		}
	}

	return nil, nil
}

/*
Action: get vm info by vmid.

Args: vmid

Return: zecVmInfo *ZecVm, err error

zenlayer sdk describe(v0.2.0):

	type DescribeInstancesRequest struct {
	    *common.BaseRequest
	    InstanceIds []string `json:"instanceIds,omitempty"`										// InstanceIds 根据实例ID列表进行筛选。最大不能超过100个。
	    ZoneId *string `json:"zoneId,omitempty"`	 											// ZoneId 实例所属的可用区ID。
	    ImageId *string `json:"imageId,omitempty"`	 											// ImageId 镜像ID。
	    Ipv4Address *string `json:"ipv4Address,omitempty"`										// Ipv4Address 根据实例关联的IPv4过滤。
	    Ipv6Address *string `json:"ipv6Address,omitempty"`										// Ipv6Address 根据实例关联的IPv6信息过滤。
	    Status *string `json:"status,omitempty"`												// Status 根据实例的状态过滤。
	    Name *string `json:"name,omitempty"`													// Name 根据实例显示名称过滤。该字段支持模糊搜索。
	    PageSize *int `json:"pageSize,omitempty"`												// PageSize 返回的分页大小。
	    PageNum *int `json:"pageNum,omitempty"`	 												// PageNum 返回的分页数。
	    ResourceGroupId *string `json:"resourceGroupId,omitempty"`								// ResourceGroupId 根据资源组ID过滤。
	}

	type DescribeInstancesResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		TotalCount *int `json:"totalCount,omitempty"`											// TotalCount 符合条件的数据总数。
		DataSet []*InstanceInfo `json:"dataSet,omitempty"`										// DataSet 实例列表的数据。
	}

		type InstanceInfo struct {
			InstanceId *string `json:"instanceId,omitempty"`									// InstanceId 实例唯一ID。
			InstanceName *string `json:"instanceName,omitempty"`								// InstanceName 实例显示名称。
			ZoneId *string `json:"zoneId,omitempty"`											// ZoneId 实例所属的可用区ID。
			InstanceType *string `json:"instanceType,omitempty"`								// InstanceType CPU 规格。如果是GPU实例，该字段取值为null。
			Cpu *int `json:"cpu,omitempty"`														// Cpu CPU 核数。单位：个。
			Memory *int `json:"memory,omitempty"`												// Memory 内存容量。单位：GiB。
			ImageId *string `json:"imageId,omitempty"`											// ImageId 镜像ID。
			ImageName *string `json:"imageName,omitempty"`										// ImageName 镜像名称。
			TimeZone *string `json:"timeZone,omitempty"`										// TimeZone 设置的系统时区信息。
			NicNetworkType *string `json:"nicNetworkType,omitempty"`							// NicNetworkType 网卡模式。
			Status *string `json:"status,omitempty"`											// Status 实例状态。
			SystemDisk *SystemDisk `json:"systemDisk,omitempty"`								// SystemDisk 系统盘信息。
			DataDisks []*DataDisk `json:"dataDisks,omitempty"`									// DataDisks 实例上挂在的数据盘信息。
			PublicIpAddresses []string `json:"publicIpAddresses,omitempty"`						// PublicIpAddresses 实例上公网IPv4列表。
			PrivateIpAddresses []string `json:"privateIpAddresses,omitempty"`					// PrivateIpAddresses 实例上内网IP列表。
			KeyId *string `json:"keyId,omitempty"`												// KeyId 安装的SSH密钥ID。
			SubnetId *string `json:"subnetId,omitempty"`										// SubnetId 实例主网卡关联的子网ID。
			SecurityGroupId *string `json:"securityGroupId,omitempty"`							// SecurityGroupId 实例主网卡关联的安全组ID。
			EnableAgent *bool `json:"enableAgent,omitempty"`									// EnableAgent 是否开启QGA Agent。
			EnableAgentMonitor *bool `json:"enableAgentMonitor,omitempty"`						// EnableAgentMonitor 是否开启QGA 监控采集。
			EnableIpForward *bool `json:"enableIpForward,omitempty"`							// EnableIpForward 是否开启IP转发。
			CreateTime *string `json:"createTime,omitempty"`									// CreateTime 创建时间。
			ExpiredTime *string `json:"expiredTime,omitempty"`									// ExpiredTime 到期时间。
			ResourceGroupId *string `json:"resourceGroupId,omitempty"`							// ResourceGroupId 实例所属的资源组ID。
			ResourceGroupName *string `json:"resourceGroupName,omitempty"`						// ResourceGroupName 实例所属的资源组名称。
			Nics []*NicInfo `json:"nics,omitempty"`												// Nics 实例上绑定的网卡信息。
		}

		type DescribeInstancesResponse struct {
		    *common.BaseResponse
		    RequestId *string `json:"requestId,omitempty"`
		    Response *DescribeInstancesResponseParams `json:"response,omitempty"`
		}
*/
func (cm *zecCloudManager) FindInstance(vmid string) (*ZecVm, error) {

	funcName := "zecCloudManager:FindInstance:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if vmid == "" {
		return nil, errors.New(ERRORLOG + "args error" + ".Vmid=" + vmid)
	}

	request := zec.NewDescribeInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, vmid)

	response, err := cm.zecClient.DescribeInstances(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], vmid[%s]", ERRORLOG, err.Error(), vmid)
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], vmid[%s]", ERRORLOG, err.Error(), vmid)
		return nil, err
	}
	if response.Response.TotalCount == nil || response.Response == nil {
		return nil, errors.New(ERRORLOG + " nil response. vmid=" + vmid)
	}

	if *response.Response.TotalCount == 0 {
		return nil, nil
	}
	if *response.Response.TotalCount != 1 {
		return nil, errors.New(ERRORLOG + "Vmid Repeat, vmid=" + vmid)
	}
	if response.Response.DataSet[0] == nil {
		return nil, errors.New(ERRORLOG + " nil response. vmid=" + vmid)
	}

	zecVmInfo := NewZecVm()
	if response.Response.DataSet[0].InstanceId != nil {
		zecVmInfo.ZecVm_Id = *response.Response.DataSet[0].InstanceId
	}
	if response.Response.DataSet[0].Status != nil {
		zecVmInfo.ZecVm_Status = *response.Response.DataSet[0].Status
	}
	if response.Response.DataSet[0].ZoneId != nil {
		zecVmInfo.ZecVm_Zone = *response.Response.DataSet[0].ZoneId
	}

	zecVmInfo.ZecVm_Type = driver.BasicVmType.Int()

	return zecVmInfo, nil
}

/*
Action: get zec zone info.

Args: zoneid

Return: error

zenlayer sdk describe(v0.2.0):

	type DescribeZonesRequest struct {
	    *common.BaseRequest
	    ZoneIds []string `json:"zoneIds,omitempty"`							// ZoneIds 根据可用区ID过滤。
	}

	type DescribeZonesResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		ZoneSet []*ZoneInfo `json:"zoneSet,omitempty"`						// ZoneSet 可用区列表。
	}

	type ZoneInfo struct {
		ZoneId *string `json:"zoneId,omitempty"`							// ZoneId 可用区ID。
		RegionId *string `json:"regionId,omitempty"`						// RegionId 可用区所在的节点ID。
		ZoneName *string `json:"zoneName,omitempty"`						// ZoneName 可用区名称。
		SupportSecurityGroup *bool `json:"supportSecurityGroup,omitempty"`	// SupportSecurityGroup 可用区是否支持安全组。该字段已废弃，当前所有节点均支持安全组。
	}

	type DescribeZonesResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *DescribeZonesResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) GetZone(zoneId string) error {

	funcName := "zecCloudManager:GetZone:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	request := zec.NewDescribeZonesRequest()

	response, err := cm.zecClient.DescribeZones(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], zoneId[%s]", ERRORLOG, err.Error(), zoneId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], zoneId[%s]", ERRORLOG, err.Error(), zoneId)
		return err
	}
	_ = response
	return nil
}

/*
Action: check zec vm status, only support running vm.

Args: vmid(string)

Return: error

zenlayer sdk describe(v0.2.0):

	type DescribeInstancesStatusRequest struct {
	    *common.BaseRequest
	    InstanceIds []string `json:"instanceIds,omitempty"`						// InstanceIds 要查询的实例ID列表。
	    PageSize *int `json:"pageSize,omitempty"`								// PageSize 分页大小。
	    PageNum *int `json:"pageNum,omitempty"`									// PageNum 分页页数。
	    ResourceGroupId *string `json:"resourceGroupId,omitempty"`				// ResourceGroupId 根据资源组ID过滤。
	}

	type DescribeInstancesStatusResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		TotalCount *int `json:"totalCount,omitempty"`							// TotalCount 符合条件的数据总数。
		DataSet []*InstanceStatus `json:"dataSet,omitempty"`					// DataSet 实例状态数据。
	}

	type InstanceStatus struct {
		InstanceId *string `json:"instanceId,omitempty"`						// InstanceId 实例的ID。
		InstanceStatus *string `json:"instanceStatus,omitempty"`				// InstanceStatus 实例的状态。
	}

	type DescribeInstancesStatusResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *DescribeInstancesStatusResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) GetVmStatus(vmId string) (exist bool, status string, err error) {

	funcName := "zecCloudManager:GetVmStatus:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if vmId == "" {
		return false, "", errors.New(ERRORLOG + "error args" + ",Vmid=" + vmId)
	}

	request := zec.NewDescribeInstancesStatusRequest()
	request.InstanceIds = append(request.InstanceIds, vmId)

	response, err := cm.zecClient.DescribeInstancesStatus(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], vmId[%s]", ERRORLOG, err.Error(), vmId)
			return false, "", err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], vmId[%s]", ERRORLOG, err.Error(), vmId)
		return false, "", err
	}
	if response.Response.TotalCount == nil || response.Response == nil {
		return false, "", errors.New(ERRORLOG + " nil response. vmid=" + vmId)
	}

	if *response.Response.TotalCount == 0 {
		return false, "", nil
	}

	if *response.Response.TotalCount != 1 {
		return false, "", errors.New(ERRORLOG + "VmId repeat. vmid=" + vmId)
	}
	if response.Response.DataSet[0].InstanceStatus == nil {
		return false, "", errors.New(ERRORLOG + " nil response. vmid=" + vmId)
	}

	return true, *response.Response.DataSet[0].InstanceStatus, nil
}

/*
Action:	Create a snapshot from source volume

Args:	snapshotname, resource vol id

Return:	snapshot Id

zenlayer sdk describe(v0.2.0):

	type CreateSnapshotRequest struct {
	    *common.BaseRequest
	    DiskId *string `json:"diskId,omitempty"`							// DiskId 云硬盘ID。
	    SnapshotName *string `json:"snapshotName,omitempty"`				// SnapshotName 快照名称。
	    RetentionTime *string `json:"retentionTime,omitempty"`				// RetentionTime 保留的到期时间。格式为：yyyy-MM-ddTHH:mm:ssZ如果不传，则代表永久保留。指定时间必须在当前时间24小时后。
	}

	type CreateSnapshotResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		SnapshotId *string `json:"snapshotId,omitempty"`					// SnapshotId 创建的快照ID。
	}

	type CreateSnapshotResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *CreateSnapshotResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) CreateSnapshot(snapshotName string, resourceId string) (snapshotId string, err error) {
	funcName := "zecCloudManager:CreateSnapshot:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if snapshotName == "" || resourceId == "" {
		return "", errors.New(ERRORLOG + "error args" + ",snapshotName=" + snapshotName + ",resourceId=" + resourceId)
	}

	request := zec.NewCreateSnapshotRequest()
	request.DiskId = &resourceId
	request.SnapshotName = &snapshotName

	response, err := cm.zecClient.CreateSnapshot(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], snapshotName[%s], resourcevolId[%s]", ERRORLOG, err.Error(), snapshotName, resourceId)
			return "", err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], snapshotName[%s], resourcevolId[%s]", ERRORLOG, err.Error(), snapshotName, resourceId)
		return "", err
	}

	if response.Response.SnapshotId == nil {
		return "", errors.New(ERRORLOG + " nil response. snapshotName=" + snapshotName)
	}

	snapshotId = *response.Response.SnapshotId

	if err = cm.waitDiskSnapStatus(SnapStatusAvailable, snapshotId); err != nil {
		return "", err
	}

	_ = INFOLOG
	return snapshotId, nil
}

/*
Action:	Delete a snapshot

Args: snapshot Id

Return:

zenlayer sdk describe(v0.2.0):

	type DeleteSnapshotsRequest struct {
	    *common.BaseRequest
	    SnapshotIds []string `json:"snapshotIds,omitempty"`					// SnapshotIds 快照ID列表。
	}

	type DeleteSnapshotsResponseParams struct {
		RequestId *string `json:"requestId,omitempty"`
		SnapshotIds []string `json:"snapshotIds,omitempty"`					// SnapshotIds 操作失败的快照ID。
	}

	type DeleteSnapshotsResponse struct {
	    *common.BaseResponse
	    RequestId *string `json:"requestId,omitempty"`
	    Response *DeleteSnapshotsResponseParams `json:"response,omitempty"`
	}
*/
func (cm *zecCloudManager) DeleteSnapshot(snapshotId string) (err error) {
	funcName := "zecCloudManager:DeleteSnapshot:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if snapshotId == "" {
		return errors.New(ERRORLOG + "error args" + ",snapshotId=" + snapshotId)
	}

	request := zec.NewDeleteSnapshotsRequest()
	request.SnapshotIds = append(request.SnapshotIds, snapshotId)

	response, err := cm.zecClient.DeleteSnapshots(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil && err.(*common.ZenlayerCloudSdkError).Code == INVALID_DISK_SNAPSHOT_NOT_FOUND {
			klog.Infof("%s API INFO, snapshotId[%s] Not Found", INFOLOG, snapshotId)
			return nil
		}
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
			return err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
		return err
	}

	if len(response.Response.SnapshotIds) != 0 {
		klog.Errorf("%s API ERROR, delete snap fail, snapshotId[%s]", ERRORLOG, snapshotId)
		return fmt.Errorf("%s API ERROR, DeleteSnapshots response fail snapIds is not nil, req snapshotId[%s], response fail snapshotId[%v]", ERRORLOG, snapshotId, response.Response.SnapshotIds)
	}

	if err = cm.waitDiskSnapStatus(SnapStatusDeleted, snapshotId); err != nil {
		return err
	}

	_ = response
	return nil
}

/*
Action:

Args:

Return:

zenlayer sdk describe(v0.2.0):

		type DescribeSnapshotsRequest struct {
		    *common.BaseRequest
		    SnapshotIds []string `json:"snapshotIds,omitempty"`						// SnapshotIds 根据快照ID列表进行过滤。
		    ZoneId *string `json:"zoneId,omitempty"`								// ZoneId 快照所属的可用区ID。
		    Status *string `json:"status,omitempty"`								// Status 根据快照的状态过滤。
		    DiskIds []string `json:"diskIds,omitempty"`								// DiskIds 按照快照所属的Disk ID列表 过滤。
		    DiskType *string `json:"diskType,omitempty"`							// DiskType 根据快照的云盘类型过滤。
		    SnapshotType *string `json:"snapshotType,omitempty"`					// SnapshotType 根据快照类型过滤。
		    SnapshotName *string `json:"snapshotName,omitempty"`					// SnapshotName 根据快照显示名称过滤。该字段支持模糊搜索。
		    PageSize *int `json:"pageSize,omitempty"`								// PageSize 返回的分页大小。
		    PageNum *int `json:"pageNum,omitempty"`									// PageNum 返回的分页数。
		    ResourceGroupId *string `json:"resourceGroupId,omitempty"`				// ResourceGroupId 根据资源组ID过滤。
		}

		type DescribeSnapshotsResponseParams struct {
			RequestId *string `json:"requestId,omitempty"`
			TotalCount *int `json:"totalCount,omitempty"`							// TotalCount 满足过滤条件的快照总数。
			DataSet []*SnapshotInfo `json:"dataSet,omitempty"`						// DataSet 返回的快照列表数据。
		}

		type SnapshotInfo struct {
			SnapshotId *string `json:"snapshotId,omitempty"`						// SnapshotId 快照唯一ID。
			SnapshotName *string `json:"snapshotName,omitempty"`					// SnapshotName 快照显示名称。
			ZoneId *string `json:"zoneId,omitempty"`								// ZoneId 快照所属的可用区ID。
			Status *string `json:"status,omitempty"`								// Status 快照的状态。
			SnapshotType *string `json:"snapshotType,omitempty"`					// SnapshotType 快照的类型。
			RetentionTime *string `json:"retentionTime,omitempty"`					// RetentionTime 快照的保留到期时间。如果取不到值，说明快照为永久保留。
			DiskId *string `json:"diskId,omitempty"`								// DiskId 云盘ID。
			CreateTime *string `json:"createTime,omitempty"`						// CreateTime 创建时间。RFC3339
			DiskAbility *bool `json:"diskAbility,omitempty"`						// DiskAbility 是否具备创建disk的能力。
			ResourceGroup *ResourceGroupInfo `json:"resourceGroup,omitempty"`		// ResourceGroup 所属的资源组信息。
		}

		type ResourceGroupInfo struct {
	    	ResourceGroupId *string `json:"resourceGroupId,omitempty"`				// ResourceGroupId 资源组ID。
	    	ResourceGroupName *string `json:"resourceGroupName,omitempty"`			// ResourceGroupName 资源组名称。
		}

		type DescribeSnapshotsResponse struct {
		    *common.BaseResponse
		    RequestId *string `json:"requestId,omitempty"`
		    Response *DescribeSnapshotsResponseParams `json:"response,omitempty"`
		}
*/
func (cm *zecCloudManager) FindSnapshot(snapshotId string) (*ZecVolumeSnap, error) {
	funcName := "zecCloudManager:FindSnapshot:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if snapshotId == "" {
		return nil, errors.New(ERRORLOG + "error args" + ",snapshotId=" + snapshotId)
	}

	request := zec.NewDescribeSnapshotsRequest()
	request.SnapshotIds = append(request.SnapshotIds, snapshotId)

	response, err := cm.zecClient.DescribeSnapshots(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
		return nil, err
	}

	if response.Response == nil || response.Response.TotalCount == nil {
		return nil, errors.New(ERRORLOG + " nil response. snapshotId=" + snapshotId)
	}

	if *(response.Response.TotalCount) == 0 {
		return nil, nil
	}
	if *(response.Response.TotalCount) != 1 {
		return nil, errors.New(ERRORLOG + "snapshotId Repeat. snapshotid=" + snapshotId)
	}

	snapinfo := NewZecVolumeSnap()
	if response.Response.DataSet[0] == nil {
		return nil, errors.New(ERRORLOG + " nil response. snapshotId=" + snapshotId)
	}

	if response.Response.DataSet[0].SnapshotId != nil {
		snapinfo.ZecVolumeSnap_Id = *(response.Response.DataSet[0].SnapshotId)
	}
	if response.Response.DataSet[0].SnapshotName != nil {
		snapinfo.ZecVolumeSnap_Name = *(response.Response.DataSet[0].SnapshotName)
	}
	if response.Response.DataSet[0].Status != nil {
		snapinfo.ZecVolumeSnap_status = *(response.Response.DataSet[0].Status)
	}
	if response.Response.DataSet[0].DiskId != nil {
		snapinfo.ZecVolumeSnap_SrcDiskId = *(response.Response.DataSet[0].DiskId)
	}
	if response.Response.DataSet[0].DiskAbility != nil {
		snapinfo.ZecVolumeSnap_DiskAbility = *(response.Response.DataSet[0].DiskAbility)
	}
	if response.Response.DataSet[0].ZoneId != nil {
		snapinfo.ZecVolumeSnap_ZoneId = *(response.Response.DataSet[0].ZoneId)
	}
	if response.Response.DataSet[0].ResourceGroup.ResourceGroupId != nil {
		snapinfo.ZecVolumeSnap_ResourceGroupId = *(response.Response.DataSet[0].ResourceGroup.ResourceGroupId)
	}
	if response.Response.DataSet[0].CreateTime != nil {
		snapinfo.ZecVolumeSnap_CreateTime = *(response.Response.DataSet[0].CreateTime)
	}

	_ = INFOLOG
	return snapinfo, nil
}

func (cm *zecCloudManager) FindSnapshotByName(snapshotName string, srcVolId string, zoneId string, ResourceGroupId string) (*ZecVolumeSnap, error) {
	funcName := "zecCloudManager:FindSnapshotByName:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if snapshotName == "" {
		return nil, errors.New(ERRORLOG + "error args" + ",snapshotName=" + snapshotName)
	}

	sn := "^" + snapshotName + "$"
	var pagesize int = 20
	var pagenum int = 1
	request := zec.NewDescribeSnapshotsRequest()
	request.SnapshotName = &sn
	request.ZoneId = &zoneId
	request.ResourceGroupId = &ResourceGroupId
	request.PageSize = &pagesize
	request.PageNum = &pagenum

	response, err := cm.zecClient.DescribeSnapshots(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Errorf("%s API ERROR, error[%v], snapshotName[%s]", ERRORLOG, err.Error(), snapshotName)
			return nil, err
		}
	} else if err != nil {
		klog.Errorf("%s API ERROR, error[%v], snapshotName[%s]", ERRORLOG, err.Error(), snapshotName)
		return nil, err
	}
	if response.Response == nil || response.Response.TotalCount == nil {
		return nil, errors.New(ERRORLOG + " nil response. snapshotName=" + snapshotName)
	}

	if *(response.Response.TotalCount) == 0 {
		return nil, nil
	}

	//分页处理
	count := *response.Response.TotalCount
	var SnapListAll []*zec.SnapshotInfo = make([]*zec.SnapshotInfo, count)

	//如果需要分页
	if count > *request.PageSize {
		apicallnum := count / *request.PageSize
		if count%*request.PageSize != 0 {
			apicallnum++
		}

		//从第一页开始遍历
		var firstpage int = 1
		request.PageNum = &firstpage
		var end int = 0
		for ra := 0; ra < apicallnum; ra++ {
			response, err := cm.zecClient.DescribeSnapshots(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], snapshotName[%s]", ERRORLOG, err.Error(), snapshotName)
					return nil, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], snapshotName[%s]", ERRORLOG, err.Error(), snapshotName)
				return nil, err
			}
			klog.Infof("%s Traversal once,sum count[%d], pagenum[%d], pagesize[%d], rsp count[%d]", INFOLOG, count, request.PageNum, request.PageSize, len(response.Response.DataSet))
			*request.PageNum++
			for i := 0; i < len(response.Response.DataSet); i++ {
				SnapListAll[i+end] = response.Response.DataSet[i]
			}
			end += len(response.Response.DataSet)
		}
	} else { //不需要分页
		for i := 0; i < count; i++ {
			SnapListAll[i] = response.Response.DataSet[i]
		}
	}

	if len(SnapListAll) != count {
		klog.Errorf("%s API ERROR, Find SnapListAll len[%d] not equal to response.TotalCount[%d]", ERRORLOG, len(SnapListAll), count)
		return nil, fmt.Errorf("%s API ERROR, Find SnapListAll len[%d] not equal to response.TotalCount[%d]", ERRORLOG, len(SnapListAll), count)
	}

	zecSnapInfos := make([]*ZecVolumeSnap, count)
	for c := 0; c < count; c++ {
		snapinfo := NewZecVolumeSnap()
		if SnapListAll[c].SnapshotId != nil {
			snapinfo.ZecVolumeSnap_Id = *SnapListAll[c].SnapshotId
		}
		if SnapListAll[c].SnapshotName != nil {
			snapinfo.ZecVolumeSnap_Name = *SnapListAll[c].SnapshotName
		}
		if SnapListAll[c].Status != nil {
			snapinfo.ZecVolumeSnap_status = *SnapListAll[c].Status
		}
		if SnapListAll[c].DiskId != nil {
			snapinfo.ZecVolumeSnap_SrcDiskId = *SnapListAll[c].DiskId
		}
		if SnapListAll[c].DiskAbility != nil {
			snapinfo.ZecVolumeSnap_DiskAbility = *SnapListAll[c].DiskAbility
		}
		if SnapListAll[c].ZoneId != nil {
			snapinfo.ZecVolumeSnap_ZoneId = *SnapListAll[c].ZoneId
		}
		if SnapListAll[c].ResourceGroup.ResourceGroupId != nil {
			snapinfo.ZecVolumeSnap_ResourceGroupId = *SnapListAll[c].ResourceGroup.ResourceGroupId
		}
		if SnapListAll[c].CreateTime != nil {
			snapinfo.ZecVolumeSnap_CreateTime = *SnapListAll[c].CreateTime
		}

		zecSnapInfos[c] = snapinfo
		//进此分支说明找到了同名的snapshot,如果srcdisk相同则认为是已存在的snapshot，否则返回报错
		if srcVolId == zecSnapInfos[c].ZecVolumeSnap_SrcDiskId {
			klog.Infof("%s find snapname[%s], snapid[%s], srcVolid[%s]", INFOLOG, snapshotName, zecSnapInfos[c].ZecVolumeSnap_Id, srcVolId)
			return zecSnapInfos[c], nil
		} else {
			klog.Errorf("%s find same name snapshot, but srcVolId is not required. snapshot srcvol[%s], req srcvol[%s]", ERRORLOG, zecSnapInfos[c].ZecVolumeSnap_SrcDiskId, srcVolId)
			return zecSnapInfos[c], fmt.Errorf("%s find same name snapshot, but srcVolId is not required. snapshot srcvol[%s], req srcvol[%s]", ERRORLOG, zecSnapInfos[c].ZecVolumeSnap_SrcDiskId, srcVolId)
		}
	}

	return nil, nil
}

/*
Action: wait Apisdk CreateDisk/ReleaseDisk/AttachDisk/DetachDisk Asynchronous operation done

Args: status(string) // needed disk status.    volId(string) //cloud disk ID

Return: error
*/
func (cm *zecCloudManager) waitDiskStatus(status string, volId string) (err error) {

	funcName := "zecCloudManager:waitDiskStatus:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if volId == "" {
		return errors.New(ERRORLOG + "error args" + ".Volid=" + volId)
	}

	request := zec.NewDescribeDisksRequest()
	request.DiskIds = append(request.DiskIds, volId)

	if DiskStatusAvailable == status {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
				return true, err
			}
			if response.Response.TotalCount == nil || response.Response.DataSet[0].DiskStatus == nil {
				return true, errors.New(ERRORLOG + " nil response.volid=" + volId)
			}

			if *response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId Repeat or missing, volid=" + volId)
			}
			if *response.Response.DataSet[0].DiskStatus == DiskStatusAvailable {
				return true, nil
			} else {
				klog.Infof("%s Wait Avaliable 3 sec again.VolId[%s]", INFOLOG, volId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	} else if status == DiskStatusDeleted {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
				return true, err
			}
			if response.Response.TotalCount == nil {
				return true, errors.New(ERRORLOG + " nil response.volid=" + volId)
			}

			if *response.Response.TotalCount == 0 {
				return true, nil
			} else {
				klog.Infof("%s Wait Deleted 3 sec again.VolId[%s]", INFOLOG, volId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err

	} else if status == DiskStatusInUse {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
				return true, err
			}
			if response.Response.TotalCount == nil || response.Response.DataSet[0].DiskStatus == nil {
				return true, errors.New(ERRORLOG + " nil response.volid=" + volId)
			}

			if *response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId repeat or missing, volid=" + volId)
			}

			if *response.Response.DataSet[0].DiskStatus == DiskStatusInUse {
				return true, nil
			} else {
				klog.Infof("%s Wait In_use 3 sec again.VolId[%s]", INFOLOG, volId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	} else if status == DiskStatusStable {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], volId[%s]", ERRORLOG, err.Error(), volId)
				return true, err
			}
			if response.Response.TotalCount == nil || response.Response.DataSet[0].DiskStatus == nil {
				return true, errors.New(ERRORLOG + " nil response.volid=" + volId)
			}

			if *response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId repeat or missing, volid=" + volId)
			}

			if *response.Response.DataSet[0].DiskStatus == DiskStatusAvailable || *response.Response.DataSet[0].DiskStatus == DiskStatusInUse {
				return true, nil
			} else {
				klog.Infof("%s Wait Stable 3 sec again.VolId[%s]", INFOLOG, volId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	}

	return errors.New(ERRORLOG + "error status" + ".VolId=" + volId)
}

/*
Action: wait Apisdk CreateSnapshot/DeleteSnapshots Asynchronous operation done

Args: status(string) // needed snap status.    volId(string) //snap ID

Return: error
*/
func (cm *zecCloudManager) waitDiskSnapStatus(status string, snapshotId string) (err error) {

	funcName := "zecCloudManager:waitDiskSnapStatus:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if snapshotId == "" {
		return errors.New(ERRORLOG + "error args" + ".snapshotId=" + snapshotId)
	}

	request := zec.NewDescribeSnapshotsRequest()
	request.SnapshotIds = append(request.SnapshotIds, snapshotId)

	if SnapStatusAvailable == status {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeSnapshots(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
				return true, err
			}
			if response.Response.TotalCount == nil || response.Response.DataSet[0].Status == nil {
				return true, errors.New(ERRORLOG + " nil response.snapid=" + snapshotId)
			}

			if *response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "snapshotId Repeat or missing, snapshotid=" + snapshotId)
			}
			if *response.Response.DataSet[0].Status == SnapStatusAvailable {
				return true, nil
			} else {
				klog.Infof("%s Wait SnapAvaliable 3 sec again.snapshotId[%s]", INFOLOG, snapshotId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	} else if SnapStatusDeleted == status {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeSnapshots(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
					return true, err
				}
			} else if err != nil {
				klog.Errorf("%s API ERROR, error[%v], snapshotId[%s]", ERRORLOG, err.Error(), snapshotId)
				return true, err
			}
			if response.Response.TotalCount == nil {
				return true, errors.New(ERRORLOG + " nil response.snapid=" + snapshotId)
			}

			if *response.Response.TotalCount == 0 {
				return true, nil
			} else {
				klog.Infof("%s Wait SnapDeleted 3 sec again.snapshotId[%s]", INFOLOG, snapshotId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	}
	return errors.New(ERRORLOG + "error status" + ".snapshotId=" + snapshotId)
}

/*
Action: not support
*/
func (cm *zecCloudManager) CloneVolume() (err error) {
	funcName := "zecCloudManager:CloneVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	_ = ERRORLOG
	_ = INFOLOG
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) FindTag() (err error) {
	funcName := "zecCloudManager:FindTag:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	_ = ERRORLOG
	_ = INFOLOG
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) IsValidTags() bool {
	funcName := "zecCloudManager:IsValidTags:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	_ = ERRORLOG
	_ = INFOLOG
	return false
}

/*
Action: not support
*/
func (cm *zecCloudManager) AttachTags() (err error) {
	funcName := "zecCloudManager:AttachTags:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	_ = ERRORLOG
	_ = INFOLOG
	return nil
}
