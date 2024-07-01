package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/common/entity/enum"
	"mxclub/domain/common/repo"
	"mxclub/pkg/api"
)

func init() {
	jet.Provide(NewMiniConfigService)
}

type MiniConfigService struct {
	miniConfigRepo repo.IMiniConfigRepo
}

func NewMiniConfigService(repo repo.IMiniConfigRepo) *MiniConfigService {
	return &MiniConfigService{miniConfigRepo: repo}
}

func (svc MiniConfigService) List(ctx jet.Ctx, params *api.PageParams) ([]*vo.MiniConfigVO, int64, error) {
	list, count, err := svc.miniConfigRepo.ListAroundCache(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]*vo.MiniConfigVO, 0, len(list))
	for _, val := range list {
		vos = append(vos, &vo.MiniConfigVO{
			ID:         val.ID,
			ConfigName: val.ConfigName,
			Content:    val.Content,
		})
	}
	return vos, count, nil
}

func (svc MiniConfigService) Get(id int64) (*vo.MiniConfigVO, error) {
	val, err := svc.miniConfigRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &vo.MiniConfigVO{
		ID:         val.ID,
		ConfigName: val.ConfigName,
		Content:    val.Content,
	}, nil
}

func (svc MiniConfigService) Add(ctx jet.Ctx, configName string, content []map[string]any) error {
	exists, err := svc.miniConfigRepo.ExistConfig(ctx, configName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("配置文件已经存在")
	}
	if enum.MiniConfigEnum(configName).IsNotValid() {
		return errors.New(fmt.Sprintf("配置文件类型[%s]不存在", configName))
	}
	return svc.miniConfigRepo.AddConfig(ctx, configName, content)
}
