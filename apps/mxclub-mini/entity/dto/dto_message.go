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
	OrdersId      uint // db id
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
		OrdersId:      orderId, // db id
		Ext:           ext,
	}
}

// NewDispatchMessageWithFinalPrice 这里的orderId是表的主键，不是流水号
func NewDispatchMessageWithFinalPrice(
	messageTo uint, orderId uint, region string, roleId string, ext string, finalPrice float64) *MessageDTO {
	return &MessageDTO{
		MessageType: enum.DISPATCH_MESSAGE,
		Title:       "新派单",
		Content: fmt.Sprintf("新订单，区域：%v，角色：%v，订单总价格：¥ %v",
			region, roleId, utils.RoundToTwoDecimalPlaces(finalPrice)),
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
		OrdersId:      orderId, // db id
		Ext:           ext,
	}
}

// NewInviteMessage 邀请 这里的orderId是表的主键，不是流水号
func NewInviteMessage(messageTo uint, orderId uint, region string, roleId string, ext string) *MessageDTO {
	return &MessageDTO{
		MessageType:   enum.DISPATCH_MESSAGE,
		Title:         "新派单",
		Content:       fmt.Sprintf("新订单，区域：%v，角色：%v", region, roleId),
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
		OrdersId:      orderId,
		Ext:           ext,
	}
}
