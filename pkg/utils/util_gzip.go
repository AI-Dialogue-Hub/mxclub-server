package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"io"
	"sync"
)

var (
	gzipLogger     = xlog.NewWith("gzip")
	gzipWriterPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
)

// GzipCompress 压缩数据（支持[]byte和string）
func GzipCompress(data interface{}) ([]byte, error) {
	input, err := anyToByte(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(input); err != nil {
		return nil, fmt.Errorf("gzip write failed: %w", err)
	}

	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("gzip close failed: %w", err)
	}

	return buf.Bytes(), nil
}

func anyToByte(data interface{}) ([]byte, error) {
	var input []byte
	switch v := data.(type) {
	case string:
		input = []byte(v)
	case []byte:
		input = v
	default:
		return nil, fmt.Errorf("unsupported type: %T", data)
	}
	return input, nil
}

func GzipCompressOptimized(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzipWriterPool.Get().(*gzip.Writer)
	gz.Reset(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	gzipWriterPool.Put(gz)
	return buf.Bytes(), nil
}

func MustGzipCompressToString(data interface{}) string {
	defer RecoverByPrefix(gzipLogger, "MustGzipCompressToString")
	defer TraceElapsedWithPrefix(gzipLogger, "MustGzipCompressToString")()
	input, err := anyToByte(data)
	if err != nil {
		return ""
	}

	compress, err := GzipCompressOptimized(input)
	if err != nil {
		gzipLogger.Errorf("GzipCompress error: %v", err)
		return ""
	}
	return string(compress)
}

// GzipDecompress 解压数据（返回[]byte）
func GzipDecompress(compressed []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressed)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("gzip new reader failed: %w", err)
	}
	defer func(gz *gzip.Reader) {
		err = gz.Close()
		if err != nil {

		}
	}(gz)

	var res bytes.Buffer
	if _, err := io.Copy(&res, gz); err != nil {
		return nil, fmt.Errorf("gzip copy failed: %w", err)
	}

	return res.Bytes(), nil
}

func MustGzipDecompress(compressed []byte) []byte {
	defer RecoverByPrefix(gzipLogger, "MustGzipDecompress")
	defer TraceElapsedWithPrefix(gzipLogger, "MustGzipDecompress")()
	res, err := GzipDecompress(compressed)
	if err != nil {
		gzipLogger.Errorf("GzipCompress error: %v", err)
		return []byte{}
	}
	return res
}
