package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"os"
	"testing"
)

func TestProxy(t *testing.T) {
	xlog.Errorf("%v", os.Getenv("HTTP_PROXY"))
	xlog.Errorf("%v", os.Getenv("HTTPS_PROXY"))
}
