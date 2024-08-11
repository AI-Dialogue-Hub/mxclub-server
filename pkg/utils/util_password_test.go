package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	genReqId := xlog.GenReqId()
	t.Logf("%v", genReqId)
	t.Logf("%v", EncryptPassword(genReqId))
}
