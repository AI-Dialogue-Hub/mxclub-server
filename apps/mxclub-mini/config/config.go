package config

import (
	"flag"
	"gorm.io/gorm"
	"log"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/go-playground/validator/v10"
)

var (
	config   = new(Config)
	confFile = flag.String("f", "./configs/dpp_server.yaml", "config file")
)

func init() {
	flag.Parse()
	if err := utils.YamlToStruct(*confFile, config); err != nil {
		log.Fatalf("config parse error:%v", err.Error())
	}
	if err := validator.New().Struct(config); err != nil {
		log.Fatalf("config error:%v", err.Error())
	}
	// mysql
	jet.Provide(func() *xmysql.MySqlConfig { return config.Mysql })
	// ============== 耗时加载全部异步化 =============================
	var (
		c1 = make(chan struct{})
		c2 = make(chan struct{})
	)
	go func() {
		if db, err := xmysql.ConnectDB(config.Mysql); err != nil {
			panic(err)
		} else {
			// gorm
			jet.Provide(func() *gorm.DB { return db })
		}
		c1 <- struct{}{}
	}()
	go func() {
		// redis
		xredis.NewRedisClient(config.Redis)
		c2 <- struct{}{}
	}()
	<-c2
	<-c1
}

type Config struct {
	Server *Server `yaml:"server" validate:"required"`
	File   File    `yaml:"file" validate:"required"`

	Mysql *xmysql.MySqlConfig `yaml:"mysql" validate:"required"`
	Redis *xredis.RedisConfig `yaml:"redis" validate:"required"`
}

type Server struct {
	Port string `yaml:"port"`
}

type File struct {
	Domain             string `yaml:"domain" validate:"required"`
	FilePath           string `yaml:"file_path" validate:"required"`
	MaxRequestBodySize int    `yaml:"max_request_body_size" validate:"lte=1200,gte=1" reg_err_info:"不合法"`
}

func GetConfig() *Config {
	return config
}
