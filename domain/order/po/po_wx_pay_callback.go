package po

import (
	"mxclub/pkg/common/xmysql"
	"time"
)

// WxPayCallback represents the structure of the wx_pay_callback table
type WxPayCallback struct {
	ID         int64       `json:"id" gorm:"column:id"`
	OutTradeNo string      `json:"out_trade_no" gorm:"column:out_trade_no"`
	RawData    xmysql.JSON `json:"raw_data" gorm:"column:raw_data"` // 原始回调数据

	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`           // 创建时间
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`           // 更新时间
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"` // 删除时间, 使用指针类型来表示可以为空
}

func (*WxPayCallback) TableName() string {
	return "wx_pay_callback"
}
