package po

import (
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
)

// OrderTransfer 代表转单表的数据模型
type OrderTransfer struct {
	gorm.Model
	OrderId      uint64            `gorm:"column:order_id;not_null" json:"order_id" comment:"订单Id"`
	ExecutorFrom int               `gorm:"column:executor_from;not_null" json:"executor_from" comment:"转单人Id，一般是车头"`
	ExecutorTo   int               `gorm:"column:executor_to;default:-1" json:"executor_to" comment:"转单to Id，必须在线 && 不在订单中"`
	Status       enum.TransferEnum `gorm:"column:status" json:"status" comment:"转单状态"`
}

func (*OrderTransfer) TableName() string {
	return "orders_transfer"
}
