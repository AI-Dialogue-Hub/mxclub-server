package cron

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"gopkg.in/robfig/cron.v2"
	"mxclub/apps/mxclub-mini/config"
	"mxclub/apps/mxclub-mini/service"
	"sync"
)

func init() {
	jet.Provide(NewCronService)
	jet.Invoke(func(s ICronService) {
		s.RunCron()
	})
}

type ICronService interface {
	RunCron()
}

func NewCronService(config *config.Config, orderService *service.OrderService) ICronService {
	return &CronService{
		c:            cron.New(),
		once:         new(sync.Once),
		config:       config,
		logger:       xlog.NewWith("cron_service"),
		orderService: orderService,
	}
}

type CronService struct {
	c      *cron.Cron
	once   *sync.Once
	config *config.Config
	logger *xlog.Logger
	// ================
	orderService *service.OrderService
}

// RunCron 注意在集群情况下需要指定单台机器执行定时任务，防止多次执行
func (cronService *CronService) RunCron() {
	cronService.logger.Infof("[RunCron]...")
	cronService.c.AddFunc("0 0 3 * * *", func() {
		cronService.logger.Infof("[RunCron Func]...")
		cronService.orderService.SyncDeductionInfo()
	})
	cronService.once.Do(func() {
		cronService.c.Start()
	})
}
