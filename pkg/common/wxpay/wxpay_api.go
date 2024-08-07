package wxpay

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"io"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
	"net/http"
	"net/http/httptest"
	"time"
)

func Prepay(ctx jet.Ctx, prePayRequestDTO *prepayRequestDTO) (prepayDTO *PrePayDTO, err error) {
	outTradeNo := prePayRequestDTO.OutTradeNo
	request := jsapi.PrepayRequest{
		Appid:         core.String(wxPayConfig.AppId),
		Mchid:         core.String(wxPayConfig.MchID),
		Description:   core.String("明星电竞-代打订单"),
		OutTradeNo:    core.String(outTradeNo),
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
		AppId:      wxPayConfig.AppId,
		TimeStamp:  timeStampStr,
		NonceStr:   nonceStr,
		Package:    packageStr,
		SignType:   "RSA",
		PaySign:    signature,
		OutTradeNo: outTradeNo,
	}
	return
}

// Refunds 退款 仅支持全部退款
func Refunds(ctx jet.Ctx, transaction *payments.Transaction, outRefundNo, reason string) error {
	log := ctx.Logger()
	request := refunddomestic.CreateRequest{
		TransactionId: transaction.TransactionId,
		OutTradeNo:    transaction.OutTradeNo,
		OutRefundNo:   core.String(outRefundNo),
		Reason:        core.String(reason),
		NotifyUrl:     core.String("https://mx.fengxianhub.top/wxpay/refunds/notify"),
		//FundsAccount:  refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
		Amount: &refunddomestic.AmountReq{
			Currency: core.String("CNY"),
			//From: []refunddomestic.FundsFromItem{{
			//	Account: refunddomestic.ACCOUNT_AVAILABLE.Ptr(),
			//	Amount:  transaction.Amount.PayerTotal,
			//}},
			Refund: transaction.Amount.PayerTotal,
			Total:  transaction.Amount.Total,
		},
	}
	resp, result, err := refundsApiService.Create(context.Background(), request)
	if err != nil {
		// 处理错误
		log.Printf("call Create err:%s", err)
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
		return err
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
	return nil
}

// DecryptWxpayCallBack 解密支付成功后的回调
func DecryptWxpayCallBack(ctx jet.Ctx) (*payments.Transaction, error) {
	defer utils.RecoverAndLogError(ctx)
	var (
		req *http.Request
		err error
	)
	req, err = xjet.ConvertFastHTTPRequestToStandard(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		ctx.Logger().Errorf("req is nil, %v", req)
		req = httptest.NewRequest(
			http.MethodGet, "http://127.0.0.1", io.NopCloser(bytes.NewBuffer(ctx.Request().Body())),
		)
		req.Header.Add("test", "test")
	}
	transaction := new(payments.Transaction)
	notifyReq, err := notifyHandler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Infof("notifyReq Summary:%v\n", notifyReq.Summary)
	ctx.Logger().Infof("transactionId:%v", *transaction.TransactionId)
	return transaction, nil
}
