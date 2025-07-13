package req

import (
	"mxclub/domain/lottery/entity/enum"
	"time"
)

type LotteryPrizeReq struct {
	Id                    uint                       `json:"id"`
	ProductAttributeID    uint64                     `json:"productAttributeId"`
	ActivityId            uint                       `json:"activityId"`
	PrizeLevel            enum.LotteryPrizeLevelEnum `json:"prizeLevel"`
	PrizeName             string                     `json:"prizeName"`
	PrizeType             enum.PrizeTypeEnum         `json:"prizeType" validate:"required"`
	PrizeValue            float64                    `json:"prizeValue"`
	TotalQuantity         int                        `json:"totalQuantity"`
	RemainingQuantity     int                        `json:"remainingQuantity"`
	DailyLimit            *int                       `json:"dailyLimit"`
	UserDailyLimit        int                        `json:"userDailyLimit"`
	UserTotalLimit        int                        `json:"userTotalLimit"`
	PrizeImage            string                     `json:"prizeImage"`
	WinMessage            string                     `json:"winMessage"`
	DisplayProbability    float64                    `json:"displayProbability"`
	ActualProbability     float64                    `json:"actualProbability"`
	ProbabilityAdjustment float64                    `json:"probabilityAdjustment"`
	SortOrder             int                        `json:"sortOrder"`
	IsActive              bool                       `json:"isActive"`
	StartTime             *time.Time                 `json:"startTime"`
	EndTime               *time.Time                 `json:"endTime"`
}

type LotteryActivityReq struct {
	ID                  uint                    `json:"id"`
	ActivityPrice       float64                 `json:"activity_price" validate:"required" reg_err_info:"价格字段不能为空"`
	ActivityTitle       string                  `json:"activity_title" validate:"required"`
	ActivitySubtitle    string                  `json:"activity_subtitle" validate:"required"`
	ActivityDesc        string                  `json:"activity_desc" validate:"required"`
	EntryURL            string                  `json:"entry_url"`
	EntryImage          string                  `json:"entry_image" validate:"required"`
	BannerImage         string                  `json:"banner_image"`
	BackgroundImage     string                  `json:"background_image"`
	ActivityRules       string                  `json:"activity_rules"`
	PrizePoolID         *uint                   `json:"prize_pool_id"`
	StartTime           time.Time               `json:"start_time"`
	EndTime             time.Time               `json:"end_time"`
	ParticipateTimes    int                     `json:"participate_times"`
	ShareAddTimes       int                     `json:"share_add_times"`
	TotalPrizeCount     *int                    `json:"total_prize_count"`
	RemainingPrizeCount *int                    `json:"remaining_prize_count"`
	ActivityStatus      enum.ActivityStatusEnum `json:"activity_status"`
	DisplayOrder        int                     `json:"display_order"`
	IsFeatured          bool                    `json:"is_featured"`
	IsHot               bool                    `json:"is_hot"`
}

type LotteryActivityStatusReq struct {
	LotteryActivityId     uint                    `json:"lottery_activity_id"`
	LotteryActivityStatus enum.ActivityStatusEnum `json:"lottery_activity_status"`
}
