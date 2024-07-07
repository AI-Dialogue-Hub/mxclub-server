package req

import "mxclub/pkg/common/xmysql"

type ProductReq struct {
	ID               uint64      `json:"id,omitempty"`
	Type             uint        `json:"type,omitempty"`
	Title            string      `json:"title,omitempty" validate:"required" reg_error_info:"不能为空"`
	Price            float64     `json:"price,omitempty" validate:"gt=0" reg_error_info:"原价不能为0"`
	DiscountRuleID   int         `json:"discount_rule_id,omitempty"`
	DiscountPrice    float64     `json:"discount_price,omitempty"`
	FinalPrice       float64     `json:"final_price,omitempty"`
	Description      string      `json:"description,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
	Images           xmysql.JSON `json:"images,omitempty"`
	DetailImages     xmysql.JSON `json:"detail_images"`
}
