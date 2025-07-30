package zec

import "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/common"


type CreateVpcResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateVpcResponseParams `json:"response"`
}

type CreateVpcResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	VpcId string `json:"vpcId,omitempty"`
}

type CreateVpcRequest struct {
	*common.BaseRequest

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	Mtu int `json:"mtu,omitempty"`

	EnablePriIpv6 bool `json:"enablePriIpv6,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`
}

type DeleteVpcResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DeleteVpcRequest struct {
	*common.BaseRequest

	VpcId string `json:"vpcId,omitempty"`
}

type ModifyVpcsAttributeResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyVpcsAttributeRequest struct {
	*common.BaseRequest

	VpcIds []string `json:"vpcIds,omitempty"`

	Name string `json:"name,omitempty"`
}

type DescribeVpcsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeVpcsResponseParams `json:"response"`
}

type DescribeVpcsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*VpcInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type VpcInfo struct {
	VpcId           string `json:"vpcId,omitempty"`
	Name            string `json:"name,omitempty"`
	CidrBlock       string `json:"cidrBlock,omitempty"`
	Ipv6CidrBlock   string `json:"ipv6CidrBlock,omitempty"`
	Mtu             int    `json:"mtu,omitempty"`
	IsDefault       bool   `json:"isDefault,omitempty"`
	CreateTime      string `json:"createTime,omitempty"`
	UsageIpv4Count  int    `json:"usageIpv4Count,omitempty"`
	UsageIpv6Count  int    `json:"usageIpv6Count,omitempty"`
	SecurityGroupId string `json:"securityGroupId,omitempty"`
}

type DescribeVpcsRequest struct {
	*common.BaseRequest

	VpcIds []string `json:"vpcIds,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribeSubnetRegionsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeSubnetRegionsResponseParams `json:"response"`
}

type DescribeSubnetRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionSet []*RegionInfo `json:"regionSet,omitempty"`
}

type RegionInfo struct {
	RegionId      string `json:"regionId,omitempty"`
	RegionName    string `json:"regionName,omitempty"`
	RegionTitle   string `json:"regionTitle,omitempty"`
	EnablePubIpv6 bool   `json:"enablePubIpv6,omitempty"`
}

type DescribeSubnetRegionsRequest struct {
	*common.BaseRequest

	RegionIds []string `json:"regionIds,omitempty"`
}

type CreateSubnetResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateSubnetResponseParams `json:"response"`
}

type CreateSubnetResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	SubnetId string `json:"subnetId,omitempty"`
}

type CreateSubnetRequest struct {
	*common.BaseRequest

	VpcId string `json:"vpcId,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	StackType string `json:"stackType,omitempty"`

	Ipv6Type string `json:"ipv6Type,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`
}

type DeleteSubnetResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DeleteSubnetRequest struct {
	*common.BaseRequest

	SubnetId string `json:"subnetId,omitempty"`
}

type ModifySubnetsAttributeResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifySubnetsAttributeRequest struct {
	*common.BaseRequest

	SubnetIds []string `json:"subnetIds,omitempty"`

	Name string `json:"name,omitempty"`
}

type DescribeSubnetsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeSubnetsResponseParams `json:"response"`
}

type DescribeSubnetsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*SubnetInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type SubnetInfo struct {
	SubnetId       string `json:"subnetId,omitempty"`
	RegionId       string `json:"regionId,omitempty"`
	Name           string `json:"name,omitempty"`
	CidrBlock      string `json:"cidrBlock,omitempty"`
	Ipv6CidrBlock  string `json:"ipv6CidrBlock,omitempty"`
	StackType      string `json:"stackType,omitempty"`
	Ipv6Type       string `json:"ipv6Type,omitempty"`
	VpcId          string `json:"vpcId,omitempty"`
	VpcName        string `json:"vpcName,omitempty"`
	UsageIpv4Count int    `json:"usageIpv4Count,omitempty"`
	UsageIpv6Count int    `json:"usageIpv6Count,omitempty"`
	CreateTime     string `json:"createTime,omitempty"`
	IsDefault      bool   `json:"isDefault,omitempty"`
}

type DescribeSubnetsRequest struct {
	*common.BaseRequest

	SubnetIds []string `json:"subnetIds,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type ModifySubnetStackTypeResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *ModifySubnetStackTypeResponseParams `json:"response"`
}

type ModifySubnetStackTypeResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Ipv6CidrBlock string `json:"ipv6CidrBlock,omitempty"`
}

type ModifySubnetStackTypeRequest struct {
	*common.BaseRequest

	SubnetId string `json:"subnetId,omitempty"`

	StackType string `json:"stackType,omitempty"`

	Ipv6Type string `json:"ipv6Type,omitempty"`
}

type CreateRouteResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateRouteResponseParams `json:"response"`
}

type CreateRouteResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	RouteId string `json:"routeId,omitempty"`
}

type CreateRouteRequest struct {
	*common.BaseRequest

	VpcId string `json:"vpcId,omitempty"`

	IpVersion string `json:"ipVersion,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	Priority int `json:"priority,omitempty"`

	NextHotId string `json:"nextHotId,omitempty"`

	Name string `json:"name,omitempty"`
}

type DeleteRouteResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DeleteRouteRequest struct {
	*common.BaseRequest

	VpcId string `json:"vpcId,omitempty"`

	RouteId string `json:"routeId,omitempty"`
}

type DescribeRoutesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeRoutesResponseParams `json:"response"`
}

type DescribeRoutesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*RouteInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type RouteInfo struct {
	RouteId string `json:"routeId,omitempty"`

	Name string `json:"name,omitempty"`

	VpcId string `json:"vpcId,omitempty"`

	VpcName string `json:"vpcName,omitempty"`

	IpVersion string `json:"ipVersion,omitempty"`

	Type string `json:"type,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	Priority int `json:"priority,omitempty"`

	NextHopId string `json:"nextHopId,omitempty"`

	NextHopName string `json:"nextHopName,omitempty"`

	NextHopType string `json:"nextHopType,omitempty"`

	CreateTime string `json:"createTime,omitempty"`
}

type DescribeRoutesRequest struct {
	*common.BaseRequest

	RouteIds []string `json:"routeIds,omitempty"`

	IpVersion string `json:"ipVersion,omitempty"`

	Type string `json:"type,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

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

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
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

	NicNetworkType string `json:"nicNetworkType,omitempty"`

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

type DescribeInstancesStatusResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
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

type ModifyInstancesResourceGroupResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyInstancesResourceGroupRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`
}

type ModifyInstancesAttributeResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyInstancesAttributeRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`

	InstanceName string `json:"instanceName,omitempty"`
}

type RebootInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *RebootInstancesResponseParams `json:"response"`
}

type RebootInstancesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type RebootInstancesRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type ResetInstanceResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ResetInstanceRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`

	ImageId string `json:"imageId,omitempty"`

	Password string `json:"password,omitempty"`

	KeyId string `json:"keyId,omitempty"`

	Timezone string `json:"timezone,omitempty"`

	EnableAgent bool `json:"enableAgent,omitempty"`
}

type ResetInstancePasswordResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ResetInstancePasswordRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`

	Password string `json:"password,omitempty"`
}

type StartInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *StartInstancesResponseParams `json:"response"`
}

type StartInstancesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type StartInstancesRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type StopInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *StopInstancesResponseParams `json:"response"`
}

type StopInstancesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type StopInstancesRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`

	ForceShutdown bool `json:"forceShutdown,omitempty"`
}

////////////

type DescribeCidrsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeCidrsResponseParams `json:"response"`
}

type DescribeCidrsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*CidrInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type CidrInfo struct {
	CidrId string `json:"cidrId,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`

	UsedCount int `json:"usedCount,omitempty"`

	Source string `json:"source,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	Netmask int `json:"netmask,omitempty"`

	PoolId string `json:"poolId,omitempty"`

	CreateTime string `json:"createTime,omitempty"`

	ExpiredTime string `json:"expiredTime,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	ResourceGroupName string `json:"resourceGroupName,omitempty"`
}

