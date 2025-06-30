package po

import (
	"gorm.io/gorm"
	"time"
)

// ActivityPrizeRelation 活动奖品关联表
type ActivityPrizeRelation struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint           `gorm:"not null" json:"activity_id"`
	PrizeID    uint           `gorm:"not null" json:"prize_id"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (ActivityPrizeRelation) TableName() string {
	return "activity_prize_relations"
}
