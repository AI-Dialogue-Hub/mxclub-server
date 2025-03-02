package po

import (
	"gorm.io/gorm"
	"time"
)

// ProductSale 定义产品销售记录结构体
type ProductSale struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint      `gorm:"type:int" json:"product_id"`                                // 产品ID
	SalesVolume uint      `gorm:"type:int" json:"sales_volume"`                              // 销售量
	SaleDate    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"sale_date"` // 销售日期

	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`                             // 创建时间
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 返回表名
func (ProductSale) TableName() string {
	return "product_sales"
}
