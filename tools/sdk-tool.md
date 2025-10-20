# SDK demo

Refer to the document for more information [API doc](https://docs.console.zenlayer.com/api-reference/cn/zec/disk)		

```go
package main

import (
	"fmt"
	zec "github.com/zenlayer/zenlayercloud-sdk-go/zenlayercloud/zec20250901"
)

var AK string = ""
var SK string = ""

var client *zec.Client

//show all regions
func DescribeDiskRegions() {
	request := zec.NewDescribeDiskRegionsRequest()
	rep, _ := client.DescribeDiskRegions(request)
}

//Manually uninstall a disk from a virtual machine
func DetachDisks() {
	request := zec.NewDetachDisksRequest()
	request.DiskIds = append(request.DiskIds, "")
	var check bool = false
	request.InstanceCheckFlag = &check
	client.DetachDisks(request)
}

//show disk info
func DescribeDisks() {
	request := zec.NewDescribeDisksRequest()
	rsp, _ := client.DescribeDisks(request)
}

func main() {
	client, _ = zec.NewClientWithSecretKey(AK, SK)
	DetachDisks()
}

```