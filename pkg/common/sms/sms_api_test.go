package sms

import (
	"log"
	"mxclub/pkg/utils"
	"testing"
)

type testConfigDTO struct {
	Config *Config `yaml:"sms_config"`
}

var (
	testConfigPath = "E:\\workspace\\goland\\config\\mxclub.yml"
	testConfig     = new(testConfigDTO)
	svc            ISmsService
)

func setUp() {
	// 读取配置文件
	if err := utils.YamlToStruct(testConfigPath, testConfig); err != nil {
		log.Fatalf("config parse error:%v", err.Error())
	}
	svc = NewAliSmsService(testConfig.Config)
}

func TestSendSms(t *testing.T) {
	setUp()
	c := testConfig.Config
	err := svc.SendDispatchMsg(c.TestPhone, "测试昵称", "测试角色")
	if err != nil {
		t.Errorf("test SendDispatchMsg failed, err: %v", err)
		t.Failed()
	}
}
