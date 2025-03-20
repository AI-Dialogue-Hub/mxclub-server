package po

import "time"

// SubNotifyRecord 定义了 notify_records 表的结构
type SubNotifyRecord struct {
	ID         uint       `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID     uint       `gorm:"not null" json:"user_id"`                                  // 用户id
	TemplateID string     `gorm:"size:100;not null" json:"template_id"`                     // 用户订阅消息模板Id
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`              // 创建时间
	UpdatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP;->;<-:update" json:"updated_at"` // 更新时间
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`                                     // 删除时间，使用指针以支持软删除
}

// TableName sub_notify_records
func (SubNotifyRecord) TableName() string {
	return "sub_notify_records"
}
