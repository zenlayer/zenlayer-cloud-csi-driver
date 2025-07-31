# SDK demo

Refer to the document for more information [API doc](https://docs.console.zenlayer.com/api-reference/cn/zec/disk)		

```go
package main

import (
	"fmt"
	zec "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/zec20240401"
)

var AK string = ""
var SK string = ""

var client *zec.Client

//show all regions
func DescribeDiskRegions() {
	request := zec.NewDescribeDiskRegionsRequest()
	rep, _ := client.DescribeDiskRegions(request)
	fmt.Println(rep.Response.RegionIds)
}

//Manually uninstall a disk from a virtual machine
func DetachDisks() {
	detachrequest := zec.NewDetachDisksRequest()
	detachrequest.DiskIds = append(detachrequest.DiskIds, "")
	var check bool = false
	detachrequest.InstanceCheckFlag = &check
	client.DetachDisks(detachrequest)
}

//show disk info
func DescribeDisks() {
	describrequest := zec.NewDescribeDisksRequest()
	describrequest.ZoneId = "asia-north-1a"
	describrequest.DiskIds = append(describrequest.DiskIds, "")
	rsp, _ := client.DescribeDisks(describrequest)
	fmt.Println(rsp.Response.DataSet[0])
}

func main() {
	client, _ = zec.NewClientWithSecretKey(AK, SK)
	DetachDisks()
}

```