package req

import "mxclub/pkg/api"

type TransferListReq struct {
	api.PageParams
	Status int `form:"status"`
}
