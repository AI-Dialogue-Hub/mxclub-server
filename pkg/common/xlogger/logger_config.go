package xlogger

type LoggerConfig struct {
	Filename   string `yaml:"filename" validate:"required"`    // 日志文件路径和名称
	MaxSize    int    `yaml:"max_size" validate:"required"`    // 单个日志文件的最大大小（以 MB 为单位）
	MaxBackups int    `yaml:"max_backups" validate:"required"` // 最多保留的旧日志文件数量
	MaxAge     int    `yaml:"max_age" validate:"required"`     // 保留的旧日志文件的最大天数
	Compress   bool   `yaml:"compress" validate:"required"`    // 是否压缩旧日志文件
}
