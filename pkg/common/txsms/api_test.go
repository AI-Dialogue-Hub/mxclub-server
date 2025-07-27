//go:build ignore
// +build ignore

package txsms

import (
	"fmt"
	"log"
	"mxclub/pkg/utils"
	"testing"
)

type testConfigDTO struct {
	Config *TxSmsConfig `yaml:"tx_sms_config" validate:"required"`
}

var (
	testConfigPath = "E:\\workspace\\goland\\Config\\mxclub.yml"
	testConfig     = new(testConfigDTO)
	svc            *TxSmsService
)

func setUp() {
	// 读取配置文件
	if err := utils.YamlToStruct(testConfigPath, testConfig); err != nil {
		log.Fatalf("Config parse error:%v", err.Error())
	}
	svc = NewTxSmsService(testConfig.Config)
}

func init() {
	setUp()
}

func TestTxSmsService_SendDispatchMsg(t *testing.T) {
	err := svc.SendDispatchMsg("17670459756")
	if err != nil {
		fmt.Printf("发送短信失败: %v\n", err)
		return
	}
	fmt.Println("短信发送成功")
}
