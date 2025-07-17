package req

type LotteryStartReq struct {
	ActivityId uint `json:"activity_id"`
}

type LotteryCanDrawReq struct {
	ActivityId uint `json:"activity_id" form:"activity_id" validate:"required" reg_err_info:"不能为空"`
}
