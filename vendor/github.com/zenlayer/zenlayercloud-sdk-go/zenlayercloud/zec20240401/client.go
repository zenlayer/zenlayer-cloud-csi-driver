package zec

import (
	"github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/common"
)

const (
	APIVersion = "2024-04-01"
	SERVICE    = "zec"
)

type Client struct {
	common.Client
}

func NewClientWithSecretKey(secretKeyId, secretKeyPassword string) (client *Client, err error) {
	return NewClient(common.NewConfig(), secretKeyId, secretKeyPassword)
}

func NewClient(config *common.Config, secretKeyId, secretKeyPassword string) (client *Client, err error) {
	client = &Client{}

	err = client.InitWithCredential(common.NewCredential(secretKeyId, secretKeyPassword))
	if err != nil {
		return nil, err
	}
	err = client.WithConfig(config)

	if err != nil {
		return nil, err
	}
	return client, nil
}


func NewCreateVpcRequest() (request *CreateVpcRequest) {
	request = &CreateVpcRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateVpc")

	return
}

func NewCreateVpcResponse() (response *CreateVpcResponse) {
	response = &CreateVpcResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateVpc(request *CreateVpcRequest) (response *CreateVpcResponse, err error) {
	response = NewCreateVpcResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteVpcRequest() (request *DeleteVpcRequest) {
	request = &DeleteVpcRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteVpc")

	return
}

func NewDeleteVpcResponse() (response *DeleteVpcResponse) {
	response = &DeleteVpcResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteVpc(request *DeleteVpcRequest) (response *DeleteVpcResponse, err error) {
	response = NewDeleteVpcResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyVpcsAttributeRequest() (request *ModifyVpcsAttributeRequest) {
	request = &ModifyVpcsAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyVpcsAttribute")

	return
}

func NewModifyVpcsAttributeResponse() (response *ModifyVpcsAttributeResponse) {
	response = &ModifyVpcsAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyVpcsAttribute(request *ModifyVpcsAttributeRequest) (response *ModifyVpcsAttributeResponse, err error) {
	response = NewModifyVpcsAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeVpcsRequest() (request *DescribeVpcsRequest) {
	request = &DescribeVpcsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeVpcs")

	return
}

func NewDescribeVpcsResponse() (response *DescribeVpcsResponse) {
	response = &DescribeVpcsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcs(request *DescribeVpcsRequest) (response *DescribeVpcsResponse, err error) {
	response = NewDescribeVpcsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeSubnetRegionsRequest() (request *DescribeSubnetRegionsRequest) {
	request = &DescribeSubnetRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeSubnetRegions")

	return
}

func NewDescribeSubnetRegionsResponse() (response *DescribeSubnetRegionsResponse) {
	response = &DescribeSubnetRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeSubnetRegions(request *DescribeSubnetRegionsRequest) (response *DescribeSubnetRegionsResponse, err error) {
	response = NewDescribeSubnetRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateSubnetRequest() (request *CreateSubnetRequest) {
	request = &CreateSubnetRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateSubnet")

	return
}

func NewCreateSubnetResponse() (response *CreateSubnetResponse) {
	response = &CreateSubnetResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateSubnet(request *CreateSubnetRequest) (response *CreateSubnetResponse, err error) {
	response = NewCreateSubnetResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteSubnetRequest() (request *DeleteSubnetRequest) {
	request = &DeleteSubnetRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteSubnet")

	return
}

func NewDeleteSubnetResponse() (response *DeleteSubnetResponse) {
	response = &DeleteSubnetResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteSubnet(request *DeleteSubnetRequest) (response *DeleteSubnetResponse, err error) {
	response = NewDeleteSubnetResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifySubnetsAttributeRequest() (request *ModifySubnetsAttributeRequest) {
	request = &ModifySubnetsAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifySubnetsAttribute")

	return
}

func NewModifySubnetsAttributeResponse() (response *ModifySubnetsAttributeResponse) {
	response = &ModifySubnetsAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifySubnetsAttribute(request *ModifySubnetsAttributeRequest) (response *ModifySubnetsAttributeResponse, err error) {
	response = NewModifySubnetsAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeSubnetsRequest() (request *DescribeSubnetsRequest) {
	request = &DescribeSubnetsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeSubnets")

	return
}

func NewDescribeSubnetsResponse() (response *DescribeSubnetsResponse) {
	response = &DescribeSubnetsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeSubnets(request *DescribeSubnetsRequest) (response *DescribeSubnetsResponse, err error) {
	response = NewDescribeSubnetsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifySubnetStackTypeRequest() (request *ModifySubnetStackTypeRequest) {
	request = &ModifySubnetStackTypeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifySubnetStackType")

	return
}

func NewModifySubnetStackTypeResponse() (response *ModifySubnetStackTypeResponse) {
	response = &ModifySubnetStackTypeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifySubnetStackType(request *ModifySubnetStackTypeRequest) (response *ModifySubnetStackTypeResponse, err error) {
	response = NewModifySubnetStackTypeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateRouteRequest() (request *CreateRouteRequest) {
	request = &CreateRouteRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateRoute")

	return
}

func NewCreateRouteResponse() (response *CreateRouteResponse) {
	response = &CreateRouteResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateRoute(request *CreateRouteRequest) (response *CreateRouteResponse, err error) {
	response = NewCreateRouteResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteRouteRequest() (request *DeleteRouteRequest) {
	request = &DeleteRouteRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteRoute")

	return
}

func NewDeleteRouteResponse() (response *DeleteRouteResponse) {
	response = &DeleteRouteResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteRoute(request *DeleteRouteRequest) (response *DeleteRouteResponse, err error) {
	response = NewDeleteRouteResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeRoutesRequest() (request *DescribeRoutesRequest) {
	request = &DescribeRoutesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeRoutes")

	return
}

func NewDescribeRoutesResponse() (response *DescribeRoutesResponse) {
	response = &DescribeRoutesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeRoutes(request *DescribeRoutesRequest) (response *DescribeRoutesResponse, err error) {
	response = NewDescribeRoutesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeInstancesRequest() (request *DescribeInstancesRequest) {
	request = &DescribeInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeInstances")

	return
}

func NewDescribeInstancesResponse() (response *DescribeInstancesResponse) {
	response = &DescribeInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeInstances(request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
	response = NewDescribeInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeInstancesStatusRequest() (request *DescribeInstancesStatusRequest) {
	request = &DescribeInstancesStatusRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeInstancesStatus")

	return
}

func NewDescribeInstancesStatusResponse() (response *DescribeInstancesStatusResponse) {
	response = &DescribeInstancesStatusResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeInstancesStatus(request *DescribeInstancesStatusRequest) (response *DescribeInstancesStatusResponse, err error) {
	response = NewDescribeInstancesStatusResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyInstancesResourceGroupRequest() (request *ModifyInstancesResourceGroupRequest) {
	request = &ModifyInstancesResourceGroupRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyInstancesResourceGroup")

	return
}

func NewModifyInstancesResourceGroupResponse() (response *ModifyInstancesResourceGroupResponse) {
	response = &ModifyInstancesResourceGroupResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyInstancesResourceGroup(request *ModifyInstancesResourceGroupRequest) (response *ModifyInstancesResourceGroupResponse, err error) {
	response = NewModifyInstancesResourceGroupResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyInstancesAttributeRequest() (request *ModifyInstancesAttributeRequest) {
	request = &ModifyInstancesAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyInstancesAttribute")

	return
}

func NewModifyInstancesAttributeResponse() (response *ModifyInstancesAttributeResponse) {
	response = &ModifyInstancesAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyInstancesAttribute(request *ModifyInstancesAttributeRequest) (response *ModifyInstancesAttributeResponse, err error) {
	response = NewModifyInstancesAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewRebootInstancesRequest() (request *RebootInstancesRequest) {
	request = &RebootInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "RebootInstances")

	return
}

func NewRebootInstancesResponse() (response *RebootInstancesResponse) {
	response = &RebootInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RebootInstances(request *RebootInstancesRequest) (response *RebootInstancesResponse, err error) {
	response = NewRebootInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewResetInstanceRequest() (request *ResetInstanceRequest) {
	request = &ResetInstanceRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ResetInstance")

	return
}

func NewResetInstanceResponse() (response *ResetInstanceResponse) {
	response = &ResetInstanceResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ResetInstance(request *ResetInstanceRequest) (response *ResetInstanceResponse, err error) {
	response = NewResetInstanceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewResetInstancePasswordRequest() (request *ResetInstancePasswordRequest) {
	request = &ResetInstancePasswordRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ResetInstancePassword")

	return
}

func NewResetInstancePasswordResponse() (response *ResetInstancePasswordResponse) {
	response = &ResetInstancePasswordResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ResetInstancePassword(request *ResetInstancePasswordRequest) (response *ResetInstancePasswordResponse, err error) {
	response = NewResetInstancePasswordResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStartInstancesRequest() (request *StartInstancesRequest) {
	request = &StartInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StartInstances")

	return
}

func NewStartInstancesResponse() (response *StartInstancesResponse) {
	response = &StartInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StartInstances(request *StartInstancesRequest) (response *StartInstancesResponse, err error) {
	response = NewStartInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStopInstancesRequest() (request *StopInstancesRequest) {
	request = &StopInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StopInstances")

	return
}

func NewStopInstancesResponse() (response *StopInstancesResponse) {
	response = &StopInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StopInstances(request *StopInstancesRequest) (response *StopInstancesResponse, err error) {
	response = NewStopInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeCidrsRequest() (request *DescribeCidrsRequest) {
	request = &DescribeCidrsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeCidrs")

	return
}

func NewDescribeCidrsResponse() (response *DescribeCidrsResponse) {
	response = &DescribeCidrsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeCidrs(request *DescribeCidrsRequest) (response *DescribeCidrsResponse, err error) {
	response = NewDescribeCidrsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribePoolsRequest() (request *DescribePoolsRequest) {
	request = &DescribePoolsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribePools")

	return
}

func NewDescribePoolsResponse() (response *DescribePoolsResponse) {
	response = &DescribePoolsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribePools(request *DescribePoolsRequest) (response *DescribePoolsResponse, err error) {
	response = NewDescribePoolsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeCidrRegionsRequest() (request *DescribeCidrRegionsRequest) {
	request = &DescribeCidrRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeCidrRegions")

	return
}

func NewDescribeCidrRegionsResponse() (response *DescribeCidrRegionsResponse) {
	response = &DescribeCidrRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeCidrRegions(request *DescribeCidrRegionsRequest) (response *DescribeCidrRegionsResponse, err error) {
	response = NewDescribeCidrRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeCidrPriceRequest() (request *DescribeCidrPriceRequest) {
	request = &DescribeCidrPriceRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeCidrPrice")

	return
}

func NewDescribeCidrPriceResponse() (response *DescribeCidrPriceResponse) {
	response = &DescribeCidrPriceResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeCidrPrice(request *DescribeCidrPriceRequest) (response *DescribeCidrPriceResponse, err error) {
	response = NewDescribeCidrPriceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateCidrRequest() (request *CreateCidrRequest) {
	request = &CreateCidrRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateCidr")

	return
}

func NewCreateCidrResponse() (response *CreateCidrResponse) {
	response = &CreateCidrResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateCidr(request *CreateCidrRequest) (response *CreateCidrResponse, err error) {
	response = NewCreateCidrResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteCidrRequest() (request *DeleteCidrRequest) {
	request = &DeleteCidrRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteCidr")

	return
}

func NewDeleteCidrResponse() (response *DeleteCidrResponse) {
	response = &DeleteCidrResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteCidr(request *DeleteCidrRequest) (response *DeleteCidrResponse, err error) {
	response = NewDeleteCidrResponse()
	err = c.ApiCall(request, response)
	return
}

func NewRenewCidrRequest() (request *RenewCidrRequest) {
	request = &RenewCidrRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "RenewCidr")

	return
}

func NewRenewCidrResponse() (response *RenewCidrResponse) {
	response = &RenewCidrResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RenewCidr(request *RenewCidrRequest) (response *RenewCidrResponse, err error) {
	response = NewRenewCidrResponse()
	err = c.ApiCall(request, response)
	return
}

func NewAvailableLanIpRequest() (request *AvailableLanIpRequest) {
	request = &AvailableLanIpRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AvailableLanIp")

	return
}

func NewAvailableLanIpResponse() (response *AvailableLanIpResponse) {
	response = &AvailableLanIpResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AvailableLanIp(request *AvailableLanIpRequest) (response *AvailableLanIpResponse, err error) {
	response = NewAvailableLanIpResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeDiskRegionsRequest() (request *DescribeDiskRegionsRequest) {
	request = &DescribeDiskRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeDiskRegions")

	return
}

func NewDescribeDiskRegionsResponse() (response *DescribeDiskRegionsResponse) {
	response = &DescribeDiskRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeDiskRegions(request *DescribeDiskRegionsRequest) (response *DescribeDiskRegionsResponse, err error) {
	response = NewDescribeDiskRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateDisksRequest() (request *CreateDisksRequest) {
	request = &CreateDisksRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateDisks")

	return
}

func NewCreateDisksResponse() (response *CreateDisksResponse) {
	response = &CreateDisksResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateDisks(request *CreateDisksRequest) (response *CreateDisksResponse, err error) {
	response = NewCreateDisksResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeDisksRequest() (request *DescribeDisksRequest) {
	request = &DescribeDisksRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeDisks")

	return
}

func NewDescribeDisksResponse() (response *DescribeDisksResponse) {
	response = &DescribeDisksResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeDisks(request *DescribeDisksRequest) (response *DescribeDisksResponse, err error) {
	response = NewDescribeDisksResponse()
	err = c.ApiCall(request, response)
	return
}

func NewAttachDisksRequest() (request *AttachDisksRequest) {
	request = &AttachDisksRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AttachDisks")

	return
}

func NewAttachDisksResponse() (response *AttachDisksResponse) {
	response = &AttachDisksResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AttachDisks(request *AttachDisksRequest) (response *AttachDisksResponse, err error) {
	response = NewAttachDisksResponse()
	err = c.ApiCall(request, response)
	return
}

func NewChangeDisksAttachRequest() (request *ChangeDisksAttachRequest) {
	request = &ChangeDisksAttachRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ChangeDisksAttach")

	return
}

func NewChangeDisksAttachResponse() (response *ChangeDisksAttachResponse) {
	response = &ChangeDisksAttachResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ChangeDisksAttach(request *ChangeDisksAttachRequest) (response *ChangeDisksAttachResponse, err error) {
	response = NewChangeDisksAttachResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDetachDisksRequest() (request *DetachDisksRequest) {
	request = &DetachDisksRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DetachDisks")

	return
}

func NewDetachDisksResponse() (response *DetachDisksResponse) {
	response = &DetachDisksResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DetachDisks(request *DetachDisksRequest) (response *DetachDisksResponse, err error) {
	response = NewDetachDisksResponse()
	err = c.ApiCall(request, response)
	return
}

func NewResizeDiskRequest() (request *ResizeDiskRequest) {
	request = &ResizeDiskRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ResizeDisk")

	return
}

func NewResizeDiskResponse() (response *ResizeDiskResponse) {
	response = &ResizeDiskResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ResizeDisk 将一个云硬盘扩容到指定大小。
func (c *Client) ResizeDisk(request *ResizeDiskRequest) (response *ResizeDiskResponse, err error) {
	response = NewResizeDiskResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyDisksAttributesRequest() (request *ModifyDisksAttributesRequest) {
	request = &ModifyDisksAttributesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyDisksAttributes")

	return
}

func NewModifyDisksAttributesResponse() (response *ModifyDisksAttributesResponse) {
	response = &ModifyDisksAttributesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyDisksAttributes(request *ModifyDisksAttributesRequest) (response *ModifyDisksAttributesResponse, err error) {
	response = NewModifyDisksAttributesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewInquiryPriceCreateDisksRequest() (request *InquiryPriceCreateDisksRequest) {
	request = &InquiryPriceCreateDisksRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "InquiryPriceCreateDisks")

	return
}

func NewInquiryPriceCreateDisksResponse() (response *InquiryPriceCreateDisksResponse) {
	response = &InquiryPriceCreateDisksResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) InquiryPriceCreateDisks(request *InquiryPriceCreateDisksRequest) (response *InquiryPriceCreateDisksResponse, err error) {
	response = NewInquiryPriceCreateDisksResponse()
	err = c.ApiCall(request, response)
	return
}

func NewTerminateDiskRequest() (request *TerminateDiskRequest) {
	request = &TerminateDiskRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "TerminateDisk")

	return
}

func NewTerminateDiskResponse() (response *TerminateDiskResponse) {
	response = &TerminateDiskResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) TerminateDisk(request *TerminateDiskRequest) (response *TerminateDiskResponse, err error) {
	response = NewTerminateDiskResponse()
	err = c.ApiCall(request, response)
	return
}

func NewReleaseDiskRequest() (request *ReleaseDiskRequest) {
	request = &ReleaseDiskRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ReleaseDisk")

	return
}

func NewReleaseDiskResponse() (response *ReleaseDiskResponse) {
	response = &ReleaseDiskResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ReleaseDisk(request *ReleaseDiskRequest) (response *ReleaseDiskResponse, err error) {
	response = NewReleaseDiskResponse()
	err = c.ApiCall(request, response)
	return
}

func NewRenewDiskRequest() (request *RenewDiskRequest) {
	request = &RenewDiskRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "RenewDisk")

	return
}

func NewRenewDiskResponse() (response *RenewDiskResponse) {
	response = &RenewDiskResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RenewDisk(request *RenewDiskRequest) (response *RenewDiskResponse, err error) {
	response = NewRenewDiskResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyDisksResourceGroupRequest() (request *ModifyDisksResourceGroupRequest) {
	request = &ModifyDisksResourceGroupRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyDisksResourceGroup")

	return
}

func NewModifyDisksResourceGroupResponse() (response *ModifyDisksResourceGroupResponse) {
	response = &ModifyDisksResourceGroupResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyDisksResourceGroup(request *ModifyDisksResourceGroupRequest) (response *ModifyDisksResourceGroupResponse, err error) {
	response = NewModifyDisksResourceGroupResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeDiskCategoryRequest() (request *DescribeDiskCategoryRequest) {
	request = &DescribeDiskCategoryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeDiskCategory")

	return
}

func NewDescribeDiskCategoryResponse() (response *DescribeDiskCategoryResponse) {
	response = &DescribeDiskCategoryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeDiskCategory(request *DescribeDiskCategoryRequest) (response *DescribeDiskCategoryResponse, err error) {
	response = NewDescribeDiskCategoryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeSecurityGroupsRequest() (request *DescribeSecurityGroupsRequest) {
	request = &DescribeSecurityGroupsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeSecurityGroups")

	return
}

func NewDescribeSecurityGroupsResponse() (response *DescribeSecurityGroupsResponse) {
	response = &DescribeSecurityGroupsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeSecurityGroups(request *DescribeSecurityGroupsRequest) (response *DescribeSecurityGroupsResponse, err error) {
	response = NewDescribeSecurityGroupsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifySecurityGroupsAttributeRequest() (request *ModifySecurityGroupsAttributeRequest) {
	request = &ModifySecurityGroupsAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifySecurityGroupsAttribute")

	return
}

func NewModifySecurityGroupsAttributeResponse() (response *ModifySecurityGroupsAttributeResponse) {
	response = &ModifySecurityGroupsAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifySecurityGroupsAttribute(request *ModifySecurityGroupsAttributeRequest) (response *ModifySecurityGroupsAttributeResponse, err error) {
	response = NewModifySecurityGroupsAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeSecurityGroupRuleRequest() (request *DescribeSecurityGroupRuleRequest) {
	request = &DescribeSecurityGroupRuleRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeSecurityGroupRule")

	return
}

func NewDescribeSecurityGroupRuleResponse() (response *DescribeSecurityGroupRuleResponse) {
	response = &DescribeSecurityGroupRuleResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeSecurityGroupRule(request *DescribeSecurityGroupRuleRequest) (response *DescribeSecurityGroupRuleResponse, err error) {
	response = NewDescribeSecurityGroupRuleResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateSecurityGroupRequest() (request *CreateSecurityGroupRequest) {
	request = &CreateSecurityGroupRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateSecurityGroup")

	return
}

func NewCreateSecurityGroupResponse() (response *CreateSecurityGroupResponse) {
	response = &CreateSecurityGroupResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateSecurityGroup(request *CreateSecurityGroupRequest) (response *CreateSecurityGroupResponse, err error) {
	response = NewCreateSecurityGroupResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteSecurityGroupRequest() (request *DeleteSecurityGroupRequest) {
	request = &DeleteSecurityGroupRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteSecurityGroup")

	return
}

func NewDeleteSecurityGroupResponse() (response *DeleteSecurityGroupResponse) {
	response = &DeleteSecurityGroupResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteSecurityGroup(request *DeleteSecurityGroupRequest) (response *DeleteSecurityGroupResponse, err error) {
	response = NewDeleteSecurityGroupResponse()
	err = c.ApiCall(request, response)
	return
}

func NewConfigureSecurityGroupRulesRequest() (request *ConfigureSecurityGroupRulesRequest) {
	request = &ConfigureSecurityGroupRulesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ConfigureSecurityGroupRules")

	return
}

func NewConfigureSecurityGroupRulesResponse() (response *ConfigureSecurityGroupRulesResponse) {
	response = &ConfigureSecurityGroupRulesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ConfigureSecurityGroupRules(request *ConfigureSecurityGroupRulesRequest) (response *ConfigureSecurityGroupRulesResponse, err error) {
	response = NewConfigureSecurityGroupRulesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewAssignSecurityGroupVpcRequest() (request *AssignSecurityGroupVpcRequest) {
	request = &AssignSecurityGroupVpcRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AssignSecurityGroupVpc")

	return
}

func NewAssignSecurityGroupVpcResponse() (response *AssignSecurityGroupVpcResponse) {
	response = &AssignSecurityGroupVpcResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AssignSecurityGroupVpc(request *AssignSecurityGroupVpcRequest) (response *AssignSecurityGroupVpcResponse, err error) {
	response = NewAssignSecurityGroupVpcResponse()
	err = c.ApiCall(request, response)
	return
}

func NewUnAssignSecurityGroupVpcRequest() (request *UnAssignSecurityGroupVpcRequest) {
	request = &UnAssignSecurityGroupVpcRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "UnAssignSecurityGroupVpc")

	return
}

func NewUnAssignSecurityGroupVpcResponse() (response *UnAssignSecurityGroupVpcResponse) {
	response = &UnAssignSecurityGroupVpcResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) UnAssignSecurityGroupVpc(request *UnAssignSecurityGroupVpcRequest) (response *UnAssignSecurityGroupVpcResponse, err error) {
	response = NewUnAssignSecurityGroupVpcResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeNetworkInterfacesRequest() (request *DescribeNicsRequest) {
	request = &DescribeNicsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeNetworkInterfaces")

	return
}

func NewDescribeNetworkInterfacesResponse() (response *DescribeNicsResponse) {
	response = &DescribeNicsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeNetworkInterfaces(request *DescribeNicsRequest) (response *DescribeNicsResponse, err error) {
	response = NewDescribeNetworkInterfacesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyNetworkInterfacesAttributeRequest() (request *ModifyNicsAttributeRequest) {
	request = &ModifyNicsAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyNetworkInterfacesAttribute")

	return
}

func NewModifyNetworkInterfacesAttributeResponse() (response *ModifyNicsAttributeResponse) {
	response = &ModifyNicsAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyNetworkInterfacesAttribute(request *ModifyNicsAttributeRequest) (response *ModifyNicsAttributeResponse, err error) {
	response = NewModifyNetworkInterfacesAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateNetworkInterfaceRequest() (request *CreateNicRequest) {
	request = &CreateNicRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateNetworkInterface")

	return
}

func NewCreateNetworkInterfaceResponse() (response *CreateNicResponse) {
	response = &CreateNicResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateNetworkInterface(request *CreateNicRequest) (response *CreateNicResponse, err error) {
	response = NewCreateNetworkInterfaceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteNetworkInterfaceRequest() (request *DeleteNicRequest) {
	request = &DeleteNicRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteNetworkInterface")

	return
}

func NewDeleteNetworkInterfaceResponse() (response *DeleteNicResponse) {
	response = &DeleteNicResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteNetworkInterface(request *DeleteNicRequest) (response *DeleteNicResponse, err error) {
	response = NewDeleteNetworkInterfaceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewAttachNetworkInterfaceRequest() (request *AttachNicRequest) {
	request = &AttachNicRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AttachNetworkInterface")

	return
}

func NewAttachNetworkInterfaceResponse() (response *AttachNicResponse) {
	response = &AttachNicResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AttachNetworkInterface(request *AttachNicRequest) (response *AttachNicResponse, err error) {
	response = NewAttachNetworkInterfaceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewAssignNetworkInterfaceIpv6Request() (request *AssignNicIpv6Request) {
	request = &AssignNicIpv6Request{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AssignNetworkInterfaceIpv6")

	return
}

func NewAssignNetworkInterfaceIpv6Response() (response *AssignNicIpv6Response) {
	response = &AssignNicIpv6Response{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AssignNetworkInterfaceIpv6(request *AssignNicIpv6Request) (response *AssignNicIpv6Response, err error) {
	response = NewAssignNetworkInterfaceIpv6Response()
	err = c.ApiCall(request, response)
	return
}

func NewUnassignNetworkInterfaceIpv4Request() (request *UnAssignNicIpv4Request) {
	request = &UnAssignNicIpv4Request{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "UnassignNetworkInterfaceIpv4")

	return
}

func NewUnassignNetworkInterfaceIpv4Response() (response *UnAssignNicIpv4Response) {
	response = &UnAssignNicIpv4Response{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) UnassignNetworkInterfaceIpv4(request *UnAssignNicIpv4Request) (response *UnAssignNicIpv4Response, err error) {
	response = NewUnassignNetworkInterfaceIpv4Response()
	err = c.ApiCall(request, response)
	return
}

func NewAssignNetworkInterfaceIpv4Request() (request *AssignNicIpv4Request) {
	request = &AssignNicIpv4Request{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AssignNetworkInterfaceIpv4")

	return
}

func NewAssignNetworkInterfaceIpv4Response() (response *AssignNicIpv4Response) {
	response = &AssignNicIpv4Response{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AssignNetworkInterfaceIpv4(request *AssignNicIpv4Request) (response *AssignNicIpv4Response, err error) {
	response = NewAssignNetworkInterfaceIpv4Response()
	err = c.ApiCall(request, response)
	return
}

func NewBatchAssignNetworkInterfaceIpv4Request() (request *BatchAssignNicIpv4Request) {
	request = &BatchAssignNicIpv4Request{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "BatchAssignNetworkInterfaceIpv4")

	return
}

func NewBatchAssignNetworkInterfaceIpv4Response() (response *BatchAssignNicIpv4Response) {
	response = &BatchAssignNicIpv4Response{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) BatchAssignNetworkInterfaceIpv4(request *BatchAssignNicIpv4Request) (response *BatchAssignNicIpv4Response, err error) {
	response = NewBatchAssignNetworkInterfaceIpv4Response()
	err = c.ApiCall(request, response)
	return
}

func NewDetachNetworkInterfaceRequest() (request *DetachNicRequest) {
	request = &DetachNicRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DetachNetworkInterface")

	return
}

func NewDetachNetworkInterfaceResponse() (response *DetachNicResponse) {
	response = &DetachNicResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DetachNetworkInterface(request *DetachNicRequest) (response *DetachNicResponse, err error) {
	response = NewDetachNetworkInterfaceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeNetworkInterfaceRegionsRequest() (request *DescribeNicRegionsRequest) {
	request = &DescribeNicRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeNetworkInterfaceRegions")

	return
}

func NewDescribeNetworkInterfaceRegionsResponse() (response *DescribeNicRegionsResponse) {
	response = &DescribeNicRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeNetworkInterfaceRegions(request *DescribeNicRegionsRequest) (response *DescribeNicRegionsResponse, err error) {
	response = NewDescribeNetworkInterfaceRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewInquiryPricePublicIpv6Request() (request *InquiryPricePublicIpv6Request) {
	request = &InquiryPricePublicIpv6Request{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "InquiryPricePublicIpv6")

	return
}

func NewInquiryPricePublicIpv6Response() (response *InquiryPricePublicIpv6Response) {
	response = &InquiryPricePublicIpv6Response{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) InquiryPricePublicIpv6(request *InquiryPricePublicIpv6Request) (response *InquiryPricePublicIpv6Response, err error) {
	response = NewInquiryPricePublicIpv6Response()
	err = c.ApiCall(request, response)
	return
}

func NewCreateZecInstancesRequest() (request *CreateZecInstancesRequest) {
	request = &CreateZecInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateZecInstances")

	return
}

func NewCreateZecInstancesResponse() (response *CreateZecInstancesResponse) {
	response = &CreateZecInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateZecInstances(request *CreateZecInstancesRequest) (response *CreateZecInstancesResponse, err error) {
	response = NewCreateZecInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateInstancesRequest() (request *CreateInstancesRequest) {
	request = &CreateInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateInstances")

	return
}

func NewCreateInstancesResponse() (response *CreateInstancesResponse) {
	response = &CreateInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateInstances(request *CreateInstancesRequest) (response *CreateInstancesResponse, err error) {
	response = NewCreateInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeZoneInstanceConfigInfosRequest() (request *DescribeZoneInstanceConfigInfosRequest) {
	request = &DescribeZoneInstanceConfigInfosRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeZoneInstanceConfigInfos")

	return
}

func NewDescribeZoneInstanceConfigInfosResponse() (response *DescribeZoneInstanceConfigInfosResponse) {
	response = &DescribeZoneInstanceConfigInfosResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeZoneInstanceConfigInfos(request *DescribeZoneInstanceConfigInfosRequest) (response *DescribeZoneInstanceConfigInfosResponse, err error) {
	response = NewDescribeZoneInstanceConfigInfosResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeZonesRequest() (request *DescribeZonesRequest) {
	request = &DescribeZonesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeZones")

	return
}

func NewDescribeZonesResponse() (response *DescribeZonesResponse) {
	response = &DescribeZonesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeZones(request *DescribeZonesRequest) (response *DescribeZonesResponse, err error) {
	response = NewDescribeZonesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeImagesRequest() (request *DescribeImagesRequest) {
	request = &DescribeImagesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeImages")

	return
}

func NewDescribeImagesResponse() (response *DescribeImagesResponse) {
	response = &DescribeImagesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeImages(request *DescribeImagesRequest) (response *DescribeImagesResponse, err error) {
	response = NewDescribeImagesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeKeyPairsRequest() (request *DescribeKeyPairsRequest) {
	request = &DescribeKeyPairsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeKeyPairs")

	return
}

func NewDescribeKeyPairsResponse() (response *DescribeKeyPairsResponse) {
	response = &DescribeKeyPairsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeKeyPairs(request *DescribeKeyPairsRequest) (response *DescribeKeyPairsResponse, err error) {
	response = NewDescribeKeyPairsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipRegionsRequest() (request *DescribeEipRegionsRequest) {
	request = &DescribeEipRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipRegions")

	return
}

func NewDescribeEipRegionsResponse() (response *DescribeEipRegionsResponse) {
	response = &DescribeEipRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeEipRegions(request *DescribeEipRegionsRequest) (response *DescribeEipRegionsResponse, err error) {
	response = NewDescribeEipRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipRemoteRegionsRequest() (request *DescribeEipRemoteRegionsRequest) {
	request = &DescribeEipRemoteRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipRemoteRegions")

	return
}

func NewDescribeEipRemoteRegionsResponse() (response *DescribeEipRemoteRegionsResponse) {
	response = &DescribeEipRemoteRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeEipRemoteRegions(request *DescribeEipRemoteRegionsRequest) (response *DescribeEipRemoteRegionsResponse, err error) {
	response = NewDescribeEipRemoteRegionsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipInternetChargeTypesRequest() (request *DescribeEipInternetChargeTypesRequest) {
	request = &DescribeEipInternetChargeTypesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipInternetChargeTypes")

	return
}

func NewDescribeEipInternetChargeTypesResponse() (response *DescribeEipInternetChargeTypesResponse) {
	response = &DescribeEipInternetChargeTypesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeEipInternetChargeTypes(request *DescribeEipInternetChargeTypesRequest) (response *DescribeEipInternetChargeTypesResponse, err error) {
	response = NewDescribeEipInternetChargeTypesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipsRequest() (request *DescribeEipsRequest) {
	request = &DescribeEipsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEips")

	return
}

func NewDescribeEipsResponse() (response *DescribeEipsResponse) {
	response = &DescribeEipsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeEips(request *DescribeEipsRequest) (response *DescribeEipsResponse, err error) {
	response = NewDescribeEipsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateEipsRequest() (request *CreateEipsRequest) {
	request = &CreateEipsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateEips")

	return
}

func NewCreateEipsResponse() (response *CreateEipsResponse) {
	response = &CreateEipsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateEips(request *CreateEipsRequest) (response *CreateEipsResponse, err error) {
	response = NewCreateEipsResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteEipRequest() (request *DeleteEipRequest) {
	request = &DeleteEipRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteEip")

	return
}

func NewDeleteEipResponse() (response *DeleteEipResponse) {
	response = &DeleteEipResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteEip(request *DeleteEipRequest) (response *DeleteEipResponse, err error) {
	response = NewDeleteEipResponse()
	err = c.ApiCall(request, response)
	return
}

func NewRenewEipRequest() (request *RenewEipRequest) {
	request = &RenewEipRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "RenewEip")

	return
}

func NewRenewEipResponse() (response *RenewEipResponse) {
	response = &RenewEipResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RenewEip(request *RenewEipRequest) (response *RenewEipResponse, err error) {
	response = NewRenewEipResponse()
	err = c.ApiCall(request, response)
	return
}

func NewBatchAttachEipLanIpRequest() (request *BatchAttachEipLanIpRequest) {
	request = &BatchAttachEipLanIpRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "BatchAttachEipLanIp")

	return
}

func NewBatchAttachEipLanIpResponse() (response *BatchAttachEipLanIpResponse) {
	response = &BatchAttachEipLanIpResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) BatchAttachEipLanIp(request *BatchAttachEipLanIpRequest) (response *BatchAttachEipLanIpResponse, err error) {
	response = NewBatchAttachEipLanIpResponse()
	err = c.ApiCall(request, response)
	return
}



func NewAssociateEipAddressRequest() (request *AssociateEipAddressRequest) {
    request = &AssociateEipAddressRequest{
        BaseRequest: &common.BaseRequest{},
    }
    request.Init().InitWithApiInfo(SERVICE, APIVersion, "AssociateEipAddress")
    return
}

func NewAssociateEipAddressResponse() (response *AssociateEipAddressResponse) {
    response = &AssociateEipAddressResponse{
        BaseResponse: &common.BaseResponse{},
    }
    return
}

func (c *Client) AssociateEipAddress(request *AssociateEipAddressRequest) (response *AssociateEipAddressResponse, err error) {
    response = NewAssociateEipAddressResponse()
    err = c.ApiCall(request, response)
    return
}


func NewDetachEipLanIpRequest() (request *DetachEipLanIpRequest) {
	request = &DetachEipLanIpRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DetachEipLanIp")

	return
}

func NewDetachEipLanIpResponse() (response *DetachEipLanIpResponse) {
	response = &DetachEipLanIpResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DetachEipLanIp(request *DetachEipLanIpRequest) (response *DetachEipLanIpResponse, err error) {
	response = NewDetachEipLanIpResponse()
	err = c.ApiCall(request, response)
	return
}


func (c *Client) UnassociateEipAddress(request *UnassociateEipAddressRequest) (response *UnassociateEipAddressResponse, err error) {
    response = NewUnassociateEipAddressResponse()
    err = c.ApiCall(request, response)
    return
}

func NewUnassociateEipAddressRequest() (request *UnassociateEipAddressRequest) {
    request = &UnassociateEipAddressRequest{
        BaseRequest: &common.BaseRequest{},
    }
    request.Init().InitWithApiInfo(SERVICE, APIVersion, "UnassociateEipAddress")
    return
}

func NewUnassociateEipAddressResponse() (response *UnassociateEipAddressResponse) {
    response = &UnassociateEipAddressResponse{
        BaseResponse: &common.BaseResponse{},
    }
    return
}


func NewReplaceEipAddressRequest() (request *ReplaceEipAddressRequest) {
    request = &ReplaceEipAddressRequest{
        BaseRequest: &common.BaseRequest{},
    }
    request.Init().InitWithApiInfo(SERVICE, APIVersion, "ReplaceEipAddress")
    return
}

func NewReplaceEipAddressResponse() (response *ReplaceEipAddressResponse) {
    response = &ReplaceEipAddressResponse{
        BaseResponse: &common.BaseResponse{},
    }
    return
}

func (c *Client) ReplaceEipAddress(request *ReplaceEipAddressRequest) (response *ReplaceEipAddressResponse, err error) {
    response = NewReplaceEipAddressResponse()
    err = c.ApiCall(request, response)
    return
}


func NewConfigEipEgressIpRequest() (request *ConfigEipEgressIpRequest) {
	request = &ConfigEipEgressIpRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ConfigEipEgressIp")

	return
}

func NewConfigEipEgressIpResponse() (response *ConfigEipEgressIpResponse) {
	response = &ConfigEipEgressIpResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ConfigEipEgressIp(request *ConfigEipEgressIpRequest) (response *ConfigEipEgressIpResponse, err error) {
	response = NewConfigEipEgressIpResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipPriceRequest() (request *DescribeEipPriceRequest) {
	request = &DescribeEipPriceRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipPrice")

	return
}

func NewDescribeEipPriceResponse() (response *DescribeEipPriceResponse) {
	response = &DescribeEipPriceResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeEipPrice(request *DescribeEipPriceRequest) (response *DescribeEipPriceResponse, err error) {
	response = NewDescribeEipPriceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewChangeEipInternetChargeTypeRequest() (request *ChangeEipInternetChargeTypeRequest) {
	request = &ChangeEipInternetChargeTypeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ChangeEipInternetChargeType")

	return
}

func NewChangeEipInternetChargeTypeResponse() (response *ChangeEipInternetChargeTypeResponse) {
	response = &ChangeEipInternetChargeTypeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ChangeEipInternetChargeType(request *ChangeEipInternetChargeTypeRequest) (response *ChangeEipInternetChargeTypeResponse, err error) {
	response = NewChangeEipInternetChargeTypeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewInquiryPriceCreateInstanceRequest() (request *InquiryPriceCreateInstanceRequest) {
	request = &InquiryPriceCreateInstanceRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "InquiryPriceCreateInstance")

	return
}

func NewInquiryPriceCreateInstanceResponse() (response *InquiryPriceCreateInstanceResponse) {
	response = &InquiryPriceCreateInstanceResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) InquiryPriceCreateInstance(request *InquiryPriceCreateInstanceRequest) (response *InquiryPriceCreateInstanceResponse, err error) {
	response = NewInquiryPriceCreateInstanceResponse()
	err = c.ApiCall(request, response)
	return
}

func NewReleaseInstancesRequest() (request *ReleaseInstancesRequest) {
	request = &ReleaseInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ReleaseInstances")

	return
}

func NewReleaseInstancesResponse() (response *ReleaseInstancesResponse) {
	response = &ReleaseInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ReleaseInstances(request *ReleaseInstancesRequest) (response *ReleaseInstancesResponse, err error) {
	response = NewReleaseInstancesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyInstanceTypeRequest() (request *ModifyInstanceTypeRequest) {
	request = &ModifyInstanceTypeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyInstanceType")

	return
}

func NewModifyInstanceTypeResponse() (response *ModifyInstanceTypeResponse) {
	response = &ModifyInstanceTypeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyInstanceType(request *ModifyInstanceTypeRequest) (response *ModifyInstanceTypeResponse, err error) {
	response = NewModifyInstanceTypeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeTimeZonesRequest() (request *DescribeTimeZonesRequest) {
	request = &DescribeTimeZonesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeTimeZones")

	return
}

func NewDescribeTimeZonesResponse() (response *DescribeTimeZonesResponse) {
	response = &DescribeTimeZonesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeTimeZones(request *DescribeTimeZonesRequest) (response *DescribeTimeZonesResponse, err error) {
	response = NewDescribeTimeZonesResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStartIpForwardRequest() (request *StartIpForwardRequest) {
	request = &StartIpForwardRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StartIpForward")

	return
}

func NewStartIpForwardResponse() (response *StartIpForwardResponse) {
	response = &StartIpForwardResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StartIpForward(request *StartIpForwardRequest) (response *StartIpForwardResponse, err error) {
	response = NewStartIpForwardResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStopIpForwardRequest() (request *StopIpForwardRequest) {
	request = &StopIpForwardRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StopIpForward")

	return
}

func NewStopIpForwardResponse() (response *StopIpForwardResponse) {
	response = &StopIpForwardResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StopIpForward(request *StopIpForwardRequest) (response *StopIpForwardResponse, err error) {
	response = NewStopIpForwardResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStartAgentMonitorRequest() (request *StartAgentMonitorRequest) {
	request = &StartAgentMonitorRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StartAgentMonitor")

	return
}

func NewStartAgentMonitorResponse() (response *StartAgentMonitorResponse) {
	response = &StartAgentMonitorResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StartAgentMonitor(request *StartAgentMonitorRequest) (response *StartAgentMonitorResponse, err error) {
	response = NewStartAgentMonitorResponse()
	err = c.ApiCall(request, response)
	return
}

func NewStopAgentMonitorRequest() (request *StopAgentMonitorRequest) {
	request = &StopAgentMonitorRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "StopAgentMonitor")

	return
}

func NewStopAgentMonitorResponse() (response *StopAgentMonitorResponse) {
	response = &StopAgentMonitorResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) StopAgentMonitor(request *StopAgentMonitorRequest) (response *StopAgentMonitorResponse, err error) {
	response = NewStopAgentMonitorResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeVncUrlRequest() (request *DescribeVncUrlRequest) {
	request = &DescribeVncUrlRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeVncUrl")

	return
}

func NewDescribeVncUrlResponse() (response *DescribeVncUrlResponse) {
	response = &DescribeVncUrlResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVncUrl(request *DescribeVncUrlRequest) (response *DescribeVncUrlResponse, err error) {
	response = NewDescribeVncUrlResponse()
	err = c.ApiCall(request, response)
	return
}

func NewChangeNicNetworkTypeRequest() (request *ChangeNicNetworkTypeRequest) {
	request = &ChangeNicNetworkTypeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ChangeNicNetworkType")

	return
}

func NewChangeNicNetworkTypeResponse() (response *ChangeNicNetworkTypeResponse) {
	response = &ChangeNicNetworkTypeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ChangeNicNetworkType(request *ChangeNicNetworkTypeRequest) (response *ChangeNicNetworkTypeResponse, err error) {
	response = NewChangeNicNetworkTypeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipTrafficRequest() (request *DescribeEipTrafficRequest) {
	request = &DescribeEipTrafficRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipTraffic")

	return
}

func NewDescribeEipTrafficResponse() (response *DescribeEipTrafficResponse) {
	response = &DescribeEipTrafficResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeEipTraffic
func (c *Client) DescribeEipTraffic(request *DescribeEipTrafficRequest) (response *DescribeEipTrafficResponse, err error) {
	response = NewDescribeEipTrafficResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateBorderGatewayRequest() (request *CreateBorderGatewayRequest) {
	request = &CreateBorderGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateBorderGateway")

	return
}

func NewCreateBorderGatewayResponse() (response *CreateBorderGatewayResponse) {
	response = &CreateBorderGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// CreateBorderGateway
func (c *Client) CreateBorderGateway(request *CreateBorderGatewayRequest) (response *CreateBorderGatewayResponse, err error) {
	response = NewCreateBorderGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteBorderGatewayRequest() (request *DeleteBorderGatewayRequest) {
	request = &DeleteBorderGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteBorderGateway")

	return
}

func NewDeleteBorderGatewayResponse() (response *DeleteBorderGatewayResponse) {
	response = &DeleteBorderGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DeleteBorderGateway
func (c *Client) DeleteBorderGateway(request *DeleteBorderGatewayRequest) (response *DeleteBorderGatewayResponse, err error) {
	response = NewDeleteBorderGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeBorderGatewaysRequest() (request *DescribeBorderGatewaysRequest) {
	request = &DescribeBorderGatewaysRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeBorderGateways")

	return
}

func NewDescribeBorderGatewaysResponse() (response *DescribeBorderGatewaysResponse) {
	response = &DescribeBorderGatewaysResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeBorderGateways
func (c *Client) DescribeBorderGateways(request *DescribeBorderGatewaysRequest) (response *DescribeBorderGatewaysResponse, err error) {
	response = NewDescribeBorderGatewaysResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyBorderGatewaysAttributeRequest() (request *ModifyBorderGatewaysAttributeRequest) {
	request = &ModifyBorderGatewaysAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyBorderGatewaysAttribute")

	return
}

func NewModifyBorderGatewaysAttributeResponse() (response *ModifyBorderGatewaysAttributeResponse) {
	response = &ModifyBorderGatewaysAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ModifyBorderGatewaysAttribute
func (c *Client) ModifyBorderGatewaysAttribute(request *ModifyBorderGatewaysAttributeRequest) (response *ModifyBorderGatewaysAttributeResponse, err error) {
	response = NewModifyBorderGatewaysAttributeResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyBorderGatewayAsnRequest() (request *ModifyBorderGatewayAsnRequest) {
	request = &ModifyBorderGatewayAsnRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyBorderGatewayAsn")

	return
}

func NewModifyBorderGatewayAsnResponse() (response *ModifyBorderGatewayAsnResponse) {
	response = &ModifyBorderGatewayAsnResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ModifyBorderGatewayAsn
func (c *Client) ModifyBorderGatewayAsn(request *ModifyBorderGatewayAsnRequest) (response *ModifyBorderGatewayAsnResponse, err error) {
	response = NewModifyBorderGatewayAsnResponse()
	err = c.ApiCall(request, response)
	return
}



func NewDescribeAvailableNatsRequest() (request *DescribeAvailableNatsRequest) {
	request = &DescribeAvailableNatsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeAvailableNats")

	return
}

func NewDescribeAvailableNatsResponse() (response *DescribeAvailableNatsResponse) {
	response = &DescribeAvailableNatsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeAvailableNats
func (c *Client) DescribeAvailableNats(request *DescribeAvailableNatsRequest) (response *DescribeAvailableNatsResponse, err error) {
	response = NewDescribeAvailableNatsResponse()
	err = c.ApiCall(request, response)
	return
}


func NewAssignBorderGatewayRequest() (request *AssignBorderGatewayRequest) {
	request = &AssignBorderGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "AssignBorderGateway")

	return
}

func NewAssignBorderGatewayResponse() (response *AssignBorderGatewayResponse) {
	response = &AssignBorderGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}


// AssignBorderGateway
func (c *Client) AssignBorderGateway(request *AssignBorderGatewayRequest) (response *AssignBorderGatewayResponse, err error) {
	response = NewAssignBorderGatewayResponse()
	err = c.ApiCall(request, response)
	return
}



func NewUnassignBorderGatewayRequest() (request *UnassignBorderGatewayRequest) {
	request = &UnassignBorderGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "UnassignBorderGateway")

	return
}

func NewUnassignBorderGatewayResponse() (response *UnassignBorderGatewayResponse) {
	response = &UnassignBorderGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}


// UnassignBorderGateway
func (c *Client) UnassignBorderGateway(request *UnassignBorderGatewayRequest) (response *UnassignBorderGatewayResponse, err error) {
	response = NewUnassignBorderGatewayResponse()
	err = c.ApiCall(request, response)
	return
}


func NewCreateNatGatewayRequest() (request *CreateNatGatewayRequest) {
	request = &CreateNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateNatGateway")

	return
}

func NewCreateNatGatewayResponse() (response *CreateNatGatewayResponse) {
	response = &CreateNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// CreateNatGateway 创建NAT网关。
func (c *Client) CreateNatGateway(request *CreateNatGatewayRequest) (response *CreateNatGatewayResponse, err error) {
	response = NewCreateNatGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeNatGatewaysRequest() (request *DescribeNatGatewaysRequest) {
	request = &DescribeNatGatewaysRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeNatGateways")

	return
}

func NewDescribeNatGatewaysResponse() (response *DescribeNatGatewaysResponse) {
	response = &DescribeNatGatewaysResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeNatGateways 查询NAT网关列表。
func (c *Client) DescribeNatGateways(request *DescribeNatGatewaysRequest) (response *DescribeNatGatewaysResponse, err error) {
	response = NewDescribeNatGatewaysResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeNatGatewayDetailRequest() (request *DescribeNatGatewayDetailRequest) {
	request = &DescribeNatGatewayDetailRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeNatGatewayDetail")

	return
}

func NewDescribeNatGatewayDetailResponse() (response *DescribeNatGatewayDetailResponse) {
	response = &DescribeNatGatewayDetailResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeNatGatewayDetail 查询NAT网关详情规则表。
func (c *Client) DescribeNatGatewayDetail(request *DescribeNatGatewayDetailRequest) (response *DescribeNatGatewayDetailResponse, err error) {
	response = NewDescribeNatGatewayDetailResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteNatGatewayRequest() (request *DeleteNatGatewayRequest) {
	request = &DeleteNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteNatGateway")

	return
}

func NewDeleteNatGatewayResponse() (response *DeleteNatGatewayResponse) {
	response = &DeleteNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DeleteNatGateway 删除一个指定的NAT网关。
func (c *Client) DeleteNatGateway(request *DeleteNatGatewayRequest) (response *DeleteNatGatewayResponse, err error) {
	response = NewDeleteNatGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewRenewNatGatewayRequest() (request *RenewNatGatewayRequest) {
	request = &RenewNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "RenewNatGateway")

	return
}

func NewRenewNatGatewayResponse() (response *RenewNatGatewayResponse) {
	response = &RenewNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// RenewNatGateway 续费一个指定的NAT网关。
func (c *Client) RenewNatGateway(request *RenewNatGatewayRequest) (response *RenewNatGatewayResponse, err error) {
	response = NewRenewNatGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewInquiryPriceCreateNatGatewayRequest() (request *InquiryPriceCreateNatGatewayRequest) {
	request = &InquiryPriceCreateNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "InquiryPriceCreateNatGateway")

	return
}

func NewInquiryPriceCreateNatGatewayResponse() (response *InquiryPriceCreateNatGatewayResponse) {
	response = &InquiryPriceCreateNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// InquiryPriceCreateNatGateway 查询创建NAT网关的价格。
func (c *Client) InquiryPriceCreateNatGateway(request *InquiryPriceCreateNatGatewayRequest) (response *InquiryPriceCreateNatGatewayResponse, err error) {
	response = NewInquiryPriceCreateNatGatewayResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateSnatEntryRequest() (request *CreateSnatEntryRequest) {
	request = &CreateSnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateSnatEntry")

	return
}

func NewCreateSnatEntryResponse() (response *CreateSnatEntryResponse) {
	response = &CreateSnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// CreateSnatEntry 创建SNAT规则。
func (c *Client) CreateSnatEntry(request *CreateSnatEntryRequest) (response *CreateSnatEntryResponse, err error) {
	response = NewCreateSnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifySnatEntryRequest() (request *ModifySnatEntryRequest) {
	request = &ModifySnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifySnatEntry")

	return
}

func NewModifySnatEntryResponse() (response *ModifySnatEntryResponse) {
	response = &ModifySnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ModifySnatEntry 修改SNAT规则。
func (c *Client) ModifySnatEntry(request *ModifySnatEntryRequest) (response *ModifySnatEntryResponse, err error) {
	response = NewModifySnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteSnatEntryRequest() (request *DeleteSnatEntryRequest) {
	request = &DeleteSnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteSnatEntry")

	return
}

func NewDeleteSnatEntryResponse() (response *DeleteSnatEntryResponse) {
	response = &DeleteSnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DeleteSnatEntry 删除SNAT规则。
func (c *Client) DeleteSnatEntry(request *DeleteSnatEntryRequest) (response *DeleteSnatEntryResponse, err error) {
	response = NewDeleteSnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewCreateDnatEntryRequest() (request *CreateDnatEntryRequest) {
	request = &CreateDnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "CreateDnatEntry")

	return
}

func NewCreateDnatEntryResponse() (response *CreateDnatEntryResponse) {
	response = &CreateDnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// CreateDnatEntry 创建DNAT规则。
func (c *Client) CreateDnatEntry(request *CreateDnatEntryRequest) (response *CreateDnatEntryResponse, err error) {
	response = NewCreateDnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewModifyDnatEntryRequest() (request *ModifyDnatEntryRequest) {
	request = &ModifyDnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyDnatEntry")

	return
}

func NewModifyDnatEntryResponse() (response *ModifyDnatEntryResponse) {
	response = &ModifyDnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ModifyDnatEntry 修改DNAT规则。
func (c *Client) ModifyDnatEntry(request *ModifyDnatEntryRequest) (response *ModifyDnatEntryResponse, err error) {
	response = NewModifyDnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDeleteDnatEntryRequest() (request *DeleteDnatEntryRequest) {
	request = &DeleteDnatEntryRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DeleteDnatEntry")

	return
}

func NewDeleteDnatEntryResponse() (response *DeleteDnatEntryResponse) {
	response = &DeleteDnatEntryResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DeleteDnatEntry 删除DNAT规则。
func (c *Client) DeleteDnatEntry(request *DeleteDnatEntryRequest) (response *DeleteDnatEntryResponse, err error) {
	response = NewDeleteDnatEntryResponse()
	err = c.ApiCall(request, response)
	return
}



func NewDescribeNatGatewayRegionsRequest() (request *DescribeNatGatewayRegionsRequest) {
	request = &DescribeNatGatewayRegionsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeNatGatewayRegions")

	return
}

func NewDescribeNatGatewayRegionsResponse() (response *DescribeNatGatewayRegionsResponse) {
	response = &DescribeNatGatewayRegionsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeNatGatewayRegions(request *DescribeNatGatewayRegionsRequest) (response *DescribeNatGatewayRegionsResponse, err error) {
	response = NewDescribeNatGatewayRegionsResponse()
	err = c.ApiCall(request, response)
	return
}


func NewModifyEipBandwidthRequest() (request *ModifyEipBandwidthRequest) {
	request = &ModifyEipBandwidthRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "ModifyEipBandwidth")

	return
}

func NewModifyEipBandwidthResponse() (response *ModifyEipBandwidthResponse) {
	response = &ModifyEipBandwidthResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// ModifyEipBandwidth 调整弹性公网IP的带宽限速。
func (c *Client) ModifyEipBandwidth(request *ModifyEipBandwidthRequest) (response *ModifyEipBandwidthResponse, err error) {
	response = NewModifyEipBandwidthResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeEipMonitorDataRequest() (request *DescribeEipMonitorDataRequest) {
	request = &DescribeEipMonitorDataRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeEipMonitorData")

	return
}

func NewDescribeEipMonitorDataResponse() (response *DescribeEipMonitorDataResponse) {
	response = &DescribeEipMonitorDataResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeEipMonitorData 查询弹性公网IP监控指标。
func (c *Client) DescribeEipMonitorData(request *DescribeEipMonitorDataRequest) (response *DescribeEipMonitorDataResponse, err error) {
	response = NewDescribeEipMonitorDataResponse()
	err = c.ApiCall(request, response)
	return
}

func NewDescribeInstanceMonitorDataRequest() (request *DescribeInstanceMonitorDataRequest) {
	request = &DescribeInstanceMonitorDataRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeInstanceMonitorData")

	return
}

func NewDescribeInstanceMonitorDataResponse() (response *DescribeInstanceMonitorDataResponse) {
	response = &DescribeInstanceMonitorDataResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeInstanceMonitorData 查询实例监控指标。
func (c *Client) DescribeInstanceMonitorData(request *DescribeInstanceMonitorDataRequest) (response *DescribeInstanceMonitorDataResponse, err error) {
	response = NewDescribeInstanceMonitorDataResponse()
	err = c.ApiCall(request, response)
	return
}



func NewDescribeAvailableBorderGatewayRequest() (request *DescribeAvailableBorderGatewayRequest) {
	request = &DescribeAvailableBorderGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().InitWithApiInfo(SERVICE, APIVersion, "DescribeAvailableBorderGateway")

	return
}

func NewDescribeAvailableBorderGatewayResponse() (response *DescribeAvailableBorderGatewayResponse) {
	response = &DescribeAvailableBorderGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

// DescribeAvailableBorderGateway 获取可绑定NAT的边界网关。
func (c *Client) DescribeAvailableBorderGateway(request *DescribeAvailableBorderGatewayRequest) (response *DescribeAvailableBorderGatewayResponse, err error) {
	response = NewDescribeAvailableBorderGatewayResponse()
	err = c.ApiCall(request, response)
	return
}