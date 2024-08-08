package config

import (
	"flag"
	"gorm.io/gorm"
	"log"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/common/xupload"
	"mxclub/pkg/utils"
	"strings"
	"sync"

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

	if db, err := xmysql.ConnectDB(config.Mysql); err != nil {
		panic(err)
	} else {
		// gorm
		jet.Provide(func() *gorm.DB { return db })
	}
	// redis
	xredis.NewRedisClient(config.Redis)
	// wxPay
	wxpay.InitWxPay(config.WxPayConfig)
	// oss or localStorage
	xupload.SetUp(config.UploadConfig)
}

type Config struct {
	Server       *Server               `yaml:"server" validate:"required"`
	WxConfig     *WxConfig             `yaml:"wx_config" validate:"required"`
	WxPayConfig  *wxpay.WxPayConfig    `yaml:"wx_pay_config" validate:"required"`
	File         File                  `yaml:"file" validate:"required"`
	UploadConfig *xupload.UploadConfig `yaml:"upload_config" validate:"required"`

	Mysql *xmysql.MySqlConfig `yaml:"mysql" validate:"required"`
	Redis *xredis.RedisConfig `yaml:"redis" validate:"required"`
}

type Server struct {
	Port    string   `yaml:"port" validate:"required"`
	JwtKey  string   `yaml:"jwt_key" validate:"required"`
	OpenApi []string `yaml:"open_api"`
}

type File struct {
	Domain             string `yaml:"domain" validate:"required"`
	FilePath           string `yaml:"file_path" validate:"required"`
	MaxRequestBodySize int    `yaml:"max_request_body_size" validate:"lte=1200,gte=1" reg_err_info:"不合法"`
}

type WxConfig struct {
	Ak string `yaml:"ak" validate:"required"`
	Sk string `yaml:"sk" validate:"required"`
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
