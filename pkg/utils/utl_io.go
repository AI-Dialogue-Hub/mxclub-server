package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"io"
)

func HandleClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		xlog.Errorf("close error: %v", err)
	}
}
