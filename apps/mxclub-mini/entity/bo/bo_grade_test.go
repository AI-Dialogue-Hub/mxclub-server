package bo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGradeByScore(t *testing.T) {
	cases := []struct {
		Expected string
		Score    float64
	}{
		{Expected: "LV1", Score: 500 - 1},
		{Expected: "LV2", Score: 2000 - 1},
		{Expected: "LV3", Score: 5000 - 1},
		{Expected: "LV4", Score: 10000 - 1},
		{Expected: "LV5", Score: 20000 - 1},
		{Expected: "LV6", Score: 50000 - 1},

		// ============================

		{Expected: "LV2", Score: 500},
		{Expected: "LV3", Score: 2000},
		{Expected: "LV4", Score: 5000},
		{Expected: "LV5", Score: 10000},
		{Expected: "LV6", Score: 20000},
	}
	for _, caseInfo := range cases {
		t.Run(caseInfo.Expected, func(tt *testing.T) {
			assert.Equal(tt, caseInfo.Expected, GetGradeByScore(caseInfo.Score))
		})
	}

}
