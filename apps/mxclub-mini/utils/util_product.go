package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
)

// 预编译正则表达式
var (
	imgTagRegex     = regexp.MustCompile(`<img[^>]*>`)
	brTagRegex      = regexp.MustCompile(`<br[^>]*/?>`)
	imgStartRegex   = regexp.MustCompile(`<img`)
	pTagStartRegex  = regexp.MustCompile(`<p([^>]*)>`)
	styleAttrRegex  = regexp.MustCompile(`style\s*=\s*["'][^"']*["']`)
	widthAttrRegex  = regexp.MustCompile(`\s+width\s*=\s*["'][^"']*["']`)
	heightAttrRegex = regexp.MustCompile(`\s+height\s*=\s*["'][^"']*["']`)

	// 匹配style中的width属性，保留前后的内容
	styleWidthRegex = regexp.MustCompile(`(style\s*=\s*["'][^"']*)width\s*:\s*[^;]+;([^"']*["'])`)

	// 匹配已有的style属性，用于检查是否已有style
	existingStyleRegex = regexp.MustCompile(`style\s*=\s*["']([^"']*)["']`)
)

func RepairRichTextEnhanced(html string) string {
	defer utils.TraceElapsedByName(time.Now(), "RepairRichTextEnhanced")
	if html == "" {
		return ""
	}

	result := html

	// 1. 处理img标签
	result = imgTagRegex.ReplaceAllStringFunc(result, func(match string) string {
		// 移除不需要的属性
		match = styleAttrRegex.ReplaceAllString(match, "")
		match = widthAttrRegex.ReplaceAllString(match, "")
		match = heightAttrRegex.ReplaceAllString(match, "")
		return match
	})

	// 2. 处理style中的width属性
	result = styleWidthRegex.ReplaceAllStringFunc(result, func(match string) string {
		// 将width替换为max-width:100%
		return regexp.MustCompile(`width\s*:\s*[^;]+;`).ReplaceAllString(match, "max-width:100%;")
	})

	// 3. 移除<br/>标签
	//result = brTagRegex.ReplaceAllString(result, "")

	// 4. 为img标签添加样式
	result = imgStartRegex.ReplaceAllStringFunc(result, func(match string) string {
		return `<img style="max-width:100%;height:auto;display:block;margin-top:0;margin-bottom:0;"`
	})

	// 5. 可选：为p标签添加样式（如果原始JavaScript中有这个功能）
	result = pTagStartRegex.ReplaceAllStringFunc(result, func(match string) string {
		// 检查是否已有style属性
		if existingStyleRegex.MatchString(match) {
			// 在现有style中添加
			return existingStyleRegex.ReplaceAllString(match, `style="$1;margin:16px 0;line-height:1.6;"`)
		}
		// 添加新的style属性
		return `<p style="margin:16px 0;line-height:1.6;">`
	})

	return result
}

var (
	// 检查是否包含HTML标签
	// 常见的HTML标签正则
	htmlTagRegex = regexp.MustCompile(`<(?i)(p|div|span|a|img|table|tr|td|th|ul|ol|li|h[1-6]|br|hr|strong|b|em|i|u|strike|font|blockquote|pre|code)[^>]*>`)

	// 检查闭合标签
	closingTagRegex = regexp.MustCompile(`</(?i)(p|div|span|a|img|table|tr|td|th|ul|ol|li|h[1-6]|br|hr|strong|b|em|i|u|strike|font|blockquote|pre|code)[^>]*>`)

	// 检查自闭合标签
	selfClosingTagRegex = regexp.MustCompile(`<(?i)(img|br|hr|input|meta|link|base)[^>]*\/>`)

	// 检查HTML实体
	htmlEntityRegex = regexp.MustCompile(`&(?:[a-zA-Z]+|#\d+);`)

	// 检查是否包含HTML属性
	htmlAttrRegex = regexp.MustCompile(`\s+(?i)(class|id|style|href|src|alt|title|width|height|border|cellpadding|cellspacing)\s*=\s*["'][^"']*["']`)

	// 检查CSS样式
	cssStyleRegex = regexp.MustCompile(`style\s*=\s*["'][^"']*["']`)
)

// IsRichText 检查字符串是否为富文本HTML
// 返回 true 如果是富文本，false 如果是纯文本或空
func IsRichText(content string) bool {
	if content == "" {
		return false
	}

	// 去除首尾空白
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return false
	}

	// 如果是富文本，应该至少满足以下条件之一：
	return htmlTagRegex.MatchString(trimmed) ||
		closingTagRegex.MatchString(trimmed) ||
		selfClosingTagRegex.MatchString(trimmed) ||
		htmlEntityRegex.MatchString(trimmed) ||
		htmlAttrRegex.MatchString(trimmed) ||
		cssStyleRegex.MatchString(trimmed)
}
