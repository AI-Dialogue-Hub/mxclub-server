package req

import (
	"mxclub/pkg/common/xmysql"
)

type ProductReq struct {
	ID               uint64             `json:"id,omitempty"`
	Type             uint               `json:"type,omitempty"`
	Title            string             `json:"title,omitempty" validate:"required" reg_error_info:"不能为空"`
	Price            float64            `json:"price,omitempty" validate:"gt=0" reg_error_info:"原价不能为0"`
	DiscountRuleID   int                `json:"discount_rule_id,omitempty"`
	DiscountPrice    float64            `json:"discount_price,omitempty"`
	FinalPrice       float64            `json:"final_price,omitempty"`
	Description      string             `json:"description,omitempty"`
	ShortDescription string             `json:"short_description,omitempty"`
	Images           xmysql.StringArray `json:"images,omitempty"`
	DetailImages     xmysql.StringArray `json:"detail_images"`
	Thumbnail        string             `json:"thumbnail"`
}

type ProductHotReq struct {
	ID      uint64 `json:"id,omitempty"`
	IsHot   bool   `json:"isHot"`
	Visible bool   `json:"visible"`
}

// ProductListReq get嵌套结构体解析不出来
type ProductListReq struct {
	Page        int64 `json:"page" form:"page" validate:"gt=0" reg_error_info:"参数有误"`           // 页码
	PageSize    int64 `json:"page_size" form:"page_size" validate:"gt=0" reg_error_info:"参数有误"` // 分页大小
	ProductType uint  `form:"product_type"`
}

type ProductSaleReq struct {
	ProductId uint `json:"product_id"  validate:"gt=0" reg_error_info:"参数有误"`
	Sale      int  `json:"sale"`
}
