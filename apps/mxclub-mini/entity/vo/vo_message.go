package vo

import (
	"mxclub/domain/message/entity/enum"
	"time"
)

type MessageVO struct {
	ID                 int              `json:"id"`
	MessageType        enum.MessageType `json:"message_type"`
	MessageDisPlayType string           `json:"message_display_type"`
	Title              string           `json:"title"`
	Content            string           `json:"content"`
	MessageStatus      int              `json:"message_status"`
	CreatedAt          time.Time        `json:"created_at"`
}
