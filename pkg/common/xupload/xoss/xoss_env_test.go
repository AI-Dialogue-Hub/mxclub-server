package xoss

import (
	"context"
	"log"
	"mxclub/pkg/utils"
	"os"
	"testing"
)

type TestConfig struct {
	Oss *Config `yaml:"oss"`
}

var (
	testConfigPath = "E:\\workspace\\goland\\config\\mxclub.yml"
	testConfig     = new(TestConfig)
	ctx            = context.Background()
	testFile       = "E:\\workspace\\goland\\config\\file\\2024063017190050844.jpg"
)

func setUp() {
	// 读取配置文件
	if err := utils.YamlToStruct(testConfigPath, testConfig); err != nil {
		log.Fatalf("config parse error:%v", err.Error())
	}
}

func TestEnv(t *testing.T) {
	setUp()

	client = NewClient(testConfig.Oss)
	appendFile, err := client.AppendFile(ctx, testConfig.Oss.Bucket, "mxclub/2024063017190050844.jpg")
	if err != nil {
		t.Fatalf("failed to append file %v", err)
	}
	readFile, _ := os.ReadFile(testFile)
	t.Logf("append file af:%#v\n", appendFile)
	n, err := appendFile.Write(readFile)
	if err != nil {
		log.Fatalf("failed to af write %v", err)
	}
	defer appendFile.Close()
	log.Printf("af write n:%#v\n", n)
}
