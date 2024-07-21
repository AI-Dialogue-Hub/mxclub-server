package req

import "mxclub/domain/message/entity/enum"

type MessageReq struct {
	MessageType   enum.MessageType `json:"message_type"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	MessageFrom   int              `json:"message_from"`
	MessageTo     int              `json:"message_to"`
	MessageStatus int              `json:"message_status"`
	Ext           string           `json:"ext"`
}

type MessageReadReq struct {
	Id          uint   `json:"id"`
	MessageType string `json:"message_type"`
	IsRefuse    bool   `json:"isRefuse"`
}

type MessageHandleReq struct {
	MessageId         uint             `json:"message_id"`
	OrdersId          uint             `json:"order_id"` // 这里的id是db的主键
	MessageTypeNumber uint             `json:"message_type_number"`
	MessageType       enum.MessageType `json:"message_type"`
	Ext               string           `json:"ext"`
}