type DescribeCidrsRequest struct {
	*common.BaseRequest

	CidrIds []string `json:"cidrIds,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Name string `json:"name,omitempty"`

	CidrBlock string `json:"cidrBlock,omitempty"`

	Source string `json:"source,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribePoolsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribePoolsResponseParams `json:"response"`
}

type DescribePoolsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*PoolInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type PoolInfo struct {
	PoolId string `json:"poolId,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Name string `json:"name,omitempty"`

	CreateTime string `json:"createTime,omitempty"`
}

type DescribePoolsRequest struct {
	*common.BaseRequest

	PoolIds []string `json:"poolIds,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Name string `json:"name,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribeCidrRegionsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeCidrRegionsResponseParams `json:"response"`
}

type DescribeCidrRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionIds []string `json:"regionIds,omitempty"`
}

type DescribeCidrRegionsRequest struct {
	*common.BaseRequest
}

type ChargePrepaid struct {
	// Period of subscription.
	// Unit: month.
	Period int `json:"period,omitempty"`
}

type Price struct {
	// Discount.
	// For example, 80.0 means 20% off.
	Discount *float64 `json:"discount,omitempty"`

	// Discount price of prepaid resources.
	// Only used in subscription model. For pay-as-you-go model, the value is empty.
	DiscountPrice *float64 `json:"discountPrice,omitempty"`

	// Original price of prepaid resources.
	// Only used in subscription model. For pay-as-you-go model, the value is empty.
	OriginalPrice *float64 `json:"originalPrice,omitempty"`

	// Original unit price of postpaid resources.
	// Only used in pay-as-you-go model. For tiered billing, the value is empty.
	UnitPrice *float64 `json:"unitPrice,omitempty"`

	// Discount unit price of postpaid resources.
	// Only used in pay-as-you-go model. For tiered billing, the value is empty.
	DiscountUnitPrice *float64 `json:"discountUnitPrice,omitempty"`

	// Unit of postpaid billing.
	// Only used in pay-as-you-go model.
	// Value range:
	// HOUR: you will be billed by hour.
	// DAY: you will be billed by day.
	// MONTH: you will be billed by month. For example, the burstable 95th pricing model.
	ChargeUnit *string `json:"chargeUnit,omitempty"`

	ExcessUnitPrice *float64 `json:"excessUnitPrice,omitempty"`

	ExcessDiscountUnitPrice *float64 `json:"excessDiscountUnitPrice,omitempty"`

	ExcessAmountUnit *string `json:"excessAmountUnit,omitempty"`

	// Tiered price of postpaid billing.
	// Only used in pay-as-you-go model. If it is not tiered price, the value is empty.
	StepPrices []*StepPrice `json:"stepPrices,omitempty"`
}

type StepPrice struct {
	// First price range of a tiered price.
	StepStart *float64 `json:"stepStart,omitempty"`

	// Last price range of a tiered price.
	StepEnd *float64 `json:"stepEnd,omitempty"`

	// Original unit price of current price range.
	// Only used in pay-as-you-go model.
	UnitPrice *float64 `json:"unitPrice,omitempty"`

	// Discount unit price of current price range.
	// Only used in pay-as-you-go model.
	DiscountUnitPrice *float64 `json:"discountUnitPrice,omitempty"`
}

type DescribeCidrPriceResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeCidrPriceResponseParams `json:"response"`
}

type DescribeCidrPriceResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	CidrPrice *Price `json:"cidrPrice,omitempty"`
}

type DescribeCidrPriceRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	Netmask *NetmaskInfo `json:"netmask,omitempty"`
}

type NetmaskInfo struct {
	// Period of subscription.
	// Unit: month.
	Netmask int `json:"netmask,omitempty"`

	Amount int `json:"amount,omitempty"`
}

type CreateCidrResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateCidrResponseParams `json:"response"`
}

type CreateCidrResponseParams struct {
	RequestId string `json:"requestId,omitempty"`
}

type CreateCidrRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	Netmask *NetmaskInfo `json:"netmask,omitempty"`
}

type DeleteCidrResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`
}

type DeleteCidrRequest struct {
	*common.BaseRequest

	CidrId string `json:"cidrId,omitempty"`
}

type RenewCidrResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`
}

type RenewCidrRequest struct {
	*common.BaseRequest

	CidrId string `json:"cidrId,omitempty"`
}

type AvailableLanIpRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`
}

type AvailableLanIpResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *AvailableLanIpResponseParams `json:"response"`
}

type AvailableLanIpResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	LanIps []*PrivateIpInfo `json:"lanIps,omitempty"`
}

type PrivateIpInfo struct {
	LanIp        string `json:"lanIp,omitempty"`
	NicId        string `json:"nicId,omitempty"`
	NicName      string `json:"nicName,omitempty"`
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
}

type DescribeDiskRegionsRequest struct {
	*common.BaseRequest

	ChargeType string `json:"chargeType,omitempty"`
}

type DescribeDiskRegionsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeDiskRegionsResponseParams `json:"response"`
}

type DescribeDiskRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionIds []string `json:"regionIds,omitempty"`
}

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

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateDisksResponseParams `json:"response"`
}

type CreateDisksResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DiskIds []string `json:"diskIds,omitempty"`

	OrderNumber string `json:"orderNumber,omitempty"`
}

type DescribeDisksRequest struct {
	*common.BaseRequest

	DiskIds []string `json:"diskIds,omitempty"`

	DiskName string `json:"diskName,omitempty"`

	DiskStatus string `json:"diskStatus,omitempty"`

	DiskType string `json:"diskType,omitempty"`

	DiskCategory string `json:"diskCategory,omitempty"`

	InstanceId string `json:"instanceId,omitempty"`

	ZoneId string `json:"zoneId,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribeDisksResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeDisksResponseParams `json:"response"`
}

type DescribeDisksResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*DiskInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type DiskInfo struct {
	DiskId string `json:"diskId,omitempty"`

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

type ChangeDisksAttachRequest struct {
	*common.BaseRequest

	DiskIds []string `json:"diskIds,omitempty"`

	InstanceId string `json:"instanceId,omitempty"`
}

type ChangeDisksAttachResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

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


type ResizeDiskRequest struct {
	*common.BaseRequest

	// DiskId 云硬盘ID。通过DescribeDisks接口查询。
	DiskId *string `json:"diskId,omitempty"`

	// DiskSize 云硬盘扩容后的大小。单位GB。必须大于当前云硬盘大小。云盘最大限制为32768GB。
	DiskSize *int `json:"diskSize,omitempty"`

}

type ResizeDiskResponse struct {
	*common.BaseResponse

	RequestId *string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type ModifyDisksAttributesRequest struct {
	*common.BaseRequest

	DiskIds []string `json:"diskIds,omitempty"`

	DiskName string `json:"diskName,omitempty"`
}

type ModifyDisksAttributesResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type InquiryPriceCreateDisksRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	DiskSize int `json:"diskSize,omitempty"`

	DiskAmount int `json:"diskAmount,omitempty"`

	DiskCategory string `json:"diskCategory,omitempty"`
}

type InquiryPriceCreateDisksResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *InquiryPriceCreateDisksResponseParams `json:"response"`
}

type InquiryPriceCreateDisksResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataDiskPrice *Price `json:"dataDiskPrice,omitempty"`
}

type TerminateDiskRequest struct {
	*common.BaseRequest

	DiskId string `json:"diskId,omitempty"`
}

type TerminateDiskResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

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

type RenewDiskRequest struct {
	*common.BaseRequest

	DiskId string `json:"diskId,omitempty"`
}

type RenewDiskResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyDisksResourceGroupRequest struct {
	*common.BaseRequest

	DiskIds []string `json:"diskIds,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`
}

type ModifyDisksResourceGroupResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DescribeDiskCategoryRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	DiskCategory string `json:"diskCategory,omitempty"`
}

type DescribeDiskCategoryResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeDiskCategoryResponseParams `json:"response"`
}

type DescribeDiskCategoryResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	CategoryZoneSet []*DiskCategory `json:"categoryZoneSet,omitempty"`
}

type DiskCategory struct {
	ZoneId string `json:"zoneId,omitempty"`

	CategorySet []string `json:"categorySet,omitempty"`
}

