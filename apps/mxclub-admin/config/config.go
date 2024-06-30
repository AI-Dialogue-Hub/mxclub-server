package config

import (
	"flag"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
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
	if db, err := xmysql.ConnectDB(config.Mysql); err != nil {
		panic(err)
	} else {
		// gorm
		jet.Provide(func() *gorm.DB { return db })
	}
}

type Config struct {
	Server *Server             `yaml:"server" validate:"required"`
	File   File                `yaml:"file" validate:"required"`
	Mysql  *xmysql.MySqlConfig `yaml:"mysql" validate:"required"`
}

type Server struct {
	Port   string `yaml:"port" validate:"required"`
	JwtKey string `yaml:"jwt_key" validate:"required"`
}

type File struct {
	Domain   string `yaml:"domain" validate:"required"`
	FilePath string `yaml:"file_path" validate:"required"`
}

func GetConfig() *Config {
	return config
}
