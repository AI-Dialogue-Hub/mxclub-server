package po

import (
	"gorm.io/gorm"
)

type LotteryRecords struct {
	gorm.Model
	ActivityId            uint
	PrizeId               uint
	UserId                uint
	OrderId               string
	ActivityPrizeSnapshot string // 活动信息&奖品信息快照
}

func (LotteryRecords) TableName() string {
	return "lottery_records"
}
