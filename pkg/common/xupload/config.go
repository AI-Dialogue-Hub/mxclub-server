package xupload

import (
	"mxclub/pkg/common/xupload/local"
	"mxclub/pkg/common/xupload/xoss"
)

type UploadConfig struct {
	StorageType string        `yaml:"storage_type" validate:"required"`
	LocalConfig *local.Config `yaml:"local_config"`
	OssConfig   *xoss.Config  `yaml:"oss_config"`
}

func SetUp(cfg *UploadConfig) {
	local.SetUpConfig(cfg.LocalConfig)
	xoss.NewClient(cfg.OssConfig)
}
