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
	messageService  *MessageService
}

func NewApplicationService(applicationRepo repo.IAssistantApplicationRepo,
	userService *UserService,
	messageService *MessageService) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		userService:     userService,
		messageService:  messageService,
	}
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
		// 1. 给用户发消息
		_ = svc.messageService.PushSystemMessage(ctx, req.UserID, "您的打手申请已审核通过，请重新进入小程序")
	} else if req.Status == "REJECT" {
		_ = svc.messageService.PushSystemMessage(ctx, req.UserID, "您的打手申请被拒绝，详情请联系客服")
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
