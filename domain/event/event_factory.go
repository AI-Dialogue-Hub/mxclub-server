package event

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/pkg/utils"
	"sync"
)

var (
	// key eventCode value event array
	eventFactory = maps.NewConcurrentLinkedHashMap[int, []*EventBO]()
	logger       = xlog.NewWith("event_logger")
)

func RegisterEvent(registerName string, eventCode int, eventHandler func(ctx jet.Ctx) error) {
	if eventFactory.ContainsKey(eventCode) {
		events := eventFactory.MustGet(eventCode)
		events = append(events, &EventBO{RegisterName: registerName, EventCode: eventCode, EventCallBack: eventHandler})
		eventFactory.Put(eventCode, events)
	} else {
		bos := make([]*EventBO, 0)
		bos = append(bos, &EventBO{RegisterName: registerName, EventCode: eventCode, EventCallBack: eventHandler})
		eventFactory.Put(eventCode, bos)
	}
	xlog.Infof("[event#RegisterEvent] eventFactory now is => %v", eventFactory.Values())
}

func PublishEvent(eventCode int, ctx jet.Ctx) {
	defer utils.RecoverByPrefix(logger, "[event#PublishEvent]")
	events, ok := eventFactory.Get(eventCode)
	if !ok {
		ctx.Logger().Errorf("[event#PublishEvent] cannot find eventCode: %v", eventCode)
		return
	}
	wg := new(sync.WaitGroup)
	for _, event := range events {
		ctx.Logger().Infof("do PublishEvent, event:%v", utils.ObjToJsonStr(event))
		wg.Add(1)
		go func(finalEvent *EventBO) {
			defer utils.RecoverByPrefix(logger, "[event#PublishEvent]")
			defer wg.Done()
			ctx.Logger().Infof(
				"[event#PublishEvent] do event callback, register_name:%v, code:%v",
				finalEvent.RegisterName, finalEvent.EventCode)
			if err := finalEvent.EventCallBack(ctx); err != nil {
				ctx.Logger().Errorf("[event#PublishEvent] ERROR: %v", err)
			}
			ctx.Logger().Infof(
				"[event#PublishEvent] handle success, register_name:%v, code:%v",
				finalEvent.RegisterName, finalEvent.EventCode)
		}(event)
	}
	wg.Wait()
}
