package enum

var DiscountRules = map[string]struct {
	Threshold float64
	Discount  float64
}{
	"LV1": {Threshold: 500, Discount: 0.98},
	"LV2": {Threshold: 2000, Discount: 0.95},
	"LV3": {Threshold: 5000, Discount: 0.92},
	"LV4": {Threshold: 10000, Discount: 0.90},
	"LV5": {Threshold: 20000, Discount: 0.88},
	"LV6": {Threshold: 50000, Discount: 0.85},
}

func FetchDiscountByGrade(grade string) float64 {
	if discountRule, ok := DiscountRules[grade]; ok {
		return discountRule.Discount
	}
	return 1
}
