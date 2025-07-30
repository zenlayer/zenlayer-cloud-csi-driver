package common

import "fmt"

const (
	JsonParseError         = "JSON_PARSE_FAILED"
	NetworkError           = "NETWORK_ERROR"
	CredentialMissingError = "CREDENTIAL_VALUE_MISSING"
	IoError                = "IO_ERROR"
)

type ZenlayerCloudSdkError struct {
	Code      string
	Message   string
	RequestId string
}

func (e *ZenlayerCloudSdkError) Error() string {
	if e.RequestId == "" {
		return fmt.Sprintf("[ZenlayerCloudSdkError] Code=%s, Message=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("[ZenlayerCloudSdkError] Code=%s, Message=%s, RequestId=%s", e.Code, e.Message, e.RequestId)
}

func NewZenlayerCloudSdkError(code, message, requestId string) error {
	return &ZenlayerCloudSdkError{
		Code:      code,
		Message:   message,
		RequestId: requestId,
	}
}

func (e *ZenlayerCloudSdkError) GetCode() string {
	return e.Code
}

func (e *ZenlayerCloudSdkError) GetMessage() string {
	return e.Message
}

func (e *ZenlayerCloudSdkError) GetRequestId() string {
	return e.RequestId
}
