package xredis

var (
	cli RedisIface
)

func NewRedisClient(cfg *RedisConfig) RedisIface {
	if cfg.Single == true {
		cli = NewRedisSingle(cfg)
		return cli
	}
	cli = NewRedisCluster(cfg)
	return cli
}
