package model

import "time"

type SendStatus string

const (
	SendStatusPending SendStatus = "pending"
	SendStatusSent    SendStatus = "sent"
	SendStatusFailed  SendStatus = "failed"
)

type SendLog struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	TemplateID   *uint      `json:"template_id" gorm:"index"`
	TriggerCode  string     `json:"trigger_code" gorm:"size:100;not null;index:idx_trigger"`
	Language     string     `json:"language" gorm:"size:10;not null"`
	Recipient    string     `json:"recipient" gorm:"size:500;not null"`
	Subject      string     `json:"subject" gorm:"size:500;not null"`
	Params       string     `json:"params" gorm:"type:json"`
	Status       SendStatus `json:"status" gorm:"size:20;not null;default:pending;index:idx_status"`
	ErrorMessage string     `json:"error_message" gorm:"type:text"`
	SentAt       *time.Time `json:"sent_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index:idx_created"`
}

func (SendLog) TableName() string {
	return "email_send_logs"
}

func (l *SendLog) MarkSent() {
	now := time.Now()
	l.Status = SendStatusSent
	l.SentAt = &now
}

func (l *SendLog) MarkFailed(errMsg string) {
	l.Status = SendStatusFailed
	l.ErrorMessage = errMsg
}
