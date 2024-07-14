package dto

import (
	"fmt"
	"mxclub/domain/message/entity/enum"
	"mxclub/pkg/utils"
)

type MessageDTO struct {
	MessageType   enum.MessageType
	Title         string
	Content       string
	MessageFrom   int
	MessageTo     uint
	MessageStatus enum.MessageStatus
	Ext           string
}

// NewDispatchMessage 这里的orderId是表的主键，不是流水号
func NewDispatchMessage(messageTo uint, orderId uint, region string, roleId string) *MessageDTO {
	return &MessageDTO{
		MessageType:   enum.DISPATCH_MESSAGE,
		Title:         "新派单",
		Content:       fmt.Sprintf("新订单，区域：%v，角色：%v", region, roleId),
		MessageFrom:   0,
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
		Ext:           utils.ParseString(orderId),
	}
}
