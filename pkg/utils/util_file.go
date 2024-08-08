package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

func GenerateFileName(originalName string) string {
	randomNumber := rand.Intn(100000)
	currentTime := time.Now().Format("20060102150405")

	fileExt := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%s%d%s", currentTime, randomNumber, fileExt)

	return fileName
}
