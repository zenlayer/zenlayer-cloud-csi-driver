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

const (
	INTERNAL_SERVER_ERROR                     string = "INTERNAL_SERVER_ERROR"                     //服务端内部错误。
	SERVICE_TEMPORARY_UNAVAILABLE             string = "SERVICE_TEMPORARY_UNAVAILABLE"             //服务端暂时不可用，请稍后操作。
	REQUEST_TIMED_OUT                         string = "REQUEST_TIMED_OUT"                         //请求超时。
	AUTHFAILURE_INVALID_AUTHORIZATION         string = "AUTHFAILURE_INVALID_AUTHORIZATION"         //请求头部的 Authorization 不符合Zenalyer Cloud标准。
	AUTHFAILURE_SIGNATURE_FAILURE             string = "AUTHFAILURE_SIGNATURE_FAILURE"             //签名计算错误。请对照调用方式中的签名方法文档检查签名计算过程。
	AUTHFAILURE_UNAUTHORIZED_OPERATION        string = "AUTHFAILURE_UNAUTHORIZED_OPERATION"        //请求未授权，请确认接口存在或者有授予了访问权限。
	AUTHFAILURE_SIGNATURE_EXPIRED             string = "AUTHFAILURE_SIGNATURE_EXPIRED"             //签名过期。Timestamp和服务器时间不得相差超过五分钟。
	MISSING_PARAMETER                         string = "MISSING_PARAMETER"                         //缺少必须的参数。
	INVALID_PARAMETER_TYPE                    string = "INVALID_PARAMETER_TYPE"                    //参数类型不合法。
	INVALID_PARAMETER_EXCEED_MAXIMUM          string = "INVALID_PARAMETER_EXCEED_MAXIMUM"          //参数值大小超过了限制。
	INVALID_PARAMETER_LESS_MINIMUM            string = "INVALID_PARAMETER_LESS_MINIMUM"            //参数值小于了最低限制。
	INVALID_PARAMETER_EXCEED_LENGTH           string = "INVALID_PARAMETER_EXCEED_LENGTH"           //参数值数量超过了允许的最大长度。快照名或者云盘名超出zec长度限制。
	INVALID_PARAMETER_LESS_LENGTH             string = "INVALID_PARAMETER_LESS_LENGTH"             //参数值数量小于了最低要求的长度。
	INVALID_PARAMETER_PATTERN_ERROR           string = "INVALID_PARAMETER_PATTERN_ERROR"           //参数值所要求的格式不匹配。
	INVALID_PARAMETER_DATE_FORMAT_ERROR       string = "INVALID_PARAMETER_DATE_FORMAT_ERROR"       //参数值的日期格式不正确。
	INVALID_PARAMETER_VALUE                   string = "INVALID_PARAMETER_VALUE"                   //参数值没在规定的值里面。
	OPERATION_FAILED_FOR_RECYCLE_RESOURCE     string = "OPERATION_FAILED_FOR_RECYCLE_RESOURCE"     //对于回收中的资源进行操作是不被允许的。
	OPERATION_FAILED_RESOURCE_GROUP_NOT_FOUND string = "OPERATION_FAILED_RESOURCE_GROUP_NOT_FOUND" //对于指定操作的资源组不存在。
	OPERATION_FAILED_RESOURCE_NOT_FOUND       string = "OPERATION_FAILED_RESOURCE_NOT_FOUND"       //指定的资源不存在导致操作失败。
	INVALID_DISK_SNAPSHOT_NOT_FOUND           string = "INVALID_DISK_SNAPSHOT_NOT_FOUND"           //快照不存在
	INVALID_DISK_SNAPSHOT_SIZE_MISMATCH       string = "INVALID_DISK_SNAPSHOT_SIZE_MISMATCH"       //用快照创建的云盘大小小于快照大小
)
