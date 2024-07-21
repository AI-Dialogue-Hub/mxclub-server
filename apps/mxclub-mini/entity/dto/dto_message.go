package dto

import (
	"fmt"
	"mxclub/domain/message/entity/enum"
)

type MessageDTO struct {
	MessageType   enum.MessageType
	Title         string
	Content       string
	MessageFrom   int
	MessageTo     uint
	MessageStatus enum.MessageStatus
	OrderId       uint
	Ext           string
}

// NewDispatchMessage 这里的orderId是表的主键，不是流水号
func NewDispatchMessage(messageTo uint, orderId uint, region string, roleId string, ext string) *MessageDTO {
	return &MessageDTO{
		MessageType:   enum.DISPATCH_MESSAGE,
		Title:         "新派单",
		Content:       fmt.Sprintf("新订单，区域：%v，角色：%v", region, roleId),
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
		OrderId:       orderId,
		Ext:           ext,
	}
}
