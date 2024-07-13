package wxpay

type PrepayRequestDTO struct {
	OutTradeNo string // 交易单号 必须唯一, ex: 1217752501201407033233368018
	Amount     int64  // 交易金额
	Openid     string // 用户的openId

}