type DescribeSecurityGroupsRequest struct {
	*common.BaseRequest

	SecurityGroupIds []string `json:"securityGroupIds,omitempty"`

	SecurityGroupName string `json:"securityGroupName,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribeSecurityGroupsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeSecurityGroupsResponseParams `json:"response"`
}

type DescribeSecurityGroupsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*SecurityGroupInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type SecurityGroupInfo struct {
	SecurityGroupId   string   `json:"securityGroupId,omitempty"`
	SecurityGroupName string   `json:"securityGroupName,omitempty"`
	Scope             string   `json:"scope,omitempty"`
	CreateTime        string   `json:"createTime,omitempty"`
	VpcIds            []string `json:"vpcIds,omitempty"`
	IsDefault         bool     `json:"isDefault,omitempty"`
}

type ModifySecurityGroupsAttributeRequest struct {
	*common.BaseRequest

	SecurityGroupIds []string `json:"securityGroupIds,omitempty"`

	SecurityGroupName string `json:"securityGroupName,omitempty"`
}

type ModifySecurityGroupsAttributeResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DescribeSecurityGroupRuleRequest struct {
	*common.BaseRequest

	SecurityGroupId string `json:"securityGroupId,omitempty"`
}

type DescribeSecurityGroupRuleResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeSecurityGroupRuleResponseParams `json:"response"`
}

type DescribeSecurityGroupRuleResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	IngressRuleList []*RuleInfo `json:"ingressRuleList,omitempty"`

	EgressRuleList []*RuleInfo `json:"egressRuleList,omitempty"`
}

type RuleInfo struct {
	Direction  string `json:"direction,omitempty"`
	Policy     string `json:"policy,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	IpProtocol string `json:"ipProtocol,omitempty"`
	PortRange  string `json:"portRange,omitempty"`
	CidrIp     string `json:"cidrIp,omitempty"`
	Desc       string `json:"desc,omitempty"`
}

type CreateSecurityGroupRequest struct {
	*common.BaseRequest

	SecurityGroupName string `json:"securityGroupName,omitempty"`

	RuleInfos []RuleInfo `json:"ruleInfos,omitempty"`
}

type CreateSecurityGroupResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateSecurityGroupResponseParams `json:"response"`
}

type CreateSecurityGroupResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	SecurityGroupId string `json:"securityGroupId,omitempty"`
}

type DeleteSecurityGroupRequest struct {
	*common.BaseRequest

	SecurityGroupId string `json:"securityGroupId,omitempty"`
}

type DeleteSecurityGroupResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ConfigureSecurityGroupRulesRequest struct {
	*common.BaseRequest

	SecurityGroupId string `json:"securityGroupId,omitempty"`

	RuleInfos []RuleInfo `json:"ruleInfos,omitempty"`
}

type ConfigureSecurityGroupRulesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type AssignSecurityGroupVpcRequest struct {
	*common.BaseRequest

	SecurityGroupId string `json:"securityGroupId,omitempty"`

	VpcId string `json:"vpcId,omitempty"`
}

type AssignSecurityGroupVpcResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type UnAssignSecurityGroupVpcRequest struct {
	*common.BaseRequest

	SecurityGroupId string `json:"securityGroupId,omitempty"`

	VpcId string `json:"vpcId,omitempty"`
}

type UnAssignSecurityGroupVpcResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DescribeNicsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeNicsResponseParams `json:"response"`
}

type DescribeNicsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*NicInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type NicInfo struct {
	NicId         string   `json:"nicId,omitempty"`
	Name          string   `json:"name,omitempty"`
	Status        string   `json:"status,omitempty"`
	NicType       string   `json:"nicType,omitempty"`
	RegionId      string   `json:"regionId,omitempty"`
	NicSubnetType string   `json:"nicSubnetType,omitempty"`
	PublicIpList  []string `json:"publicIpList,omitempty"`
	PrivateIpList []string `json:"privateIpList,omitempty"`
	InstanceId    string   `json:"instanceId,omitempty"`
	VpcId         string   `json:"vpcId,omitempty"`
	SubnetId      string   `json:"subnetId,omitempty"`
	CreateTime    string   `json:"createTime,omitempty"`
}

type DescribeNicsRequest struct {
	*common.BaseRequest

	NicIds []string `json:"nicIds,omitempty"`

	Name string `json:"name,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	VpcId string `json:"vpcId,omitempty"`

	SubnetId string `json:"subnetId,omitempty"`

	// 根据网卡关联的实例ID过滤。
	InstanceId string `json:"instanceId,omitempty"`

	Status string `json:"status,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type ModifyNicsAttributeResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyNicsAttributeRequest struct {
	*common.BaseRequest

	NicIds []string `json:"nicIds,omitempty"`

	Name string `json:"name,omitempty"`
}

type CreateNicResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *CreateNicResponseParams `json:"response"`
}

type CreateNicResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	NicId string `json:"nicId,omitempty"`
}

type CreateNicRequest struct {
	*common.BaseRequest

	Name string `json:"name,omitempty"`

	SubnetId string `json:"subnetId,omitempty"`

	PackageSize float64 `json:"packageSize"`

	Bandwidth int `json:"bandwidth"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	ClusterId string `json:"clusterId,omitempty"`
}

type DeleteNicResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DeleteNicRequest struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`
}

type AttachNicResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type AttachNicRequest struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	InstanceId string `json:"instanceId,omitempty"`
}

type AssignNicIpv6Response struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type AssignNicIpv6Request struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	PackageSize float64 `json:"packageSize"`

	Bandwidth int `json:"bandwidth"`

	ClusterId string `json:"clusterId"`
}

type UnAssignNicIpv4Response struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type UnAssignNicIpv4Request struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	IpAddress string `json:"ipAddress,omitempty"`
}

type AssignNicIpv4Response struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type AssignNicIpv4Request struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	IpAddress string `json:"ipAddress,omitempty"`
}

type BatchAssignNicIpv4Response struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type BatchAssignNicIpv4Request struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	IpAddresses []string `json:"ipAddresses,omitempty"`
}

type DetachNicResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DetachNicRequest struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`
}

type DescribeNicRegionsResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *DescribeNicRegionsResponseParams `json:"response"`
}

type DescribeNicRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionIds []string `json:"regionIds,omitempty"`
}

type DescribeNicRegionsRequest struct {
	*common.BaseRequest
}

type InquiryPricePublicIpv6Request struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	PackageSize float64 `json:"packageSize"`

	Bandwidth int `json:"bandwidth"`
}

type InquiryPricePublicIpv6Response struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *InquiryPricePublicIpv6ResponseParams `json:"response"`
}

type InquiryPricePublicIpv6ResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	BandwidthPrice *Price `json:"bandwidthPrice,omitempty"`
}

type SystemDisk struct {
	DiskId       string `json:"diskId,omitempty"`
	DiskSize     int    `json:"diskSize,omitempty"`
	DiskCategory string `json:"diskCategory,omitempty"`
}

type DataDisk struct {
	DiskId       string `json:"diskId,omitempty"`
	DiskSize     int    `json:"diskSize,omitempty"`
	DiskName     string `json:"diskName,omitempty"`
	DiskCategory string `json:"diskCategory,omitempty"`
	Portable     bool   `json:"portable,omitempty"`
	DiskAmount   int    `json:"diskAmount,omitempty"`
}

type CreateZecInstancesRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	InstanceType string `json:"instanceType,omitempty"`

	ImageId string `json:"imageId,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	InstanceName string `json:"instanceName,omitempty"`

	InstanceCount int `json:"instanceCount,omitempty"`

	Password string `json:"password,omitempty"`

	KeyId string `json:"keyId,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	TrafficPackageSize float64 `json:"trafficPackageSize,omitempty"`

	Bandwidth int `json:"bandwidth,omitempty"`

	SubnetId string `json:"subnetId,omitempty"`

	LanIp string `json:"lanIp,omitempty"`

	SystemDisk *SystemDisk `json:"systemDisk,omitempty"`

	DataDisks []*DataDisk `json:"dataDisks,omitempty"`

	TimeZone string `json:"timeZone,omitempty"`

	EnableAgent bool `json:"enableAgent"`

	EnableIpForward bool `json:"enableIpForward"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	ClusterId string `json:"clusterId,omitempty"`

	EipBindType *string `json:"eipBindType,omitempty"`
}

type CreateZecInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateZecInstancesResponseParams `json:"response"`
}

type CreateZecInstancesResponseParams struct {
	RequestId     *string             `json:"requestId,omitempty"`
	InstanceIdSet []*string           `json:"instanceIdSet,omitempty"`
	OrderNumber   *string             `json:"orderNumber,omitempty"`
	Instances     []*DiskWithInstance `json:"instances,omitempty"`
}

type DiskWithInstance struct {
	InstanceId string   `json:"instanceId,omitempty"`
	DiskIdSet  []string `json:"diskIdSet,omitempty"`
}

type CreateInstancesRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	InstanceType string `json:"instanceType,omitempty"`

	ImageId string `json:"imageId,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	InstanceName string `json:"instanceName,omitempty"`

	InstanceCount int `json:"instanceCount,omitempty"`

	Password string `json:"password,omitempty"`

	KeyId string `json:"keyId,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	TrafficPackageSize float64 `json:"trafficPackageSize,omitempty"`

	Bandwidth float64 `json:"bandwidth,omitempty"`

	SubnetId string `json:"subnetId,omitempty"`

	SystemDisk *SystemDisk `json:"systemDisk,omitempty"`

	DataDisks []*DataDisk `json:"dataDisks,omitempty"`

	TimeZone string `json:"timeZone,omitempty"`

	EnableAgent bool `json:"enableAgent"`

	EnableIpForward bool `json:"enableIpForward"`

	EipV4Type string `json:"eipV4Type,omitempty"`
}

type CreateInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateInstancesResponseParams `json:"response"`
}

type CreateInstancesResponseParams struct {
	RequestId     *string             `json:"requestId,omitempty"`
	InstanceIdSet []*string           `json:"instanceIdSet,omitempty"`
	OrderNumber   *string             `json:"orderNumber,omitempty"`
	Instances     []*DiskWithInstance `json:"instances,omitempty"`
}

type DescribeZoneInstanceConfigInfosRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	InstanceType string `json:"instanceType,omitempty"`
}

type DescribeZoneInstanceConfigInfosResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeZoneInstanceConfigInfosResponseParams `json:"response"`
}

type DescribeZoneInstanceConfigInfosResponseParams struct {
	RequestId            string                   `json:"requestId,omitempty"`
	InstanceTypeQuotaSet []*InstanceTypeQuotaItem `json:"instanceTypeQuotaSet,omitempty"`
}

type InstanceTypeQuotaItem struct {
	ZoneId                       string   `json:"zoneId,omitempty"`
	InstanceType                 string   `json:"instanceType,omitempty"`
	CpuCount                     int      `json:"cpuCount,omitempty"`
	Memory                       int      `json:"memory,omitempty"`
	Frequency                    string   `json:"frequency,omitempty"`
	InternetMaxBandwidthOutLimit int      `json:"internetMaxBandwidthOutLimit,omitempty"`
	InstanceTypeName             string   `json:"instanceTypeName,omitempty"`
	InternetChargeTypes          []string `json:"internetChargeTypes,omitempty"`
}

type DescribeZonesRequest struct {
	*common.BaseRequest

	ZoneIds []string `json:"zoneIds,omitempty"`
}

type DescribeZonesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
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

type DescribeImagesRequest struct {
	*common.BaseRequest

	ImageIds []string `json:"imageIds,omitempty"`

	ImageName string `json:"imageName,omitempty"`

	ZoneId string `json:"zoneId,omitempty"`

	Category string `json:"category,omitempty"`

	ImageType string `json:"imageType,omitempty"`

	OsType string `json:"osType,omitempty"`

	ImageStatus string `json:"imageStatus,omitempty"`

	PageNum int `json:"pageNum,omitempty"`

	PageSize int `json:"pageSize,omitempty"`
}

type DescribeImagesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeImagesResponseParams `json:"response"`
}

type DescribeImagesResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	// Information on an image.
	DataSet []*ImageInfo `json:"dataSet,omitempty"`

	// Number of images meeting the filtering conditions.
	TotalCount int `json:"totalCount,omitempty"`
}

type ImageInfo struct {
	ImageId string `json:"imageId,omitempty"`

	ImageName string `json:"imageName,omitempty"`

	ImageType string `json:"imageType,omitempty"`

	ImageSize string `json:"imageSize,omitempty"`

	ImageDescription string `json:"imageDescription,omitempty"`

	ImageVersion string `json:"imageVersion,omitempty"`

	ImageStatus string `json:"imageStatus,omitempty"`

	NicNetworkType []string `json:"nicNetworkType,omitempty"`

	Category string `json:"category,omitempty"`

	OsType string `json:"osType,omitempty"`
}

type DescribeKeyPairsRequest struct {
	*common.BaseRequest
	KeyIds   []string `json:"keyIds,omitempty"`
	KeyName  string   `json:"keyName,omitempty"`
	PageNum  int      `json:"pageNum,omitempty"`
	PageSize int      `json:"pageSize,omitempty"`
}

type DescribeKeyPairsResponse struct {
	*common.BaseResponse
	RequestId string                          `json:"requestId,omitempty"`
	Response  *DescribeKeyPairsResponseParams `json:"response"`
}

type DescribeKeyPairsResponseParams struct {
	RequestId  string     `json:"requestId,omitempty"`
	DataSet    []*KeyPair `json:"dataSet,omitempty"`
	TotalCount int        `json:"totalCount,omitempty"`
}

type KeyPair struct {
	// 密钥对ID。
	KeyId string `json:"keyId,omitempty"`

	// 密钥对名称。
	KeyName string `json:"keyName,omitempty"`

	// 密钥对的公钥内容。
	PublicKey string `json:"publicKey,omitempty"`

	// 密钥对描述信息。
	KeyDescription string `json:"keyDescription,omitempty"`

	// 创建时间。
	CreateTime string `json:"createTime,omitempty"`
}

type DescribeEipRegionsRequest struct {
	*common.BaseRequest
}

type DescribeEipRegionsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeEipRegionsResponseParams `json:"response"`
}

type DescribeEipRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionIds []string `json:"regionIds,omitempty"`
}

type DescribeEipRemoteRegionsRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`
}

type DescribeEipRemoteRegionsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeEipRemoteRegionsResponseParams `json:"response"`
}

type DescribeEipRemoteRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	PeerRegionIds []string `json:"peerRegionIds,omitempty"`
}

type DescribeEipInternetChargeTypesRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`
}

type DescribeEipInternetChargeTypesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeEipInternetChargeTypesResponseParams `json:"response"`
}

type DescribeEipInternetChargeTypesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	InternetChargeTypes []string `json:"internetChargeTypes,omitempty"`
}

type DescribeEipsRequest struct {
	*common.BaseRequest

	EipIds []string `json:"eipIds,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Name string `json:"name,omitempty"`

	Status string `json:"status,omitempty"`

	IsDefault bool `json:"isDefault,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`

	// PrivateIpAddress 按照 EIP 绑定的内网 IP 过滤。
	PrivateIpAddress *string `json:"privateIpAddress,omitempty"`

	// IpAddress 按照 EIP 的 IP 过滤。
	// Deprecated: 请使用 IpAddresses 代替
	IpAddress *string `json:"ipAddress,omitempty"`

	// IpAddresses 按照 EIP 的 IP列表过滤。
	IpAddresses []string `json:"ipAddresses,omitempty"`

	// InstanceId 按照 EIP 绑定的实例 ID 过滤。该字段过滤出该实例所绑定的网卡上的 EIP。
	InstanceId *string `json:"instanceId,omitempty"`

	// AssociatedId 按照 EIP 绑定的资源 ID 过滤。
	AssociatedId *string `json:"associatedId,omitempty"`

	// CidrIds 按照 EIP 所属的CIDR ID列表 过滤。
	CidrIds []string `json:"cidrIds,omitempty"`
}

type DescribeEipsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeEipsResponseParams `json:"response"`
}

type DescribeEipsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*EipInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type EipInfo struct {
	EipId string `json:"eipId,omitempty"`

	Name string `json:"name,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	PeerRegionId string `json:"peerRegionId,omitempty"`

	IsDefault bool `json:"isDefault,omitempty"`

	Status string `json:"status,omitempty"`

	PublicIpAddresses []string `json:"publicIpAddresses,omitempty"`

	// PrivateIpAddress EIP 绑定的内网IP地址。
	PrivateIpAddress *string `json:"privateIpAddress,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	CidrId string `json:"cidrId,omitempty"`

	NicId string `json:"nicId,omitempty"`

	// AssociatedId EIP 绑定的资源ID。可能为实例ID、网卡ID或者NAT网关ID。
	AssociatedId *string `json:"associatedId,omitempty"`

	// AssociatedType EIP 资源类型。可能为实例ID、网卡ID或者NAT网关ID。
	AssociatedType *string `json:"associatedType,omitempty"`

	// BindType EIP 绑定类型。
	BindType *string `json:"bindType,omitempty"`

	FlowPackage float64 `json:"flowPackage,omitempty"`

	Bandwidth int `json:"bandwidth,omitempty"`

	CreateTime string `json:"createTime,omitempty"`

	ExpiredTime string `json:"expiredTime,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	ResourceGroupName string `json:"resourceGroupName,omitempty"`
}

type CreateEipsRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	Amount int `json:"amount,omitempty"`

	Name string `json:"name,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	PrimaryIsp string `json:"primaryIsp,omitempty"`

	Bandwidth int `json:"bandwidth,omitempty"`

	FlowPackage float64 `json:"flowPackage,omitempty"`

	CidrId string `json:"cidrId,omitempty"`

	PublicIp string `json:"publicIp,omitempty"`

	ResourceGroupId string `json:"resourceGroupId,omitempty"`

	ClusterId string `json:"clusterId,omitempty"`

	PeerRegionId string `json:"peerRegionId,omitempty"`
}

type CreateEipsResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *CreateEipsResponseParams `json:"response"`
}

type CreateEipsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	EipIds []string `json:"eipIds,omitempty"`

	OrderNumber string `json:"orderNumber,omitempty"`
}

type DeleteEipRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`
}

type DeleteEipResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type RenewEipRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`
}

type RenewEipResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type BatchAttachEipLanIpRequest struct {
	*common.BaseRequest

	NicId string `json:"nicId,omitempty"`

	LanIp string `json:"lanIp,omitempty"`

	EipIds []string `json:"eipIds,omitempty"`
}

type BatchAttachEipLanIpResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}



type AssociateEipAddressRequest struct {
    *common.BaseRequest

	// LoadBalancerId 负载均衡实例的ID。
	LoadBalancerId *string `json:"loadBalancerId,omitempty"`

	// NicId 虚拟网卡的ID。
	NicId *string `json:"nicId,omitempty"`

	// LanIp 要绑定的网卡上的内网IP地址。当IP绑定的是网卡, 则该字段不能为空。
	LanIp *string `json:"lanIp,omitempty"`

	// NatId NAT网关的ID。
	NatId *string `json:"natId,omitempty"`

	// EipIds 要绑定的EIP的ID。EIP 必须是未绑定状态。
	EipIds []string `json:"eipIds,omitempty"`

	// BindType 绑定类型。当绑定的是网卡时生效。默认为普通NATm模式。
	BindType *string `json:"bindType,omitempty"`
}

type AssociateEipAddressResponse struct {
    *common.BaseResponse

    RequestId string `json:"requestId,omitempty"`

    Response *AssociateEipAddressResponseParams `json:"response"`
}

type AssociateEipAddressResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	FailedEipIds []string `json:"failedEipIds,omitempty"`
}

type DetachEipLanIpRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`
}

type DetachEipLanIpResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type UnassociateEipAddressRequest struct {
    *common.BaseRequest

    EipIds []string `json:"eipIds,omitempty"`  // 需要解绑的EIP ID列表
}

type UnassociateEipAddressResponse struct {
    *common.BaseResponse

    RequestId string                                 `json:"requestId,omitempty"`

    Response  *UnassociateEipAddressResponseParams   `json:"response"`
}

type UnassociateEipAddressResponseParams struct {
    RequestId     string   `json:"requestId,omitempty"`

    FailedEipIds  []string `json:"failedEipIds,omitempty"`  // 解绑失败的EIP ID列表
}


type ReplaceEipAddressRequest struct {
    *common.BaseRequest

    ReplaceIps  []*ReplaceIp `json:"replaceIps,omitempty"`
}

type ReplaceIp struct {
	EipId string `json:"eipId,omitempty"`

	OwnIp string `json:"ownIp,omitempty"`

	TargetIp string `json:"targetIp,omitempty"`
}

type ReplaceEipAddressResponse struct {
    *common.BaseResponse

    RequestId string `json:"requestId,omitempty"`

    Response *ReplaceEipAddressResponseParams `json:"response"`
}

type ReplaceEipAddressResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	FailedEipIds []string `json:"failedEipIds,omitempty"`
}


type ConfigEipEgressIpRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`
}

type ConfigEipEgressIpResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DescribeEipPriceRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	Amount int `json:"amount,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	Bandwidth int `json:"bandwidth,omitempty"`

	FlowPackage float64 `json:"flowPackage,omitempty"`

	CidrId string `json:"cidrId,omitempty"`

	ClusterId string `json:"clusterId,omitempty"`

	PeerRegionId string `json:"peerRegionId,omitempty"`
}

type DescribeEipPriceResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *DescribeEipPriceResponseParams `json:"response"`
}

type DescribeEipPriceResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	EipPrice *Price `json:"eipPrice,omitempty"`

	BandwidthPrice *Price `json:"bandwidthPrice,omitempty"`
}

type ChangeEipInternetChargeTypeRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	ClusterId string `json:"clusterId,omitempty"`
}

type ChangeEipInternetChargeTypeResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *ChangeEipInternetChargeTypeResponseParams `json:"response"`
}

type ChangeEipInternetChargeTypeResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	OrderNumber string `json:"orderNumber,omitempty"`
}

type InquiryPriceCreateInstanceRequest struct {
	*common.BaseRequest

	ZoneId string `json:"zoneId,omitempty"`

	InstanceType string `json:"instanceType,omitempty"`

	EipV4Type string `json:"eipV4Type,omitempty"`

	InternetChargeType string `json:"internetChargeType,omitempty"`

	TrafficPackageSize float64 `json:"trafficPackageSize,omitempty"`

	Bandwidth int `json:"bandwidth,omitempty"`

	InstanceCount int `json:"instanceCount,omitempty"`

	SystemDisk *SystemDisk `json:"systemDisk,omitempty"`

	DataDisk *DataDisk `json:"dataDisk,omitempty"`
}

type InquiryPriceCreateInstanceResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *InquiryPriceCreateInstanceResponseParams `json:"response"`
}

type InquiryPriceCreateInstanceResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	SpecPrice *Price `json:"specPrice,omitempty"`

	GpuPrice *Price `json:"gpuPrice,omitempty"`

	Ipv4Price *Price `json:"ipv4Price,omitempty"`

	Ipv6Price *Price `json:"ipv6Price,omitempty"`

	Ipv4BandwidthPrice *Price `json:"ipv4BandwidthPrice,omitempty"`

	Ipv6BandwidthPrice *Price `json:"ipv6BandwidthPrice,omitempty"`

	SystemDiskPrice *Price `json:"systemDiskPrice,omitempty"`

	DataDiskPrice *Price `json:"dataDiskPrice,omitempty"`
}

type ReleaseInstancesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *ReleaseInstancesResponseParams `json:"response"`
}

type ReleaseInstancesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type ReleaseInstancesRequest struct {
	*common.BaseRequest

	InstanceIds []string `json:"instanceIds,omitempty"`
}

type ModifyInstanceTypeResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *ModifyInstanceTypeResponseParams `json:"response"`
}

type ModifyInstanceTypeResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	OrderNumber string `json:"orderNumber,omitempty"`
}

type ModifyInstanceTypeRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`

	InstanceType string `json:"instanceType,omitempty"`
}

type DescribeTimeZonesResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeTimeZonesResponseParams `json:"response"`
}

type DescribeTimeZonesResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	TimeZones []string `json:"timeZones,omitempty"`
}

