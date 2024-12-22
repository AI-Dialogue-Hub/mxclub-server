package event

import "github.com/fengyuan-liang/jet-web-fasthttp/jet"

type EventBO struct {
	RegisterName  string                  `json:"register_name"`
	EventName     string                  `json:"event_name"`
	EventCode     int                     `json:"event_code"`
	EventCallBack func(ctx jet.Ctx) error `json:"event_call_back"`
}

const (
	EventRemoveDasher = iota
)
