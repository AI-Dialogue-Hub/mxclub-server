package xmysql

type MySqlConfig struct {
	Address      string `yaml:"address" validate:"required"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Db           string `yaml:"db"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
	MaxIdleTime  int    `yaml:"max_idle_time"`
	Charset      string `yaml:"charset"`
	LogLevel     int    `yaml:"log_level"`
}
