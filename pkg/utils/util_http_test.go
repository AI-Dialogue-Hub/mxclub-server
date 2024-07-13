package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	xlog.Errorf("%v", os.Getenv("HTTP_PROXY"))
	xlog.Errorf("%v", os.Getenv("HTTPS_PROXY"))
}

func TestCert(t *testing.T) {
	xlog.Infof("%v", GenerateKey())
}

// GenerateKey generates a 32-character key consisting of numbers and uppercase/lowercase letters
func GenerateKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < 32; i++ {
		randomIndex := seededRand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
