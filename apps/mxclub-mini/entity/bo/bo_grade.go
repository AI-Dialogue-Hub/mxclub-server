package bo

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"mxclub/pkg/utils"
)

var GradeMap = utils.NewLinkedHashMapWithPairs([]*maps.Pair[float64, string]{
	{500, "LV1"},
	{2000, "LV2"},
	{5000, "LV3"},
	{10000, "LV4"},
	{20000, "LV5"},
	{50000, "LV6"},
})

func GetGradeByScore(score float64) string {
	for _, level := range GradeMap.KeySet() {
		if score < level {
			return GradeMap.MustGet(level)
		}
	}
	return "LV0"
}
