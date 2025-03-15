package service

import (
	"github.com/fengyuan-liang/GoKit/collection/stream"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/common/entity/enum"
	"mxclub/domain/common/po"
	"mxclub/domain/common/repo"
	"mxclub/pkg/api"
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

func (svc MiniConfigService) GetConfigByName(ctx jet.Ctx, configName string) (*vo.MiniConfigVO, error) {
	configPO, err := svc.configRepo.FindConfigByName(ctx, configName)
	if err != nil {
		return nil, err
	}
	return &vo.MiniConfigVO{
		ID:         configPO.ID,
		ConfigName: configPO.ConfigName,
		Content:    configPO.Content,
	}, err
}

func (svc MiniConfigService) FetchSellingPoints(ctx jet.Ctx) (vos []map[string]string) {
	configVO, err := svc.GetConfigByName(ctx, enum.SellingPoint.String())
	if err != nil {
		return
	}
	content := configVO.Content
	if err != nil {
		return
	}
	for _, chunk := range content {
		desc := chunk["desc"]
		if descStr, ok := desc.(string); ok {
			vos = append(vos, map[string]string{"text": descStr})
		}
	}
	return
}

func (svc MiniConfigService) List(ctx jet.Ctx, params *api.PageParams) ([]*vo.MiniConfigVO, int64, error) {
	list, count, err := svc.configRepo.ListAroundCache(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	collect := stream.Of[*po.MiniConfig, *vo.MiniConfigVO](list).
		Map(func(ele *po.MiniConfig) *vo.MiniConfigVO {
			return &vo.MiniConfigVO{
				ID:          ele.ID,
				ConfigName:  ele.ConfigName,
				DisPlayName: enum.MiniConfigEnum(ele.ConfigName).DisPlayName(),
				Content:     ele.Content,
			}
		}).
		CollectToSlice()
	return collect, count, nil
}
