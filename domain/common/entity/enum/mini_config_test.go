package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMiniConfigEnum_IsValid(t *testing.T) {
	assert.True(t, MiniConfigEnum("swiper").IsValid())
}
