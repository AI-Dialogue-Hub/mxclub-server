package dto

import "time"

type LotteryRecordsDTO struct {
	Id                    uint
	ActivityId            uint
	PrizeId               uint
	UserId                uint
	OrderId               string
	ActivityPrizeSnapshot string // 活动信息&奖品信息快照
	// ================================================
	ActivityTitle string
	ActivityPrice float64
	PrizeName     string
	CreatedAt     time.Time
}
