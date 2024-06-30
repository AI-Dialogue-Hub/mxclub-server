package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptPassword(password string) string {
	// 创建 MD5 哈希对象
	hashes := md5.New()

	// 将密码转换为字节数组并进行哈希计算
	hashes.Write([]byte(password))
	encryptedBytes := hashes.Sum(nil)

	// 将哈希结果转换为十六进制字符串
	encryptedPassword := hex.EncodeToString(encryptedBytes)

	return encryptedPassword
}
