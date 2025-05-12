package xmysql

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogAdapter struct{}

func NewGormLogAdapter() *GormLogAdapter {
	// init log
	gormDefaultLogger = xlog.NewWith("GormLogger")
	gormDefaultLogger.SetOutputLevel(xlog.Ldebug)
	gormDefaultLogger.SetCalldPath(6)
	return new(GormLogAdapter)
}

var gormDefaultLogger *xlog.Logger

func (l *GormLogAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	getLogger(ctx).Infof(s, i...)
}

func (l *GormLogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, affected := fc()
	getLogger(ctx).Debugf("=> [%v] [rows:%v] %s", time.Since(begin), affected, sql)
}

func (l *GormLogAdapter) LogMode(level logger.LogLevel) logger.Interface {
	gormDefaultLogger.Infof("GormLogger Mode => %v", level)
	return l
}

func (l *GormLogAdapter) Info(ctx context.Context, s string, args ...interface{}) {
	getLogger(ctx).Infof(s, args...)
}

func (l *GormLogAdapter) Warn(ctx context.Context, s string, args ...interface{}) {
	getLogger(ctx).Warnf(s, args...)
}

// WthLogger 将自定义 logger 放入 context 中
func WthLogger(ctx context.Context, logger *xlog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// getLogger 从 context 中获取自定义 logger
func getLogger(ctx context.Context) *xlog.Logger {
	val := ctx.Value("logger")
	if val == nil {
		return gormDefaultLogger
	}
	customLogger, ok := val.(*xlog.Logger)
	if !ok {
		return gormDefaultLogger
	}
	return customLogger
}

func SetLoggerPrefix(requestId string) {
	gormDefaultLogger.SetPrefix(requestId)
}
