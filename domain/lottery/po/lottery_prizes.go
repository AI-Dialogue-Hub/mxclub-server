package po

import (
	"gorm.io/gorm"
	"time"
)

// LotteryPrize 抽奖奖品表
type LotteryPrize struct {
	ID                    uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductAttributeID    uint           `gorm:"not null" json:"product_attribute_id"`
	PrizeLevel            string         `gorm:"size:50;not null" json:"prize_level"`
	PrizeName             string         `gorm:"size:100;not null" json:"prize_name"`
	PrizeType             string         `gorm:"type:enum('physical','virtual','coupon','points','empty');not null" json:"prize_type"`
	PrizeValue            float64        `gorm:"type:decimal(10,2)" json:"prize_value"`
	TotalQuantity         int            `gorm:"not null" json:"total_quantity"`
	RemainingQuantity     int            `gorm:"not null" json:"remaining_quantity"`
	DailyLimit            *int           `gorm:"default:null" json:"daily_limit"`
	UserDailyLimit        int            `gorm:"default:1" json:"user_daily_limit"`
	UserTotalLimit        int            `gorm:"default:1" json:"user_total_limit"`
	PrizeImage            string         `gorm:"size:255" json:"prize_image"`
	WinMessage            string         `gorm:"size:255;not null" json:"win_message"`
	DisplayProbability    float64        `gorm:"type:decimal(5,2);not null" json:"display_probability"`
	ActualProbability     float64        `gorm:"type:decimal(5,2);not null" json:"actual_probability"`
	ProbabilityAdjustment float64        `gorm:"type:decimal(5,2);default:1.0" json:"probability_adjustment"`
	SortOrder             int            `gorm:"default:0" json:"sort_order"`
	IsActive              bool           `gorm:"default:true" json:"is_active"`
	StartTime             *time.Time     `json:"start_time"`
	EndTime               *time.Time     `json:"end_time"`
	CreatedAt             time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (LotteryPrize) TableName() string {
	return "lottery_prize"
}
