package xredis

type RedisConfig struct {
	Address  []string `yaml:"address"`
	Db       int      `yaml:"db"`
	Password string   `yaml:"password"`
	Prefix   string   `yaml:"prefix"`
	Single   bool     `yaml:"single"`
}
