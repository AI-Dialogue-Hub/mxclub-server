package po

import (
	"gorm.io/gorm"
	"mxclub/domain/lottery/entity/enum"
	"time"
)

// LotteryPurchaseRecord 用户购买记录
type LotteryPurchaseRecord struct {
	ID               uint                     `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID           uint                     `gorm:"not null;index;comment:用户ID"`
	ActivityID       uint                     `gorm:"not null;index;comment:活动ID"`
	TransactionID    string                   `gorm:"size:64;not null;index;comment:交易流水号"`
	PurchaseAmount   float64                  `gorm:"type:decimal(10,2);not null;comment:购买金额"`
	PurchaseTime     time.Time                `gorm:"not null;default:CURRENT_TIMESTAMP;comment:购买时间"`
	PurchaseStatus   enum.PurchaseStatusEnum  `gorm:"not null;default:1;comment:购买状态(1:成功 2:失败 3:退款)"`
	PaymentMethod    enum.PurchasePaymentEnum `gorm:"comment:支付方式(1:微信 2:支付宝 3:银行卡)"`
	PaymentTime      *time.Time               `gorm:"comment:支付完成时间"`
	LotteryQualified bool                     `gorm:"default:false;comment:是否获得抽奖资格"`
	LotteryUsed      bool                     `gorm:"default:false;comment:是否已使用抽奖资格"`
	Phone            string                   `gorm:"size:40;comment:手机号"`
	RoleId           string                   `gorm:"size:255;comment:角色ID"`
	IPAddress        *string                  `gorm:"size:45;comment:用户IP地址"`
	DeviceInfo       *string                  `gorm:"size:255;comment:设备信息"`
	CreatedAt        time.Time                `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt        time.Time                `gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt        gorm.DeletedAt           `gorm:"index;comment:删除时间"`
}

func (LotteryPurchaseRecord) TableName() string {
	return "lottery_purchase_records"
}
