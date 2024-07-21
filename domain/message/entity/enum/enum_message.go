package enum

import "github.com/fengyuan-liang/GoKit/collection/maps"

type MessageType string

const (
	ALL                 MessageType = "all"
	SYSTEM_NOTIFICATION MessageType = "system_notification"
	DISPATCH_MESSAGE    MessageType = "dispatch_message"
	ACCEPTANCE_MESSAGE  MessageType = "acceptance_message"
	REMOVE_MESSAGE      MessageType = "remove_message" // 移除在进行中的队友
)

var MessageType2DisPlayNameMap = map[MessageType]string{
	ALL:                 "全部消息", // 所有人都可以查看
	SYSTEM_NOTIFICATION: "系统通知", // 系统通知
	DISPATCH_MESSAGE:    "新派单",
	ACCEPTANCE_MESSAGE:  "接单成功",
	REMOVE_MESSAGE:      "移除在进行中的订单",
}

var DisPlayName2MessageTypeMap = func() maps.IMap[string, MessageType] {
	linkedHashMap := maps.NewLinkedHashMapWithExpectedSize[string, MessageType](len(MessageType2DisPlayNameMap))
	for k, v := range MessageType2DisPlayNameMap {
		linkedHashMap.Put(v, k)
	}
	return linkedHashMap
}()

func (msg MessageType) IsValid() bool {
	_, ok := MessageType2DisPlayNameMap[msg]
	return ok
}

func (msg MessageType) ParseDisPlayName() string {
	return MessageType2DisPlayNameMap[msg]
}

type MessageStatus int

const (
	UN_READ MessageStatus = iota
	READ
)