type DescribeTimeZonesRequest struct {
	*common.BaseRequest
}

type StartIpForwardResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type StartIpForwardRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`
}

type StopIpForwardResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type StopIpForwardRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`
}

type StartAgentMonitorResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type StartAgentMonitorRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`
}

type StopAgentMonitorResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type StopAgentMonitorRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`
}

type DescribeVncUrlResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeVncUrlResponseParams `json:"response"`
}

type DescribeVncUrlResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	Url string `json:"url,omitempty"`
}

type DescribeVncUrlRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`
}

type ChangeNicNetworkTypeResponse struct {
	*common.BaseResponse

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 requestId。
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ChangeNicNetworkTypeRequest struct {
	*common.BaseRequest

	InstanceId string `json:"instanceId,omitempty"`

	NicNetworkType string `json:"nicNetworkType,omitempty"`
}

type DescribeEipTrafficRequest struct {
	*common.BaseRequest

	EipId string `json:"eipId,omitempty"`

	StartTime string `json:"startTime,omitempty"`

	EndTime string `json:"endTime,omitempty"`

	Step int `json:"step,omitempty"`

	WanIp string `json:"wanIp,omitempty"`
}

type DescribeEipTrafficResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response *EipTrafficDataResponse `json:"response"`
}

type EipTrafficDataResponse struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	DataList []*EipTrafficData `json:"dataList,omitempty"`

	In95 int64 `json:"in95,omitempty"`

	InAvg int64 `json:"inAvg,omitempty"`

	InMax int64 `json:"inMax,omitempty"`

	InMin int64 `json:"inMin,omitempty"`

	InTotal int64 `json:"inTotal,omitempty"`

	Out95 int64 `json:"out95,omitempty"`

	OutAvg int64 `json:"outAvg,omitempty"`

	OutMax int64 `json:"outMax,omitempty"`

	OutMin int64 `json:"outMin,omitempty"`

	OutTotal int64 `json:"outTotal,omitempty"`
}

type EipTrafficData struct {
	InternetRX int64 `json:"internetRX,omitempty"`

	InternetTX int64 `json:"internetTX,omitempty"`

	Time string `json:"time,omitempty"`
}

type CreateBorderGatewayRequest struct {
	*common.BaseRequest

	RegionId string `json:"regionId,omitempty"`

	VpcId string `json:"vpcId,omitempty"`

	Label string `json:"label,omitempty"`

	Asn int `json:"asn,omitempty"`

	AdvertisedSubnet string `json:"advertisedSubnet,omitempty"`

	AdvertisedCidrs []string `json:"advertisedCidrs,omitempty"`
}

type CreateBorderGatewayResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response *CreateBorderGatewayResponseParams `json:"response"`
}

type CreateBorderGatewayResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	ZbgId string `json:"zbgId,omitempty"`
}

type DeleteBorderGatewayRequest struct {
	*common.BaseRequest

	ZbgId string `json:"zbgId,omitempty"`

	Asn int `json:"asn,omitempty"`
}

type DeleteBorderGatewayResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type DescribeBorderGatewaysRequest struct {
	*common.BaseRequest

	ZbgIds []string `json:"zbgIds,omitempty"`

	Name string `json:"name,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	PageSize int `json:"pageSize,omitempty"`

	PageNum int `json:"pageNum,omitempty"`
}

type DescribeBorderGatewaysResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeBorderGatewaysResponseParams `json:"response"`
}

type DescribeBorderGatewaysResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	DataSet []*ZbgInfo `json:"dataSet,omitempty"`

	TotalCount int `json:"totalCount,omitempty"`
}

type ZbgInfo struct {
	ZbgId string `json:"zbgId,omitempty"`

	Name string `json:"name,omitempty"`

	VpcId string `json:"vpcId,omitempty"`

	RegionId string `json:"regionId,omitempty"`

	Asn int `json:"asn,omitempty"`

	InterConnectCidr string `json:"interConnectCidr,omitempty"`

	CreateTime string `json:"createTime,omitempty"`

	CloudRouterIds []string `json:"cloudRouterIds,omitempty"`

	AdvertisedSubnet string `json:"advertisedSubnet,omitempty"`

    AdvertisedCidrs []string `json:"advertisedCidrs,omitempty"`

    NatId string `json:"natId,omitempty"`
}

type ModifyBorderGatewaysAttributeRequest struct {
	*common.BaseRequest

	ZbgIds []string `json:"zbgIds,omitempty"`

	Name string `json:"name,omitempty"`

	AdvertisedSubnet string `json:"advertisedSubnet,omitempty"`

    AdvertisedCidrs []string `json:"advertisedCidrs,omitempty"`
}

type ModifyBorderGatewaysAttributeResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type ModifyBorderGatewayAsnRequest struct {
	*common.BaseRequest

	ZbgId string `json:"zbgId,omitempty"`

	Asn int `json:"asn,omitempty"`
}

type ModifyBorderGatewayAsnResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}


type DescribeAvailableNatsRequest struct {
	*common.BaseRequest

	ZbgId string `json:"zbgId,omitempty"`
}


type DescribeAvailableNatsResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response *DescribeAvailableNatsResponseParams `json:"response,omitempty"`
}

type DescribeAvailableNatsResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // NAT网关ID集合。
    NatIds []string `json:"natIds,omitempty"`
}


type AssignBorderGatewayRequest struct {
	*common.BaseRequest

	ZbgId string `json:"zbgId,omitempty"`

	NatId string `json:"natId,omitempty"`
}

type AssignBorderGatewayResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}


type UnassignBorderGatewayRequest struct {
	*common.BaseRequest

	ZbgId string `json:"zbgId,omitempty"`
}

type UnassignBorderGatewayResponse struct {
	*common.BaseResponse

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId string `json:"requestId,omitempty"`

	Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response"`
}

type CreateNatGatewayRequest struct {
    *common.BaseRequest

    // RegionId 区域节点ID。
    RegionId *string `json:"regionId,omitempty"`

    // VpcId NAT网关所属的VPC网络ID。
    VpcId *string `json:"vpcId,omitempty"`

    // Name NAT网关的名称。长度为2～63个字符。
    Name *string `json:"name,omitempty"`

    // SubnetIds NAT网关所属的Subnet子网ID集合。
    SubnetIds []string `json:"subnetIds,omitempty"`

    // ResourceGroupId 资源组ID。如果不指定，则会创建在默认资源组。
    ResourceGroupId *string `json:"resourceGroupId,omitempty"`

}

type CreateNatGatewayResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // OrderNumber 下单编号。
    OrderNumber *string `json:"orderNumber,omitempty"`

    // NatGatewayId NAT网关唯一ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

}

type CreateNatGatewayResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *CreateNatGatewayResponseParams `json:"response,omitempty"`

}

type DescribeNatGatewaysRequest struct {
    *common.BaseRequest

    // RegionId 节点ID。
    RegionId *string `json:"regionId,omitempty"`

    // VpcId 根据NAT网关所属的VPC网络 ID过滤。
    VpcId *string `json:"vpcId,omitempty"`

    // NatGatewayIds NAT网关ID集合。实例ID数量上限为100个。
    NatGatewayIds []string `json:"natGatewayIds,omitempty"`

    // Name NAT网关名称。可以在末尾加入*以支持模糊匹配。
    Name *string `json:"name,omitempty"`

    // Status NAT网关状态。
    Status *string `json:"status,omitempty"`

    // PageSize 返回的分页大小。默认为20，最大为1000。
    PageSize *int `json:"pageSize,omitempty"`

    // PageNum 返回的分页页码。默认为1。
    PageNum *int `json:"pageNum,omitempty"`

}

type DescribeNatGatewaysResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // TotalCount 满足过滤条件的NAT网关总数。
    TotalCount *int `json:"totalCount,omitempty"`

    // DataSet 返回的NAT网关列表。
    DataSet []*NatGateway `json:"dataSet,omitempty"`

}

