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
	zec "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/zec20240401"
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
zenlayer sdk describe(v0.1.27):

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

	config := common.NewConfig().WithTimeout(300).WithAutoRetry(true).WithMaxRetryTime(3)

	cloud_secretKeyId := CLOUDAK
	cloud_secretKeyPassword := CLOUDSK

	cm.zecClient, err = zec.NewClient(config, cloud_secretKeyId, cloud_secretKeyPassword)
	if err != nil || cm.zecClient == nil {
		klog.Error("ERROR:zec.NewClient error")
		return nil, err
	}
	return cm, nil
}

/*
Action: check connect to console.zenlayer.com
*/
func (cm *zecCloudManager) Probe() error {
	funcName := "zecCloudManager:Probe:"
	ERRORLOG := "ERROR:" + funcName + " "
	request := zec.NewDescribeDiskRegionsRequest()

	response, err := cm.zecClient.DescribeDiskRegions(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId)
		return err
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	return err
}

func (cm zecCloudManager) IsController() bool {
	return cm.Iscontroller
}

/*
Action: create cloud disk.

Args: volName string, volSize int, volCategory string, zoneId string, resouceGroupID string

Return: volId string, err error

zenlayer sdk describe(v0.1.27):

	type CreateDisksRequest struct {
		*common.BaseRequest
		ZoneId string `json:"zoneId,omitempty"`
		DiskName string `json:"diskName,omitempty"`
		DiskSize int `json:"diskSize,omitempty"`
		DiskAmount int `json:"diskAmount,omitempty"`
		InstanceId string `json:"instanceId,omitempty"`
		ResourceGroupId string `json:"resourceGroupId,omitempty"`
		DiskCategory string `json:"diskCategory,omitempty"`
	}

	type CreateDisksResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response *CreateDisksResponseParams `json:"response"`
	}

	type CreateDisksResponseParams struct {
		RequestId string `json:"requestId,omitempty"`
		DiskIds []string `json:"diskIds,omitempty"`
		OrderNumber string `json:"orderNumber,omitempty"`
	}
*/
func (cm *zecCloudManager) CreateVolume(volName string, volSize int, volCategory string, zoneId string, resouceGroupID string) (volId string, err error) {

	funcName := "zecCloudManager:CreateVolume:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	if volName == "" {
		return "", errors.New(ERRORLOG + "args error")
	}

	request := zec.NewCreateDisksRequest()
	request.DiskName = volName
	request.DiskSize = volSize
	request.ZoneId = zoneId
	request.DiskAmount = 1
	request.ResourceGroupId = resouceGroupID
	request.DiskCategory = volCategory

	response, err := cm.zecClient.CreateDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".volId=" + volId)
			return "", err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return "", err
	}

	volId = response.Response.DiskIds[0]

	if err = cm.waitDiskStatus(DiskStatusAvaliable, volId); err != nil {
		return "", err
	}

	return volId, nil
}

