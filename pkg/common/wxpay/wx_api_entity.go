package wxpay

type prepayRequestDTO struct {
	OutTradeNo string // 交易单号 必须唯一, ex: 1217752501201407033233368018
	Amount     int64  // 交易金额  订单总金额，单位为分
	Openid     string // 用户的openId
}

func NewPrepayRequest(amount float64, openid string) *prepayRequestDTO {
	return &prepayRequestDTO{
		OutTradeNo: generateUniqueOrderNumber(),
		Amount:     int64(amount * 100),
		Openid:     openid,
	}
}

type PrePayDTO struct {
	AppId     string `json:"app_id"`
	TimeStamp string `json:"time_stamp"` // 10位
	NonceStr  string `json:"nonce_str"`  // 32位
	Package   string `json:"package"`    // prepay_id=***
	SignType  string `json:"sign_type"`  // RSA
	PaySign   string `json:"pay_sign"`   // 签名
}
