package xjet

import "github.com/fengyuan-liang/jet-web-fasthttp/jet"

func NewCommonJetController[T jet.IJetController]() {
	jet.Provide(func() jet.ControllerResult { return jet.NewJetController(new(T)) })
}

type JetContext struct {
	Ctx jet.Ctx
}
