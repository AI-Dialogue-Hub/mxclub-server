package wxpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"math/big"
	mathRand "math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// generateNonceStr 生成32位随机字符串
func generateNonceStr() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return xlog.GenReqId()
	}
	return hex.EncodeToString(bytes)
}

var (
	block   *pem.Block
	blockMu sync.Mutex
)

// getRSASignature 生成RSA签名
func getRSASignature(parts []string) (string, error) {
	// 拼接字符串
	str := ""
	for _, part := range parts {
		str += part + "\n"
	}

	if block == nil {
		// 解析私钥
		if _, err := loadPrivateKeyBlock(wxPayConfig.PrivateKeyPath); err != nil {
			return "", errors.New("failed to parse PEM block containing the key")
		}
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPrivateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// 生成签名
	hashed := sha256.Sum256([]byte(str))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// loadPrivateKeyBlock 解析并缓存私钥的 PEM block
func loadPrivateKeyBlock(privateKeyPath string) (*pem.Block, error) {
	blockMu.Lock()
	defer blockMu.Unlock()

	// 如果缓存的路径不同或 block 为 nil，需要重新加载
	if block == nil || privateKeyPath != privateKeyPath {
		privateKeyBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}

		newBlock, _ := pem.Decode(privateKeyBytes)
		if newBlock == nil {
			return nil, errors.New("failed to parse PEM block containing the key")
		}
		block = newBlock
	}

	return block, nil
}

func init() {
	// 生成一个随机数作为后缀
	mathRand.Seed(time.Now().Unix())
}

func GenerateUniqueOrderNumber() string {

	// 获取当前时间戳
	timestamp := time.Now().Unix()

	randomNumber := mathRand.Intn(1000000) // 包含0到999999

	// 使用fmt.Sprintf来保证四位数，不足四位则前面补零
	suffix := fmt.Sprintf("%06d", randomNumber)

	// 拼接订单号
	orderNumber := strconv.FormatInt(timestamp, 10) + suffix

	return orderNumber
}

func GenerateOutRefundNo() string {
	return GenerateRandomString(25)
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]rune, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[n.Int64()]
	}
	return string(result)
}
