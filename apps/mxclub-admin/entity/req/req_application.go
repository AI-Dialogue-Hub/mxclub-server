package req

import (
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/pkg/api"
)

type ApplicationListReq struct {
	*api.PageParams
	Status string `json:"status"`
}

type ApplicationReq struct {
	vo.AssistantApplicationVO
}
