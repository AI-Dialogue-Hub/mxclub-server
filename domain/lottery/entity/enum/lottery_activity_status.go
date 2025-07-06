package enum

import "github.com/fengyuan-liang/GoKit/collection/maps"

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
	Virtual  PrizeTypeEnum = "virtual" // 一般是直接关联代打订单
	Coupon   PrizeTypeEnum = "coupon"  // 优惠券
	Points   PrizeTypeEnum = "points"  // 积分
	Empty    PrizeTypeEnum = "empty"
)

var PrizeTypeNames = func() maps.IMap[PrizeTypeEnum, string] {
	linkedHashMap := maps.NewLinkedHashMap[PrizeTypeEnum, string]()
	linkedHashMap.PutAll([]*maps.Pair[PrizeTypeEnum, string]{
		{Physical, "实物"},
		{Virtual, "虚拟物品"},
		{Coupon, "优惠券"},
		{Points, "积分"},
		{Empty, "空奖（谢谢参与）"},
	})
	return linkedHashMap
}()
