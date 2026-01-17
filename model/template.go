package model

import (
	"time"

	"gorm.io/gorm"
)

// TemplateStatus 模板状态
type TemplateStatus string

const (
	TemplateStatusDraft    TemplateStatus = "draft"
	TemplateStatusEnabled  TemplateStatus = "enabled"
	TemplateStatusDisabled TemplateStatus = "disabled"
)

// Template 邮件模板
type Template struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TriggerCode string         `json:"trigger_code" gorm:"size:100;not null;index:idx_trigger"`
	Language    string         `json:"language" gorm:"size:10;not null;default:zh-CN"`
	Name        string         `json:"name" gorm:"size:200;not null"`
	Subject     string         `json:"subject" gorm:"size:500;not null"`
	BodyHTML    string         `json:"body_html" gorm:"type:text;not null"`
	BodyText    string         `json:"body_text" gorm:"type:text"`
	Status      TemplateStatus `json:"status" gorm:"size:20;not null;default:draft;index:idx_status"`
	Cc          string         `json:"cc" gorm:"size:1000"`
	Bcc         string         `json:"bcc" gorm:"size:1000"`
	ReplyTo     string         `json:"reply_to" gorm:"size:200"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 表名
func (Template) TableName() string {
	return "email_templates"
}

// IsEnabled 是否启用
func (t *Template) IsEnabled() bool {
	return t.Status == TemplateStatusEnabled
}
