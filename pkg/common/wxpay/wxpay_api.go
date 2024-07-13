package wxpay

import (
	"context"
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"time"
)

func Prepay(ctx jet.Ctx, prePayRequestDTO *PrepayRequestDTO) (prepayId string, err error) {
	request := jsapi.PrepayRequest{
		Appid:         core.String(wxPayConfig.AppId),
		Mchid:         core.String(wxPayConfig.AppId),
		Description:   core.String("明星电竞-代打订单"),
		OutTradeNo:    core.String(prePayRequestDTO.OutTradeNo),
		TimeExpire:    core.Time(time.Now().Add(time.Minute * 15)),
		Attach:        core.String("自定义数据说明"),
		NotifyUrl:     core.String("https://mx.fengxianhub.top/wxpay/notify"),
		GoodsTag:      core.String("MX"),
		LimitPay:      []string{"wx_pay"},
		SupportFapiao: core.Bool(false),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(100),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(prePayRequestDTO.Openid),
		},
	}
	resp, result, err := jsapiApiService.Prepay(context.Background(), request)
	if err != nil || resp == nil || *resp.PrepayId == "" {
		ctx.Logger().Errorf("[Prepay]ERROR: %v\nresp:%v\nresult:%v", err.Error(), resp, result)
		err = errors.New("请求微信接口失败")
		return
	}
	prepayId = *resp.PrepayId
	return
}
