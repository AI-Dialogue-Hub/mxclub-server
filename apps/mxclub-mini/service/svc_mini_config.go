package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/common/repo"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewMiniConfigService)
}

type MiniConfigService struct {
	configRepo repo.IMiniConfigRepo
}

func NewMiniConfigService(repo repo.IMiniConfigRepo) *MiniConfigService {
	return &MiniConfigService{configRepo: repo}
}

func (svc MiniConfigService) GetConfigByName(ctx jet.Ctx) (*vo.MiniConfigVO, error) {
	configPO, err := svc.configRepo.FindSwiperConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &vo.MiniConfigVO{
		ID:         configPO.ID,
		ConfigName: configPO.ConfigName,
		Content:    utils.MustByteToMap(configPO.Content),
	}, err
}
