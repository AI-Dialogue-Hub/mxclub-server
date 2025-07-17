package enum

type PurchaseStatusEnum int

// 购买状态常量
const (
	_                     PurchaseStatusEnum = iota
	PurchaseStatusHold                       // 待支付
	PurchaseStatusSuccess                    // 成功
	PurchaseStatusFailed                     // 失败
	PurchaseStatusRefund                     // 退款
)

type PurchasePaymentEnum int

// 支付方式常量
const (
	PaymentMethodWeChat   PurchasePaymentEnum = 1 // 微信
	PaymentMethodAlipay   PurchasePaymentEnum = 2 // 支付宝
	PaymentMethodBankCard PurchasePaymentEnum = 3 // 银行卡
)
