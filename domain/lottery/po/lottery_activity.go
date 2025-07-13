package po

import (
	"gorm.io/gorm"
	"mxclub/domain/lottery/entity/enum"
	"time"
)

// LotteryActivity 抽奖活动表
type LotteryActivity struct {
	ID                  uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	FallbackPrizeId     uint                    `json:"fallback_prize_id"`
	ActivityPrice       float64                 `gorm:"type:decimal(5,2);not null" json:"activity_price"`
	ActivityTitle       string                  `gorm:"size:100;not null" json:"activity_title"`
	ActivitySubtitle    string                  `gorm:"size:255" json:"activity_subtitle"`
	ActivityDesc        string                  `gorm:"type:text" json:"activity_desc"`
	EntryURL            string                  `gorm:"size:255;not null" json:"entry_url"`
	EntryImage          string                  `gorm:"size:255" json:"entry_image"`
	BannerImage         string                  `gorm:"size:255" json:"banner_image"`
	BackgroundImage     string                  `gorm:"size:255" json:"background_image"`
	ActivityRules       string                  `gorm:"type:text;not null" json:"activity_rules"`
	PrizePoolID         *uint                   `json:"prize_pool_id"`
	StartTime           time.Time               `gorm:"not null" json:"start_time"`
	EndTime             time.Time               `gorm:"not null" json:"end_time"`
	ParticipateTimes    int                     `gorm:"default:1" json:"participate_times"`
	ShareAddTimes       int                     `gorm:"default:0" json:"share_add_times"`
	TotalPrizeCount     *int                    `json:"total_prize_count"`
	RemainingPrizeCount *int                    `json:"remaining_prize_count"`
	ActivityStatus      enum.ActivityStatusEnum `gorm:"type:enum('pending','ongoing','paused','ended');default:'pending'" json:"activity_status"`
	DisplayOrder        int                     `gorm:"default:0" json:"display_order"`
	IsFeatured          bool                    `gorm:"default:false" json:"is_featured"` // 是否推荐活动
	IsHot               bool                    `gorm:"default:false" json:"is_hot"`
	CreatorID           *uint                   `json:"creator_id"`
	CreatedAt           time.Time               `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time               `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt           gorm.DeletedAt          `gorm:"index" json:"deleted_at"`
}

func (LotteryActivity) TableName() string {
	return "lottery_activities"
}