/*
Action: delete cloud disk.

Args: volId

Return: err error

zenlayer sdk describe(v0.1.27):

	type ReleaseDiskRequest struct {
		*common.BaseRequest
		DiskId string `json:"diskId,omitempty"`
	}

	type ReleaseDiskResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response"`
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
		return errors.New(ERRORLOG + "args error" + ".volId=" + volId)
	}

	request := zec.NewReleaseDiskRequest()
	request.DiskId = volId

	response, err := cm.zecClient.ReleaseDisk(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil && err.(*common.ZenlayerCloudSdkError).Code == OPERATION_FAILED_RESOURCE_NOT_FOUND {
			klog.Info(INFOLOG + "API INFO" + err.Error() + ".ReqID=" + response.RequestId + ".volId=" + volId + " Not Found.")
			return nil
		} else {
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".volId=" + volId)
			return err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusDeleted, volId); err != nil {
		return err
	}

	return nil
}

/*
Action: attach cloud disk to Vm.

Args: volId, vmid

Return: err error

zenlayer sdk describe(v0.1.27):

	type AttachDisksRequest struct {
		*common.BaseRequest
		DiskIds []string `json:"diskIds,omitempty"`
		InstanceId string `json:"instanceId,omitempty"`
	}

	type AttachDisksResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response"`
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
	request.InstanceId = vmId

	response, err := cm.zecClient.AttachDisks(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".VolId=" + volId + ".Vmid=" + vmId)
			return err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusInUse, volId); err != nil {
		return err
	}

	return nil
}

/*
Action: detach cloud disk.

Args: volId

Return: err error

zenlayer sdk describe(v0.1.27):

	type DetachDisksRequest struct {
		*common.BaseRequest
		DiskIds []string `json:"diskIds,omitempty"`
		InstanceCheckFlag *bool `json:"instanceCheckFlag,omitempty"`
	}

	type DetachDisksResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response struct {
			RequestId string `json:"requestId,omitempty"`
		} `json:"response"`
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
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".volId=" + volId)
			return err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusAvaliable, volId); err != nil {
		return err
	}

	return nil
}

/*
Action: Resize cloud disk.

Args: volId, newSize

Return: err error

zenlayer sdk describe(v0.1.27):

	type ResizeDiskRequest struct {
		*common.BaseRequest
		DiskId *string `json:"diskId,omitempty"`
		DiskSize *int `json:"diskSize,omitempty"`	//only support expand, max size 32768GB
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
			if response.RequestId != nil {
				klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + *response.RequestId + ".volId=" + volId)
			} else {
				klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".volId=" + volId)
			}
			return err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	if err = cm.waitDiskStatus(DiskStatusStable, volId); err != nil {
		return err
	}

	return nil
}

/*
Action: get disk info by volid.

Args: volId

Return: zecVolInfo *ZecVolume, err error

zenlayer sdk describe(v0.1.27):

type DescribeDisksRequest struct {

		*common.BaseRequest
		DiskIds []string `json:"diskIds,omitempty"`			//云硬盘ID集合
		DiskName string `json:"diskName,omitempty"`			//云硬盘名称
		DiskStatus string `json:"diskStatus,omitempty"`		//云硬盘状态
		DiskType string `json:"diskType,omitempty"`			//云硬盘类型 SYSTEM系统盘，DATA数据盘
		DiskCategory string `json:"diskCategory,omitempty"`	//云硬盘种类 “Basic NVMe SSD” or “Standard NVMe SSD”
		InstanceId string `json:"instanceId,omitempty"`		//VM实例ID
		ZoneId string `json:"zoneId,omitempty"`				//可用区ID
		PageSize int `json:"pageSize,omitempty"`			//返回的分页大小 默认为20，最大为1000
		PageNum int `json:"pageNum,omitempty"`				//返回的分页数 默认是1
	}

	type DescribeDisksResponse struct {
			*common.BaseResponse
			RequestId string `json:"requestId,omitempty"`
			Response *DescribeDisksResponseParams `json:"response"`
		}

	type DescribeDisksResponseParams struct {
			RequestId string `json:"requestId,omitempty"`
			DataSet []*DiskInfo `json:"dataSet,omitempty"`
			TotalCount int `json:"totalCount,omitempty"`
		}

	type DiskInfo struct {
		DiskName string `json:"diskName,omitempty"`
		ZoneId string `json:"zoneId,omitempty"`
		DiskType string `json:"diskType,omitempty"`
		Portable bool `json:"portable,omitempty"`
		DiskCategory string `json:"diskCategory,omitempty"`
		DiskSize int `json:"diskSize,omitempty"`
		DiskStatus string `json:"diskStatus,omitempty"`
		InstanceId string `json:"instanceId,omitempty"`
		InstanceName string `json:"instanceName,omitempty"`
		ChargeType string `json:"chargeType,omitempty"`
		CreateTime string `json:"createTime,omitempty"`
		ExpiredTime string `json:"expiredTime,omitempty"`
		Period int `json:"period,omitempty"`
		ResourceGroupId string `json:"resourceGroupId,omitempty"`
		ResourceGroupName string `json:"resourceGroupName,omitempty"`
		Serial string `json:"serial,omitempty"`
	}
*/
func (cm *zecCloudManager) FindVolume(volId string) (zecVolInfo *ZecVolume, err error) {

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
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".volid=" + volId)
			return nil, err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return nil, err
	}

	if response.Response.TotalCount == 0 {
		return nil, nil
	}
	if response.Response.TotalCount != 1 {
		return nil, errors.New(ERRORLOG + "Volid Repeat" + volId)
	}

	zecVolInfo = NewZecVolume()
	zecVolInfo.ZecVolume_Id = response.Response.DataSet[0].DiskId
	zecVolInfo.ZecVolume_Name = response.Response.DataSet[0].DiskName
	zecVolInfo.ZecVolume_Size = response.Response.DataSet[0].DiskSize
	zecVolInfo.ZecVolume_Type = int(driver.StringToType(response.Response.DataSet[0].DiskCategory))
	zecVolInfo.ZecVolume_Status = response.Response.DataSet[0].DiskStatus
	zecVolInfo.ZecVolume_Zone = response.Response.DataSet[0].ZoneId
	zecVolInfo.ZecVolume_InstanceId = response.Response.DataSet[0].InstanceId
	zecVolInfo.ZecVolume_Serial = response.Response.DataSet[0].Serial

	return zecVolInfo, nil
}

