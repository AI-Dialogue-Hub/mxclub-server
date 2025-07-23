package utils

import "testing"

func TestGzipCompress(t *testing.T) {
	data := "hello world"
	compress, _ := GzipCompress(data)
	t.Logf("data: %s", MustGzipCompressToString(data))
	t.Logf("data: %s", compress)
}
