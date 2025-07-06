//go:build ignore
// +build ignore

package wxwork

import (
	"log"
	"mxclub/pkg/utils"
	"testing"
)

type testConfigDTO struct {
	Config *WxWorkConfig `yaml:"wx_work_config"`
}

var (
	testConfigPath = "E:\\workspace\\goland\\config\\mxclub.yml"
	testConfig     = new(testConfigDTO)
	svc            *WxworkService
)

func setUp() {
	// 读取配置文件
	if err := utils.YamlToStruct(testConfigPath, testConfig); err != nil {
		log.Fatalf("config parse error:%v", err.Error())
	}
	svc = NewWxworkService(testConfig.Config)
}

func init() {
	setUp()
}

func TestWxworkService_GetAccessToken(t *testing.T) {
	token, err := svc.GetAccessToken()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("token is %v", token)
}
