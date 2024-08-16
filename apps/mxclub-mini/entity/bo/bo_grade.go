package bo

import (
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"math"
	"mxclub/pkg/utils"
)

var GradeMap = utils.NewLinkedHashMapWithPairs([]*maps.Pair[float64, string]{
	{500, "LV1"},
	{2000, "LV2"},
	{5000, "LV3"},
	{10000, "LV4"},
	{20000, "LV5"},
	{50000, "LV6"},
	{math.MaxFloat64, "LV7"},
})

func GetGradeByScore(score float64) string {
	keySet := GradeMap.KeySet()
	for index, level := range keySet {
		if score >= level && score < keySet[index+1] {
			if got, ok := GradeMap.Get(level); ok {
				return got
			}
		}
	}
	return "LV0"
}
