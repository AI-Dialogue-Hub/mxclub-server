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
	vo.DeductionVO
}

type DeductionUpdateReq struct {
	vo.DeductionVO
}
