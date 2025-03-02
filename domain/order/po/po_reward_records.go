package po

import (
	"gorm.io/gorm"
	"time"
)

// RewardRecord 定义打赏记录结构体
type RewardRecord struct {
	Id           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                          // 自增主键
	PurchaserID  int       `gorm:"type:int;not null" json:"purchaser_id"`                       // 购买人的ID
	DasherID     int       `gorm:"type:int;not null" json:"dasher_id"`                          // 打手的ID
	DasherNumber string    `gorm:"type:varchar(50);not null" json:"dasher_number"`              // 打手编号
	Remarks      string    `gorm:"type:text" json:"remarks"`                                    // 备注信息
	RewardAmount float64   `gorm:"type:decimal(10,2);not null" json:"reward_amount"`            // 打赏金额
	RewardTime   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"reward_time"` // 打赏时间

	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`                             // 创建时间
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`                                                                // 删除时间，用于软删除
}

// TableName 返回表名
func (RewardRecord) TableName() string {
	return "reward_records"
}
