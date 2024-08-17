package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"mxclub/apps/mxclub-mini/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DecodeBase64ToImage decodes a Base64 string to an image and saves it to the specified file path.
func DecodeBase64ToImage(base64String string) (string, error) {
	// Split the base64 string into header and data parts (if present)
	parts := strings.Split(base64String, ",")
	var data string
	if len(parts) > 1 {
		data = parts[1]
	} else {
		data = parts[0]
	}

	// Decode the Base64 string
	imgData, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return "", err
	}

	// 生成文件名
	fileName := generateFileName("output.png")

	// 拼接文件保存路径
	savePath := filepath.Join(config.GetConfig().File.FilePath, fileName)

	// Write the decoded data to a file
	err = os.WriteFile(savePath, imgData, 0644)
	if err != nil {
		return "", err
	}

	// 构建文件URI
	fileURI := fmt.Sprintf("%s%s", config.GetConfig().File.Domain, fileName)

	return fileURI, nil
}

func generateFileName(originalName string) string {
	// rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100000)
	currentTime := time.Now().Format("20060102150405")

	fileExt := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%s%d%s", currentTime, randomNumber, fileExt)

	return fileName
}
