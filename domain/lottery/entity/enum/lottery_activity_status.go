package enum

type ActivityStatusEnum string

const (
	Pending ActivityStatusEnum = "pending"
	Ongoing ActivityStatusEnum = "ongoing"
	Paused  ActivityStatusEnum = "paused"
	Ended   ActivityStatusEnum = "ended"
)

type PrizeTypeEnum string

const (
	Physical PrizeTypeEnum = "physical"
	Virtual  PrizeTypeEnum = "virtual"
	Coupon   PrizeTypeEnum = "coupon"
	Points   PrizeTypeEnum = "points" // 积分
	Empty    PrizeTypeEnum = "empty"
)
