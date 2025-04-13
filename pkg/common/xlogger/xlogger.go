package xlogger

import "gopkg.in/natefinch/lumberjack.v2"

func NewLogger(loggerConfig *LoggerConfig) *lumberjack.Logger {
	var logger = &lumberjack.Logger{
		Filename:   loggerConfig.Filename,   // 日志文件路径和名称
		MaxSize:    loggerConfig.MaxSize,    // 单个日志文件的最大大小（以 MB 为单位）
		MaxBackups: loggerConfig.MaxBackups, // 最多保留的旧日志文件数量
		MaxAge:     loggerConfig.MaxAge,     // 保留的旧日志文件的最大天数
		Compress:   loggerConfig.Compress,   // 是否压缩旧日志文件
	}
	return logger
}
