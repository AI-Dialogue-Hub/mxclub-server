package po

import (
	"gorm.io/gorm"
	"mxclub/pkg/common/xmysql"
	"time"
)

type Product struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Type             uint        `gorm:"type:int;not null" json:"type"`
	Title            string      `gorm:"type:varchar(255);not null" json:"title"`
	Price            float64     `gorm:"type:decimal(10,2);not null" json:"price"`
	DiscountRuleID   int         `gorm:"type:int" json:"discount_rule_id"`
	DiscountPrice    float64     `gorm:"type:decimal(10,2)" json:"discount_price"`
	FinalPrice       float64     `gorm:"type:decimal(10,2)" json:"final_price"`
	Description      string      `gorm:"type:text;not null" json:"description"`
	ShortDescription string      `gorm:"type:varchar(255);not null" json:"short_description"`
	Images           xmysql.JSON `gorm:"type:json;not null" json:"images"`
	DetailImages     xmysql.JSON `gorm:"type:json;not null" json:"detail_images"`
	Thumbnail        string
	CreatedAt        time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Product) TableName() string {
	return "products"
}
