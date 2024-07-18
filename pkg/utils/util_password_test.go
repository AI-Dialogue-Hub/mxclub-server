package utils

import (
	"strconv"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	t.Logf("%v", EncryptPassword(strconv.Itoa(123456)))
}
