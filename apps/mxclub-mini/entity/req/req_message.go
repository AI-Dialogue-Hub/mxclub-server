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
