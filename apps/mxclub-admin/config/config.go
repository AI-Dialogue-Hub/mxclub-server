package config

import (
	"flag"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xlogger"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/common/xupload"
	"mxclub/pkg/utils"
	"strings"
	"sync"
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
	// log
	//initLogger(config.LoggerConfig)
	// mysql
	if db, err := xmysql.ConnectDB(config.Mysql); err != nil {
		panic(err)
	} else {
		// gorm
		jet.Provide(func() *gorm.DB { return db })
	}
	// redis
	xredis.NewRedisClient(config.Redis)
	// 加载wxpay
	wxpay.InitWxPay(config.WxPayConfig)
	// oss or localStorage
	xupload.SetUp(config.UploadConfig)
}

type Config struct {
	Server *Server             `yaml:"server" validate:"required"`
	File   File                `yaml:"file" validate:"required"`
	Mysql  *xmysql.MySqlConfig `yaml:"mysql" validate:"required"`
	Redis  *xredis.RedisConfig `yaml:"redis" validate:"required"`

	WxPayConfig  *wxpay.WxPayConfig    `yaml:"wx_pay_config" validate:"required"`
	UploadConfig *xupload.UploadConfig `yaml:"upload_config" validate:"required"`

	LoggerConfig *xlogger.LoggerConfig `yaml:"logger_config" validate:"required"`
}

type Server struct {
	Port    string   `yaml:"port" validate:"required"`
	JwtKey  string   `yaml:"jwt_key" validate:"required"`
	OpenApi []string `yaml:"open_api"`
}

type File struct {
	Domain   string `yaml:"domain" validate:"required"`
	FilePath string `yaml:"file_path" validate:"required"`
}

func GetConfig() *Config {
	return config
}

var (
	openApiSet = make(map[string]bool)
	mu         = new(sync.Mutex)
)

func IsOpenApi(url string) bool {
	mu.Lock()
	defer mu.Unlock()
	if config.Server.OpenApi == nil {
		return false
	}

	if len(openApiSet) == 0 {
		for _, path := range config.Server.OpenApi {
			if strings.HasSuffix(path, "/*") {
				prefix := strings.TrimSuffix(path, "/*")
				openApiSet[prefix] = true
			} else {
				openApiSet[path] = true
			}
		}
	}

	for prefix := range openApiSet {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}

func initLogger(loggerConfig *xlogger.LoggerConfig) {
	if loggerConfig == nil {
		panic("LoggerConfig is invalid")
	}
	xlog.SetGlobalOutput(xlogger.NewLogger(loggerConfig))
}
