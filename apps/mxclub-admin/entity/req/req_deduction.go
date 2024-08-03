package req

import "mxclub/pkg/api"

type DeductionListReq struct {
	*api.PageParams
	Ge string `json:"search_GE_createTime"` // start time
	Le string `json:"search_LE_createTime"` // end time
}
