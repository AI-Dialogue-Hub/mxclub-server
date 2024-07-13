package captcha

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"github.com/mojocn/base64Captcha"
	"time"
)

// StringCaptcha 验证码工具类
type StringCaptcha struct {
	captcha *base64Captcha.Captcha
}

// NewCaptcha 创建验证码
func NewCaptcha() *StringCaptcha {
	// store
	store := base64Captcha.DefaultMemStore

	// 包含数字和字母的字符集
	source := "123456789"

	// driver
	driver := base64Captcha.NewDriverString(
		80,     // height int
		240,    // width int
		6,      // noiseCount int
		1,      // showLineOptions int
		4,      // length int
		source, // source string
		nil,    // bgColor *color.RGBA
		nil,    // fontsStorage
		nil,    // fonts []string
	)

	captcha := base64Captcha.NewCaptcha(driver, store)
	return &StringCaptcha{
		captcha: captcha,
	}
}

// Generate 生成验证码
func (stringCaptcha *StringCaptcha) Generate() (string, string, string) {
	defer utils.TraceElapsedByName(time.Now(), "[stringCaptcha]Generate")
	id, b64s, answer, _ := stringCaptcha.captcha.Generate()
	return id, b64s, answer
}

// Verify 验证验证码
func (stringCaptcha *StringCaptcha) Verify(id string, answer string) bool {
	return stringCaptcha.captcha.Verify(id, answer, true)
}
