package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserService)
}

type UserService struct {
	userRepo repo.IUserRepo
}

func NewUserService(repo repo.IUserRepo) *UserService {
	return &UserService{userRepo: repo}
}

// =============================================================

func (svc UserService) GetUserById(ctx jet.Ctx, id int) (*vo.UserLoginVO, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopyByCtx[vo.UserLoginVO](ctx, userPO), err
}

func (svc UserService) CheckUser(ctx jet.Ctx, username string, password string) (*po.User, error) {
	userPO, err := svc.userRepo.QueryUserByAccount(username, password)
	if err != nil || userPO == nil || userPO.ID == 0 {
		ctx.Logger().Infof("user %s not exist", username)
		return nil, errors.New("账号或密码错误")
	}
	return userPO, nil
}

func (svc UserService) List(ctx jet.Ctx, params *req.UserListReq) (*api.PageResult, error) {
	if params.UserType == "all" {
		params.UserType = ""
	}
	list, count, err := svc.userRepo.ListAroundCacheByUserType(ctx, params.PageParams, enum.RoleType(params.UserType))
	if err != nil {
		ctx.Logger().Errorf("[UserService List] error:%v", err.Error())
		return nil, errors.New("查询失败")
	}
	userVOS := utils.CopySlice[*po.User, *vo.UserVO](list)
	utils.ForEach(userVOS, func(t *vo.UserVO) { t.DisPlayName = t.Role.DisPlayName() })
	return api.WrapPageResult(params.PageParams, userVOS, count), nil
}

func (svc UserService) Update(ctx jet.Ctx, req *req.UserReq) error {
	updateMap := utils.ObjToMap(req)
	updateMap["member_number"] = utils.SafeParseNumber[int](updateMap["member_number"])
	return svc.userRepo.UpdateUser(ctx, updateMap)
}

func (svc UserService) RemoveAssistant(ctx jet.Ctx, userId uint) error {
	return svc.userRepo.RemoveDasher(ctx, userId)
}
