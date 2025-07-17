package ability

import "mxclub/pkg/common/wxpay"

// GenerateUniqueOrderNumber 生成唯一订单号，9开头表示这是抽奖活动
func GenerateUniqueOrderNumber() string {
	return "9" + wxpay.GenerateUniqueOrderNumber()
}

func IsLotteryOrder(orderNo string) bool {
	if orderNo == "" || len(orderNo) < 0 {
		return false
	}
	return orderNo[0] == '9'
}
