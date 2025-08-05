package req

import (
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/pkg/api"
)

type DeductionListReq struct {
	*api.PageParams
	Ge     string `json:"search_GE_createTime"` // start time
	Le     string `json:"search_LE_createTime"` // end time
	Status string `json:"status"`
}

type DeductionAddReq struct {
	ID              uint    `json:"id"`
	DasherId        uint    `json:"dasher_id"`
	UserID          uint    `json:"user_id"`
	UserInfo        string  `json:"user_info"`
	ConfirmPersonId uint    `json:"confirm_person_id"`
	Amount          float64 `json:"amount"`
	Reason          string  `json:"reason"`
	Status          string  `json:"status"`
}

type DeductionUpdateReq struct {
	vo.DeductionVO
}