// NatGateway 描述NAT网关的信息。
type NatGateway struct {

    // NatGatewayId NAT网关唯一ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

    // VpcId NAT网关所属的VPC网络ID。
    VpcId *string `json:"vpcId,omitempty"`

    // RegionId 区域节点ID。
    RegionId *string `json:"regionId,omitempty"`

    // Status NAT网关的状态。
    Status *string `json:"status,omitempty"`

    // Name NAT网关的名称。
    Name *string `json:"name,omitempty"`

    // SubnetIds NAT网关所属的Subnet子网ID集合。
    SubnetIds []string `json:"subnetIds,omitempty"`

    // EipIds NAT网关所关联的EIP ID集合。
    EipIds []string `json:"eipIds,omitempty"`

    // ZbgId 边界网关 ID。
    ZbgId []string `json:"zbgId,omitempty"`

    // ResourceGroupId 资源组ID。
    ResourceGroupId *string `json:"resourceGroupId,omitempty"`

    // ResourceGroupName 资源组名称。
    ResourceGroupName *string `json:"resourceGroupName,omitempty"`

    // CreateTime 创建时间。按照ISO8601标准表示，并且使用UTC时间, 格式为：YYYY-MM-DDThh:mm:ssZ。
    CreateTime *string `json:"createTime,omitempty"`

    // ExpiredTime 到期时间。按照ISO8601标准表示，并且使用UTC时间, 格式为：YYYY-MM-DDThh:mm:ssZ。
    ExpiredTime *string `json:"expiredTime,omitempty"`

}

type DescribeNatGatewaysResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *DescribeNatGatewaysResponseParams `json:"response,omitempty"`

}

type DescribeNatGatewayDetailRequest struct {
    *common.BaseRequest

    // NatGatewayId NAT网关 ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

}

type DescribeNatGatewayDetailResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // NatGatewayId NAT网关唯一ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

    // Name NAT网关名称。
    Name *string `json:"name,omitempty"`

    // Snats SNAT网关规则集合。
    Snats []*SnatEntry `json:"snats,omitempty"`

    // Dnats DNAT网关规则集合。
    Dnats []*DnatEntry `json:"dnats,omitempty"`

}

// SnatEntry 描述SNAT规则的信息。
type SnatEntry struct {

    // SnatEntryId SNAT规则 ID。
    SnatEntryId *string `json:"snatEntryId,omitempty"`

    // Cidrs CIDR网段，不传默认是0.0.0.0/0。
    Cidrs []string `json:"cidrs,omitempty"`

    // EipIds SNAT规则添加的eip ID集合。
    EipIds []string `json:"eipIds,omitempty"`

}

// DnatEntry 描述DNAT规则的信息。
type DnatEntry struct {

    // DnatEntryId DNAT规则 ID。
    DnatEntryId *string `json:"dnatEntryId,omitempty"`

    // Status DNAT规则状态。
    Status *string `json:"status,omitempty"`

    // PrivateIp DNAT规则的内网IP地址。
    PrivateIp *string `json:"privateIp,omitempty"`

    // EipId DNAT规则添加的eip ID。
    EipId *string `json:"eipId,omitempty"`

    // Protocol DNAT规则的协议类型。
    Protocol *string `json:"protocol,omitempty"`

    // ListenerPort DNAT规则端口转发的外部端口或端口段，取值范围1-65535。
    ListenerPort *string `json:"listenerPort,omitempty"`

    // InternalPort DNAT规则端口转发的内部端口或端口段，取值范围1-65535。
    InternalPort *string `json:"internalPort,omitempty"`

}

type DescribeNatGatewayDetailResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *DescribeNatGatewayDetailResponseParams `json:"response,omitempty"`

}

type DeleteNatGatewayRequest struct {
    *common.BaseRequest

    // NatGatewayId NAT网关 ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

}

type DeleteNatGatewayResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type RenewNatGatewayRequest struct {
    *common.BaseRequest

    // NatGatewayId NAT网关 ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

}

type RenewNatGatewayResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type InquiryPriceCreateNatGatewayRequest struct {
    *common.BaseRequest

    // RegionId 区域节点ID。
    RegionId *string `json:"regionId,omitempty"`

}

// InquiryPriceCreateNatGatewayResponseParams 描述创建NAT网关的响应。
type InquiryPriceCreateNatGatewayResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // NatGatewayPrice NAT网关的价格。
    NatGatewayPrice *PriceItem `json:"natGatewayPrice,omitempty"`

    // CuPrice NAT网关CU的价格。
    CuPrice *PriceItem `json:"cuPrice,omitempty"`

}

// PriceItem 描述价格的信息。
type PriceItem struct {

    // Discount 折扣大小。如80.0代表8折。
    Discount *float64 `json:"discount,omitempty"`

    // DiscountPrice 后付费的单元折后价格。后付费模式使用，如果价格为阶梯价格，该项为null。
    DiscountPrice *float64 `json:"discountPrice,omitempty"`

    // OriginalPrice 预付费的原价。预付费模式使用，后付费该值为 null。
    OriginalPrice *float64 `json:"originalPrice,omitempty"`

    // UnitPrice 后付费的单元原始价格。后付费模式使用，如果价格为阶梯价格，该项为null。
    UnitPrice *float64 `json:"unitPrice,omitempty"`

    // DiscountUnitPrice 后付费的单元折后价格。后付费模式使用，如果价格为阶梯价格，该项为null。
    DiscountUnitPrice *float64 `json:"discountUnitPrice,omitempty"`

    // ChargeUnit 后付费计价单元。后付费模式使用，可取值范围：<br/>HOUR: 表示计价单元是按每小时来计算。DAY: 表示计价单元是按天来计算。MONTH: 表示计价单元是按月来计算，95计费则是这种。
    ChargeUnit *string `json:"chargeUnit,omitempty"`

    // StepPrices 后付费阶梯价格。后付费模式使用，如果非阶梯价格，该项为null。
    StepPrices []*StepPrice `json:"stepPrices,omitempty"`

    // AmountUnit 用量单位。比如Mbps, LCU等。如果为null, 代表取不到值。
    AmountUnit *string `json:"amountUnit,omitempty"`

    // ExcessUnitPrice 超量原始价格。
    ExcessUnitPrice *float64 `json:"excessUnitPrice,omitempty"`

    // ExcessDiscountUnitPrice 超量折扣后价格。
    ExcessDiscountUnitPrice *float64 `json:"excessDiscountUnitPrice,omitempty"`

    // ExcessAmountUnit 超量用量单位。如果为null, 代表取不到值。
    ExcessAmountUnit *string `json:"excessAmountUnit,omitempty"`

}


type InquiryPriceCreateNatGatewayResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *InquiryPriceCreateNatGatewayResponseParams `json:"response,omitempty"`

}

type CreateSnatEntryRequest struct {
    *common.BaseRequest

    // NatGatewayId NAT网关 ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

    // Cidrs CIDR网段，不传默认是0.0.0.0/0。
    Cidrs []string `json:"cidrs,omitempty"`

    // EipIds
    EipIds []string `json:"eipIds,omitempty"`

}

type CreateSnatEntryResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // SnatEntryId SNAT规则唯一ID。
    SnatEntryId *string `json:"snatEntryId,omitempty"`

}

type CreateSnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *CreateSnatEntryResponseParams `json:"response,omitempty"`

}

type ModifySnatEntryRequest struct {
    *common.BaseRequest

    // SnatEntryId SNAT规则 ID。
    SnatEntryId *string `json:"snatEntryId,omitempty"`

    // Cidrs CIDR网段，不传默认是0.0.0.0/0。
    Cidrs []string `json:"cidrs,omitempty"`

    // EipIds SNAT规则添加的eip ID集合。
    EipIds []string `json:"eipIds,omitempty"`

}

type ModifySnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type DeleteSnatEntryRequest struct {
    *common.BaseRequest

    // SnatEntryId SNAT规则 ID。
    SnatEntryId *string `json:"snatEntryId,omitempty"`

}

type DeleteSnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type CreateDnatEntryRequest struct {
    *common.BaseRequest

    // NatGatewayId NAT网关 ID。
    NatGatewayId *string `json:"natGatewayId,omitempty"`

    // EipId
    EipId *string `json:"eipId,omitempty"`

    // PrivateIp DNAT规则的内网IP地址。
    PrivateIp *string `json:"privateIp,omitempty"`

    // Protocol DNAT规则的协议类型。
    Protocol *string `json:"protocol,omitempty"`

    // ListenerPort DNAT规则端口转发的外部端口或端口段，取值范围1-65535。
    ListenerPort *string `json:"listenerPort,omitempty"`

    // InternalPort DNAT规则端口转发的内部端口或端口段，取值范围1-65535。
    InternalPort *string `json:"internalPort,omitempty"`

}

type CreateDnatEntryResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

    // DnatEntryId DNAT规则唯一ID。
    DnatEntryId *string `json:"dnatEntryId,omitempty"`

}

type CreateDnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *CreateDnatEntryResponseParams `json:"response,omitempty"`

}

type ModifyDnatEntryRequest struct {
    *common.BaseRequest

    // DnatEntryId DNAT规则 ID。
    DnatEntryId *string `json:"dnatEntryId,omitempty"`

    // EipId
    EipId *string `json:"eipId,omitempty"`

    // PrivateIp DNAT规则的内网IP地址。
    PrivateIp *string `json:"privateIp,omitempty"`

    // Protocol DNAT规则的协议类型。
    Protocol *string `json:"protocol,omitempty"`

    // ListenerPort DNAT规则端口转发的外部端口或端口段，取值范围1-65535。
    ListenerPort *string `json:"listenerPort,omitempty"`

    // InternalPort DNAT规则端口转发的内部端口或端口段，取值范围1-65535。
    InternalPort *string `json:"internalPort,omitempty"`

}

type ModifyDnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`

}

type DeleteDnatEntryRequest struct {
    *common.BaseRequest

    // DnatEntryId DNAT规则 ID。
    DnatEntryId *string `json:"dnatEntryId,omitempty"`

}

type DeleteDnatEntryResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response struct {
		RequestId string `json:"requestId,omitempty"`
	} `json:"response,omitempty"`
}



type DescribeNatGatewayRegionsResponse struct {
	*common.BaseResponse

	RequestId string `json:"requestId,omitempty"`

	Response *DescribeNatGatewayRegionsResponseParams `json:"response"`
}

type DescribeNatGatewayRegionsResponseParams struct {
	RequestId string `json:"requestId,omitempty"`

	RegionIds []string `json:"regionIds,omitempty"`
}

type DescribeNatGatewayRegionsRequest struct {
	*common.BaseRequest
}

// ModifyEipBandwidthRequest 修改带宽限速请求。
type ModifyEipBandwidthRequest struct {
    *common.BaseRequest

    // EipId EIP唯一标识ID。
    EipId *string `json:"eipId,omitempty"`

    // Bandwidth 调整带宽限速的目标值，有且仅当为95计费时, 需要指定。单位Mbps。
    Bandwidth *int `json:"bandwidth,omitempty"`

    // CommitBandwidth 保底带宽。单位Mbps。
    CommitBandwidth *int `json:"commitBandwidth,omitempty"`

}

type ModifyEipBandwidthResponseParams struct {

    RequestId *string `json:"requestId,omitempty"`

}

type ModifyEipBandwidthResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *ModifyEipBandwidthResponseParams `json:"response,omitempty"`

}

// DescribeEipMonitorDataRequest 查询弹性公网IP监控指标请求。
type DescribeEipMonitorDataRequest struct {
	*common.BaseRequest

	// EipId EIP唯一标识ID。
	EipId *string `json:"eipId,omitempty"`

	// MetricType EIP监控指标类型。
	MetricType *string `json:"metricType,omitempty"`

	// StartTime 查询开始时间。时间格式：yyyy-MM-ddTHH:mm:ssZ。
	StartTime *string `json:"startTime,omitempty"`

	// EndTime 查询结束时间。时间格式：yyyy-MM-ddTHH:mm:ssZ。
	EndTime *string `json:"endTime,omitempty"`

}

// DescribeEipMonitorDataResponseParams 查询弹性IP监控数据的响应信息。
type DescribeEipMonitorDataResponseParams struct {

	RequestId *string `json:"requestId,omitempty"`

	// InMaxValue 入方向最大值。
	InMaxValue *float64 `json:"inMaxValue,omitempty"`

	// InAvgValue 入方向平均值。
	InAvgValue *float64 `json:"inAvgValue,omitempty"`

	// InMinValue 入方向最小值。
	InMinValue *float64 `json:"inMinValue,omitempty"`

	// InTotalValue 入方向总和值。
	InTotalValue *float64 `json:"inTotalValue,omitempty"`

	// OutMaxValue 出方向最大值。
	OutMaxValue *float64 `json:"outMaxValue,omitempty"`

	// OutAvgValue 出方向平均值。
	OutAvgValue *float64 `json:"outAvgValue,omitempty"`

	// OutMinValue 出方向最小值。
	OutMinValue *float64 `json:"outMinValue,omitempty"`

	// OutTotalValue 出方向总和值。
	OutTotalValue *float64 `json:"outTotalValue,omitempty"`

	// LoseOutMaxValue 丢失出方向最大值。
	LoseOutMaxValue *float64 `json:"loseOutMaxValue,omitempty"`

	// LoseOutMinValue 丢失出方向最小值。
	LoseOutMinValue *float64 `json:"loseOutMinValue,omitempty"`

	// LoseOutTotalValue 丢失出方向总和值。
	LoseOutTotalValue *float64 `json:"loseOutTotalValue,omitempty"`

	// LoseInMaxValue 丢失入方向最大值。
	LoseInMaxValue *float64 `json:"loseInMaxValue,omitempty"`

	// LoseInMinValue 丢失入方向最小值。
	LoseInMinValue *float64 `json:"loseInMinValue,omitempty"`

	// LoseInTotalValue 丢失入方向总和值。
	LoseInTotalValue *float64 `json:"loseInTotalValue,omitempty"`

	// DataList 监控数据集合。
	DataList []*EipMetricValue `json:"dataList,omitempty"`

}

// EipMetricValue 描述EIP的监控指标数据。
type EipMetricValue struct {

	// Time 数据点时间。
	Time *string `json:"time,omitempty"`

	// In 入方向值。
	In *float64 `json:"in,omitempty"`

	// Out 入方向值。
	Out *float64 `json:"out,omitempty"`

	// LoseIn 丢失入方向值。
	LoseIn *float64 `json:"loseIn,omitempty"`

	// LoseOut 丢失出方向值。
	LoseOut *float64 `json:"loseOut,omitempty"`

}

type DescribeEipMonitorDataResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *DescribeEipMonitorDataResponseParams `json:"response,omitempty"`

}

// DescribeInstanceMonitorDataRequest 查询实例监控指标请求。
type DescribeInstanceMonitorDataRequest struct {
	*common.BaseRequest

	// InstanceId 实例唯一标识ID。
	InstanceId *string `json:"instanceId,omitempty"`

	// MetricType 实例监控指标类型。
	MetricType *string `json:"metricType,omitempty"`

	// StartTime 查询开始时间。时间格式：yyyy-MM-ddTHH:mm:ssZ。
	StartTime *string `json:"startTime,omitempty"`

	// EndTime 查询结束时间。时间格式：yyyy-MM-ddTHH:mm:ssZ。
	EndTime *string `json:"endTime,omitempty"`

}

// DescribeInstanceMonitorDataResponseParams 查询实例监控数据的响应信息。
type DescribeInstanceMonitorDataResponseParams struct {

	RequestId *string `json:"requestId,omitempty"`

	// MaxValue 数据点的最大值。
	MaxValue *float64 `json:"maxValue,omitempty"`

	// MinValue 数据点的最小值。
	MinValue *float64 `json:"minValue,omitempty"`

	// AvgValue 数据点的平均值。
	AvgValue *float64 `json:"avgValue,omitempty"`

	// Metrics 监控数据集合。
	Metrics []*MetricValue `json:"metrics,omitempty"`

}

// MetricValue 描述实例监控指标的数据值。
type MetricValue struct {

	// Time 数据点时间。
	Time *string `json:"time,omitempty"`

	// Value 数据点的值。
	Value *float64 `json:"value,omitempty"`

}

type DescribeInstanceMonitorDataResponse struct {
    *common.BaseResponse

    RequestId *string `json:"requestId,omitempty"`

    Response *DescribeInstanceMonitorDataResponseParams `json:"response,omitempty"`

}

