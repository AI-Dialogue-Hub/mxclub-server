package api

import (
	"encoding/json"
)

type Resp[T any] struct {
	RequestId string `json:"request_id,omitempty"` //请求ID
	Code      int    `json:"code"`                 //错误码，200 成功，其他失败
	Message   string `json:"message,omitempty"`    //错误信息
	Data      T      `json:"data,omitempty"`
}

type R[T any] struct {
	RequestId string `json:"request_id,omitempty"` //请求ID
	Code      int    `json:"code"`                 //错误码，200 成功，其他失败
	Message   string `json:"message,omitempty"`    //错误信息
	Data      T      `json:"data,omitempty"`
}

func (r *Resp[T]) Success() bool {
	return r.Code == CodeSuccess
}

func (r *Resp[T]) Error() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func Ok[T any](reqId string, datas ...T) *R[T] {
	var data T
	if len(datas) > 0 {
		data = datas[0]
	}

	return &R[T]{
		RequestId: reqId,
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
	}
}

func Failed[T any](reqId string, code int, message string) *Resp[T] {
	return &Resp[T]{
		RequestId: reqId,
		Code:      code,
		Message:   message,
	}
}

func FailedFromError[T any](reqId string, err error) *Resp[T] {
	if apiErr, ok := err.(*Resp[T]); ok {
		return &Resp[T]{
			RequestId: reqId,
			Code:      apiErr.Code,
			Message:   apiErr.Message,
		}
	} else {
		return &Resp[T]{
			RequestId: reqId,
			Code:      ErrorCodeInternal,
			Message:   err.Error(),
		}
	}
}
