package api

import (
	"encoding/json"
)

const CodeSuccess = 200

type Response struct {
	RequestId string `json:"request_id,omitempty"` //请求ID
	Code      int    `json:"code"`                 //错误码，200 成功，其他失败
	Message   string `json:"message,omitempty"`    //错误信息
	Data      any    `json:"data,omitempty"`
}

func (r *Response) Success() bool {
	return r.Code == CodeSuccess
}

func (r *Response) Error() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func Success(reqId string, datas ...any) *Response {
	var data any = nil
	if len(datas) > 0 {
		data = datas[0]
	}

	return &Response{
		RequestId: reqId,
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
	}
}

func Error(reqId string, code int, message string) *Response {
	return &Response{
		RequestId: reqId,
		Code:      code,
		Message:   message,
	}
}

func FromError(reqId string, err error) *Response {
	if apiErr, ok := err.(*Response); ok {
		return &Response{
			RequestId: reqId,
			Code:      apiErr.Code,
			Message:   apiErr.Message,
		}
	} else {
		return &Response{
			RequestId: reqId,
			Code:      ErrorCodeInternal,
			Message:   err.Error(),
		}
	}
}
