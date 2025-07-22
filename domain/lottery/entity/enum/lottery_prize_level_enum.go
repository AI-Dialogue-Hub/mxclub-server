package enum

type LotteryPrizeLevelEnum string

const (
	PrizeLevelFirst  LotteryPrizeLevelEnum = "一等奖"
	PrizeLevelSecond LotteryPrizeLevelEnum = "二等奖"
	PrizeLevelThird  LotteryPrizeLevelEnum = "三等奖"
)

func (p LotteryPrizeLevelEnum) String() string {
	return string(p)
}
