package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/common/po"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/constant"
)

func init() {
	jet.Provide(NewIMiniConfigRepo)
}

type IMiniConfigRepo interface {
	xmysql.IBaseRepo[po.MiniConfig]
	FindConfigByName(ctx jet.Ctx, configName string) (*po.MiniConfig, error)
	FindSwiperConfig(ctx jet.Ctx) (*po.MiniConfig, error)
}

func NewIMiniConfigRepo(db *gorm.DB) IMiniConfigRepo {
	repo := new(MiniConfigRepo)
	repo.Db = db.Model(new(po.MiniConfig))
	repo.Ctx = context.Background()
	return repo
}

type MiniConfigRepo struct {
	xmysql.BaseRepo[po.MiniConfig]
}

func (repo MiniConfigRepo) FindConfigByName(ctx jet.Ctx, configName string) (*po.MiniConfig, error) {
	one, err := repo.FindOne("config_name = ?", configName)
	if err != nil {
		ctx.Logger().Errorf("FindConfigByName error: %v", err.Error())
	}
	return one, err
}

const mx_mini_config_swiper = "mini_config_swiper"

func (repo MiniConfigRepo) FindSwiperConfig(ctx jet.Ctx) (*po.MiniConfig, error) {
	got, err := xredis.GetByString[po.MiniConfig](ctx, mx_mini_config_swiper)
	if err == nil {
		return got, nil
	}
	config, err := repo.FindConfigByName(ctx, "swiper")
	if err != nil {
		return nil, err
	}
	_ = xredis.SetJSONStr(mx_mini_config_swiper, config, constant.Duration_7_Day)
	return config, nil
}
