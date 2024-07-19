package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewApplicationService)
}

type ApplicationService struct {
	applicationRepo repo.IAssistantApplicationRepo
	userService     *UserService
}

func NewApplicationService(applicationRepo repo.IAssistantApplicationRepo, userService *UserService) *ApplicationService {
	return &ApplicationService{applicationRepo: applicationRepo, userService: userService}
}

func (svc ApplicationService) List(ctx jet.Ctx, params *req.ApplicationListReq) (*api.PageResult, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(int32(params.Page), int32(params.PageSize))
	if params.Status != "" {
		query.SetFilter("status = ?", params.Status)
	}
	pos, count, err := svc.applicationRepo.ListByWrapper(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("ERROR:%v", err.Error())
		return nil, errors.New("查询失败")
	}
	vos := utils.CopySlice[*po.AssistantApplication, *vo.AssistantApplicationVO](pos)
	return api.WrapPageResult(params.PageParams, vos, count), nil
}

func (svc ApplicationService) UpdateStatus(ctx jet.Ctx, req *req.ApplicationReq) (err error) {
	// 1. 如果通过 修改用户表 赋予打手权限]
	if req.Status == "PASS" {
		err = svc.userService.userRepo.UpdateUser(ctx, buildUpdateWrapper(req))
		if err != nil {
			ctx.Logger().Errorf("ERROR:%v", err.Error())
			return errors.New("更新失败")
		}
	}
	// 2. 修改申请记录
	err = svc.applicationRepo.UpdateStatus(ctx, req.ID, req.Status)
	if err != nil {
		ctx.Logger().Errorf("ERROR:%v", err.Error())
		return errors.New("更新失败")
	}
	return err
}

func buildUpdateWrapper(req *req.ApplicationReq) map[string]any {
	updateMap := map[string]any{}
	updateMap["id"] = req.UserID
	updateMap["role"] = enum.RoleAssistant
	updateMap["phone"] = req.Phone
	updateMap["member_number"] = req.MemberNumber
	updateMap["name"] = req.Name
	return updateMap
}