func (cm *zecCloudManager) FindVolumeByName(volName string, zoneID string, sizeGB int, diskType string) (zecVolInfo *ZecVolume, err error) {

	funcName := "zecCloudManager:FindVolumeByName:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "
	INFOLOG := "INFO:" + funcName + hash + " "

	if volName == "" {
		return nil, errors.New(ERRORLOG + "args error")
	}

	request := zec.NewDescribeDisksRequest()
	request.DiskName = "^" + volName + "$"

	request.PageSize = 20
	request.PageNum = 1
	response, err := cm.zecClient.DescribeDisks(request) //第一个请求用来获取TotalCount，如果不需要分页就直接用此数据，如果需要分页再发请求
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId)
			return nil, err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return nil, err
	}

	if response.Response.TotalCount == 0 {
		return nil, nil
	}

	//用于保存所有数据
	count := response.Response.TotalCount
	var DiskListAll []*zec.DiskInfo = make([]*zec.DiskInfo, count)

	//如果需要分页
	if count > request.PageSize {

		apicallnum := count / request.PageSize
		if count%request.PageSize != 0 {
			apicallnum++
		}

		//从第一页开始遍历
		request.PageNum = 1
		var end int = 0
		for ra := 0; ra < apicallnum; ra++ {
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId)
					return nil, err
				}
			} else if err != nil {
				klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
				return nil, err
			}
			klog.Infof("%s Traversal once,sum count=%d, pagenum=%d, pagesize=%d, rsp count=%d", INFOLOG, count, request.PageNum, request.PageSize, len(response.Response.DataSet))
			request.PageNum++
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

	zecVolInfos := make([]*ZecVolume, count)
	for c := 0; c < count; c++ {
		zecVolInfos[c] = &ZecVolume{
			ZecVolume_Id:         DiskListAll[c].DiskId,
			ZecVolume_Name:       DiskListAll[c].DiskName,
			ZecVolume_Size:       DiskListAll[c].DiskSize,
			ZecVolume_Type:       int(driver.StringToType(DiskListAll[c].DiskCategory)),
			ZecVolume_Zone:       DiskListAll[c].ZoneId,
			ZecVolume_InstanceId: DiskListAll[c].InstanceId,
			ZecVolume_Serial:     DiskListAll[c].Serial,
		}

		if zoneID == zecVolInfos[c].ZecVolume_Zone {
			if zecVolInfos[c].ZecVolume_Type == int(driver.StringToType(diskType)) {
				if zecVolInfos[c].ZecVolume_Size != sizeGB {
					//csi规定创建不允许创建同名不同大小的卷
					klog.Errorf("%s find different size but same volname volumeid=%s, volumename=%s, existvol size=%d, request new size=%d", ERRORLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name, zecVolInfos[c].ZecVolume_Size, sizeGB)
					return zecVolInfos[c], fmt.Errorf("%s find different size but same volname volumeid=%s, volumename=%s, existvol size=%d, request new size=%d", ERRORLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name, zecVolInfos[c].ZecVolume_Size, sizeGB)
				}
				if zecVolInfos[c].ZecVolume_Size == sizeGB {
					klog.Infof("%s find volumeid=%s, volumename=%s, serial=%s, size=%d, zone=%s.", INFOLOG, zecVolInfos[c].ZecVolume_Id, zecVolInfos[c].ZecVolume_Name,
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

zenlayer sdk describe(v0.1.27):

	type DescribeInstancesRequest struct {
		*common.BaseRequest
		InstanceIds []string `json:"instanceIds,omitempty"`
		ZoneId string `json:"zoneId,omitempty"`
		ImageId string `json:"imageId,omitempty"`
		Status string `json:"status,omitempty"`
		Name string `json:"name,omitempty"`
		Ipv4Address string `json:"ipv4Address,omitempty"`
		Ipv6Address string `json:"ipv6Address,omitempty"`
		PageSize int `json:"pageSize,omitempty"`
		PageNum int `json:"pageNum,omitempty"`
	}

	type DescribeInstancesResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response *DescribeInstancesResponseParams `json:"response"`
	}

	type DescribeInstancesResponseParams struct {
		RequestId string `json:"requestId,omitempty"`
		DataSet []*InstanceInfo `json:"dataSet,omitempty"`
		TotalCount int `json:"totalCount,omitempty"`
	}

	type InstanceInfo struct {
		InstanceId string `json:"instanceId,omitempty"`
		InstanceName string `json:"instanceName,omitempty"`
		ZoneId string `json:"zoneId,omitempty"`
		Cpu int `json:"cpu,omitempty"`
		Memory int `json:"memory,omitempty"`
		ImageId string `json:"imageId,omitempty"`
		ImageName string `json:"imageName,omitempty"`
		InstanceType *string `json:"instanceType,omitempty"`
		Status string `json:"status,omitempty"`
		SystemDisk *SystemDisk `json:"systemDisk,omitempty"`
		DataDisks []*DataDisk `json:"dataDisks,omitempty"`
		PublicIpAddresses []string `json:"publicIpAddresses,omitempty"`
		PrivateIpAddresses []string `json:"privateIpAddresses,omitempty"`
		KeyId string `json:"keyId,omitempty"`
		CreateTime string `json:"createTime,omitempty"`
		ExpiredTime string `json:"expiredTime,omitempty"`
		ResourceGroupId string `json:"resourceGroupId,omitempty"`
		ResourceGroupName string `json:"resourceGroupName,omitempty"`
	}
*/
func (cm *zecCloudManager) FindInstance(vmid string) (zecVmInfo *ZecVm, err error) {

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
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".VmId=" + vmid)
			return nil, err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return nil, err
	}

	if response.Response.TotalCount == 0 {
		return nil, nil
	}
	if response.Response.TotalCount != 1 {
		return nil, errors.New(ERRORLOG + "Vmid Repeat" + vmid)
	}

	zecVmInfo = NewZecVm()
	zecVmInfo.ZecVm_Id = response.Response.DataSet[0].InstanceId
	zecVmInfo.ZecVm_Status = response.Response.DataSet[0].Status
	zecVmInfo.ZecVm_Type = driver.BasicVmType.Int()
	zecVmInfo.ZecVm_Zone = response.Response.DataSet[0].ZoneId

	return zecVmInfo, nil
}

/*
Action: get zec zone info.

Args: zoneid

Return: error

zenlayer sdk describe(v0.1.27):

	type DescribeZonesRequest struct {
		*common.BaseRequest
		ZoneIds []string `json:"zoneIds,omitempty"`
	}

	type DescribeZonesResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response *DescribeZonesResponseParams `json:"response"`
	}

	type DescribeZonesResponseParams struct {
		RequestId string `json:"requestId,omitempty"`
		ZoneSet []*ZoneInfo `json:"zoneSet,omitempty"`
	}

	type ZoneInfo struct {
		// Zone ID. For example, SEL-A.
		ZoneId string `json:"zoneId,omitempty"`
		// Zone name.
		ZoneName string `json:"zoneName,omitempty"`
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
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId)
			return err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return err
	}

	return nil
}

/*
Action: get cloud disk support zones.

Args: nil

Return: zones []string.     error

zenlayer sdk describe(v0.1.27):

	type DescribeDiskRegionsRequest struct {
		*common.BaseRequest
		ChargeType string `json:"chargeType,omitempty"`
	}

	type DescribeDiskRegionsResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response *DescribeDiskRegionsResponseParams `json:"response"`
	}

	type DescribeDiskRegionsResponseParams struct {
		RequestId string `json:"requestId,omitempty"`
		RegionIds []string `json:"regionIds,omitempty"`
	}
*/
func (cm *zecCloudManager) GetZoneList() (zones []string, err error) {
	funcName := "zecCloudManager:GetZoneList:"
	funcInfo, hash := csicommon.EntryFunction(funcName)
	klog.Info(funcInfo)
	defer klog.Info(csicommon.ExitFunction(funcName, hash))
	ERRORLOG := "ERROR:" + funcName + hash + " "

	request := zec.NewDescribeDiskRegionsRequest()

	response, err := cm.zecClient.DescribeDiskRegions(request)
	if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
		if err != nil {
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId)
			return nil, err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return nil, err
	}

	zones = response.Response.RegionIds
	return zones, nil
}

/*
Action: check zec vm status, only support running vm.

Args: vmid(string)

Return: error

zenlayer sdk describe(v0.1.27):

	type DescribeInstancesStatusResponse struct {
		*common.BaseResponse
		RequestId string `json:"requestId,omitempty"`
		Response *DescribeInstancesStatusResponseParams `json:"response"`
	}

	type DescribeInstancesStatusResponseParams struct {
		RequestId string `json:"requestId,omitempty"`
		DataSet []*InstanceStatus `json:"dataSet,omitempty"`
		TotalCount int `json:"totalCount,omitempty"`
	}

	type InstanceStatus struct {
		InstanceId string `json:"instanceId,omitempty"`
		InstanceStatus string `json:"instanceStatus,omitempty"`
	}

	type DescribeInstancesStatusRequest struct {
		*common.BaseRequest
		InstanceIds []string `json:"instanceIds,omitempty"`
		PageSize int `json:"pageSize,omitempty"`
		PageNum int `json:"pageNum,omitempty"`
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
			klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ",Vmid=" + vmId)
			return false, "", err
		}
	} else if err != nil {
		klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
		return false, "", err
	}

	if response.Response.TotalCount == 0 {
		return false, "", nil
	}

	if response.Response.TotalCount != 1 {
		return false, "", errors.New(ERRORLOG + "VmId repeat" + vmId)
	}

	return true, response.Response.DataSet[0].InstanceStatus, nil
}

/*
Action: wait Apisdk CreateDisk/ReleaseDisk/AttachDisk/DetachDisk Asynchronous operation done

Args: status(string) // needed adisk status.    volId(string) //cloud disk ID

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

	if DiskStatusAvaliable == status {
		job := func() (stop bool, err error) {
			//ignore PageNum
			response, err := cm.zecClient.DescribeDisks(request)
			if _, ok := err.(*common.ZenlayerCloudSdkError); ok {
				if err != nil {
					klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".Volid=" + volId)
					return true, err
				}
			} else if err != nil {
				klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
				return true, err
			}

			if response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId Repeat or missing" + volId)
			}
			if response.Response.DataSet[0].DiskStatus == DiskStatusAvaliable {
				return true, nil
			} else {
				klog.Info(INFOLOG+" Wait Avaliable 3 sec again.VolId="+volId)
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
					klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".VolId=" + volId)
					return true, err
				}
			} else if err != nil {
				klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
				return true, err
			}

			if response.Response.TotalCount == 0 {
				return true, nil
			} else {
				klog.Info(INFOLOG+" Wait Deleted 3 sec again.VolId=", volId)
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
					klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".VolId=" + volId)
					return true, err
				}
			} else if err != nil {
				klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
				return true, err
			}

			if response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId repeat or missing" + volId)
			}

			if response.Response.DataSet[0].DiskStatus == DiskStatusInUse {
				return true, nil
			} else {
				klog.Info(INFOLOG+" Wait In_use 3 sec again.VolId=", volId)
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
					klog.Error(ERRORLOG + "API ERROR" + err.Error() + ".ReqID=" + response.RequestId + ".VolId=" + volId)
					return true, err
				}
			} else if err != nil {
				klog.Error(ERRORLOG + "API ERROR Unexpected " + err.Error())
				return true, err
			}

			if response.Response.TotalCount != 1 {
				return true, errors.New(ERRORLOG + "VolId repeat or missing" + volId)
			}

			if response.Response.DataSet[0].DiskStatus == DiskStatusAvaliable || response.Response.DataSet[0].DiskStatus == DiskStatusInUse {
				return true, nil
			} else {
				klog.Info(INFOLOG+" Wait Stable 3 sec again.VolId=", volId)
				return false, nil
			}
		}

		err = WaitFor(job)
		return err
	}

	return errors.New(ERRORLOG + "error status" + ".VolId=" + volId)
}

/*
Action: not support
*/
func (cm *zecCloudManager) CloneVolume() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) FindTag() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) IsValidTags() bool {
	return true
}

/*
Action: not support
*/
func (cm *zecCloudManager) AttachTags() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) FindSnapshot() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) FindSnapshotByName() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) CreateSnapshot() (err error) {
	return nil
}

/*
Action: not support
*/
func (sm *zecCloudManager) DeleteSnapshot() (err error) {
	return nil
}

/*
Action: not support
*/
func (cm *zecCloudManager) CreateVolumeFromSnapshot() (err error) {
	return nil
}
