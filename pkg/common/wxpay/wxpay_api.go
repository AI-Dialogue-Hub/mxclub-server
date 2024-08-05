package wxpay

import (
	"context"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"mxclub/pkg/utils"
	"time"
)

func Prepay(ctx jet.Ctx, prePayRequestDTO *prepayRequestDTO) (prepayDTO *PrePayDTO, err error) {
	request := jsapi.PrepayRequest{
		Appid:         core.String(wxPayConfig.AppId),
		Mchid:         core.String(wxPayConfig.MchID),
		Description:   core.String("明星电竞-代打订单"),
		OutTradeNo:    core.String(prePayRequestDTO.OutTradeNo),
		TimeExpire:    core.Time(time.Now().Add(time.Minute * 15)),
		Attach:        core.String("自定义数据说明"),
		NotifyUrl:     core.String("https://mx.fengxianhub.top/wxpay/notify"),
		GoodsTag:      core.String("MX"),
		SupportFapiao: core.Bool(false),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(prePayRequestDTO.Amount),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(prePayRequestDTO.Openid),
		},
	}
	resp, result, err := jsapiApiService.Prepay(context.Background(), request)
	if err != nil || resp == nil || *resp.PrepayId == "" {
		ctx.Logger().Errorf("[Prepay]ERROR: %v\nresp:%v\nresult:%v", err.Error(), resp, utils.ObjToJsonStr(result))
		err = errors.New("请求微信接口失败")
		return
	}
	var (
		timeStampStr = fmt.Sprintf("%v", time.Now().Unix())
		nonceStr     = generateNonceStr()
		packageStr   = fmt.Sprintf("prepay_id=%v", *resp.PrepayId)
	)
	signature, err := getRSASignature([]string{
		wxPayConfig.AppId,
		timeStampStr,
		nonceStr,
		packageStr,
	})
	if err != nil {
		ctx.Logger().Errorf("[Prepay]getRSASignature: %v\nresp:%v\nresult:%v", err.Error(), resp, result)
	}
	prepayDTO = &PrePayDTO{
		AppId:     wxPayConfig.AppId,
		TimeStamp: timeStampStr,
		NonceStr:  nonceStr,
		Package:   packageStr,
		SignType:  "RSA",
		PaySign:   signature,
	}
	return
}

// refunds 退款
func refunds() {
}
