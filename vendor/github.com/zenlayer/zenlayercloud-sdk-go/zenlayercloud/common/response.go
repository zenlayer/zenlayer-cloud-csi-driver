package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response interface {
	ParseErrorResponse(body []byte) error
}

type BaseResponse struct {
}

type ErrorResponse struct {
	RequestId string `json:"requestId"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func (r *BaseResponse) ParseErrorResponse(body []byte) (err error) {
	resp := &ErrorResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json body: %s, error: %s", body, err)
		return NewZenlayerCloudSdkError(JsonParseError, msg, "")
	}
	return NewZenlayerCloudSdkError(resp.Code, resp.Message, resp.RequestId)
}

func ParseFromHttpResponse(hr *http.Response, response Response) (err error) {
	defer hr.Body.Close()
	body, err := ioutil.ReadAll(hr.Body)
	if err != nil {
		msg := fmt.Sprintf("Fail to read response body, error %s", err)
		return NewZenlayerCloudSdkError(IoError, msg, "")
	}

	if hr.StatusCode != 200 {
		_ = fmt.Sprintf("Request fail with http status code: %s, with body: %s", hr.Status, body)

		// parse error
		return response.ParseErrorResponse(body)
	}

	//todo log.Printf("[DEBUG] Response body=%s", body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := fmt.Sprintf("Fail to parse json body: %s, error: %s", body, err)
		return NewZenlayerCloudSdkError(JsonParseError, msg, "")
	}
	return
}
