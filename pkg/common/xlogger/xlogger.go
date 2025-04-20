package xlogger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func NewLogger(loggerConfig *LoggerConfig) *lumberjack.Logger {
	// 初始化 lumberjack.Logger
	logger := &lumberjack.Logger{
		Filename:   loggerConfig.Filename,   // 日志文件路径和名称
		MaxSize:    loggerConfig.MaxSize,    // 单个日志文件的最大大小（以 MB 为单位）
		MaxBackups: loggerConfig.MaxBackups, // 最多保留的旧日志文件数量
		MaxAge:     loggerConfig.MaxAge,     // 保留的旧日志文件的最大天数
		Compress:   loggerConfig.Compress,   // 是否压缩旧日志文件
	}

	// 创建管道
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// 保存原始的 stderr
	originalStderr := os.Stderr

	// 替换 os.Stderr 为管道的写入端
	os.Stderr = w

	// 启动一个 goroutine 处理管道中的日志
	go func() {
		defer func() {
			// 恢复原始的 os.Stderr
			os.Stderr = originalStderr
			// 关闭管道的读取端
			_ = r.Close()
		}()

		// 使用 MultiWriter 同时写入 lumberjack.Logger 和原始 os.Stderr
		multiWriter := io.MultiWriter(logger, originalStderr)
		_, _ = io.Copy(multiWriter, r)
	}()

	return logger
}
