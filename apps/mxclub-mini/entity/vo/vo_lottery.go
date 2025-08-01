package vo

import (
	"mxclub/domain/lottery/entity/enum"
	"time"
)

type LotteryPrizeVO struct {
	ID                    uint               `json:"id"`
	ProductAttributeID    uint64             `json:"productAttributeId"`
	PrizeInfo             string             `json:"prizeInfo"`
	PrizeLevel            string             `json:"prizeLevel"`
	PrizeName             string             `json:"prizeName"`
	PrizeType             enum.PrizeTypeEnum `json:"prizeType"`
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
	Phone                 string             `json:"phone"`
	RoleId                string             `json:"role_id"`
	SortOrder             int                `json:"sortOrder"`
	IsActive              bool               `json:"isActive"`
	StartTime             *time.Time         `json:"startTime"`
	EndTime               *time.Time         `json:"endTime"`
}

type LotteryActivityPrizeVO struct {
	LotteryActivity *LotteryActivityVO `json:"lotteryActivity"`
	LotteryPrizes   []*LotteryPrizeVO  `json:"lotteryPrizes"`
}

type LotteryActivityVO struct {
	ID                  uint                    `json:"id"`
	ActivityPrice       float64                 `json:"activity_price"`
	ActivityTitle       string                  `json:"activity_title"`
	ActivitySubtitle    string                  `json:"activity_subtitle"`
	ActivityDesc        string                  `json:"activity_desc"`
	EntryURL            string                  `json:"entry_url"`
	EntryImage          string                  `json:"entry_image"`
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

// LotteryVO 抽奖结果
type LotteryVO struct {
	PrizeIndex int    `json:"prize_index"` // 奖品在奖池中的位置，默认是优先级
	WinMessage string `json:"win_message"`
	PrizeImage string `json:"prize_image"`
}

type LotteryRecordsVO struct {
	Id                    uint   `json:"id"`
	ActivityId            uint   `json:"activity_id"`
	PrizeId               uint   `json:"prize_id"`
	UserId                uint   `json:"user_id"`
	OrderId               string `json:"order_id"`
	ActivityPrizeSnapshot string `json:"activity_prize_snapshot"` // 活动信息&奖品信息快照
	// ================================================
	ActivityTitle string    `json:"activity_title"`
	ActivityPrice float64   `json:"activity_price"`
	PrizeName     string    `json:"prize_name"`
	AvatarUrl     string    `json:"avatar_url"`
	CreatedAt     time.Time `json:"created_at"`
}
