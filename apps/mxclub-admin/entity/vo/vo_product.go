package vo

import "mxclub/pkg/common/xmysql"

type ProductVO struct {
	ID               uint64             `json:"id,omitempty"`
	Type             uint               `json:"type"`
	Title            string             `json:"title,omitempty"`
	Price            float64            `json:"price,omitempty"`
	DiscountRuleID   int                `json:"discount_rule_id,omitempty"`
	DiscountPrice    float64            `json:"discount_price,omitempty"`
	FinalPrice       float64            `json:"final_price,omitempty"`
	Description      string             `json:"description,omitempty"`
	ShortDescription string             `json:"short_description,omitempty"`
	Images           xmysql.StringArray `json:"images,omitempty"`
	DetailImages     xmysql.StringArray `json:"detail_images"`
	Thumbnail        string             `json:"thumbnail"`
	IsHot            bool               `json:"isHot"`
	Visible          bool               `json:"visible"`
	Sale             int                `json:"sale"`
}
