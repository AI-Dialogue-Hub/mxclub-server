package req

import (
	"mxclub/domain/lottery/entity/enum"
	"time"
)

type LotteryPrizeReq struct {
	Id                    uint               `json:"id"`
	ProductAttributeID    uint64             `json:"productAttributeId"`
	ActivityId            uint               `json:"activityId"`
	PrizeLevel            string             `json:"prizeLevel"`
	PrizeName             string             `json:"prizeName"`
	PrizeType             enum.PrizeTypeEnum `json:"prizeType" validate:"required"`
	PrizeValue            float64            `json:"prizeValue"`
	TotalQuantity         int                `json:"totalQuantity"`
	RemainingQuantity     int                `json:"remainingQuantity"`
	DailyLimit            *int               `json:"dailyLimit"`
	UserDailyLimit        int                `json:"userDailyLimit"`
	UserTotalLimit        int                `json:"userTotalLimit"`
	PrizeImage            string             `json:"prizeImage"`
	WinMessage            string             `json:"winMessage"`
	DisplayProbability    float64            `json:"displayProbability"`
	ActualProbability     float64            `json:"actualProbability"`
	ProbabilityAdjustment float64            `json:"probabilityAdjustment"`
	SortOrder             int                `json:"sortOrder"`
	IsActive              bool               `json:"isActive"`
	StartTime             *time.Time         `json:"startTime"`
	EndTime               *time.Time         `json:"endTime"`
}
