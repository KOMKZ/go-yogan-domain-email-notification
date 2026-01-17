package email_notification

import "github.com/KOMKZ/go-yogan-domain-email-notification/model"

// CreateTemplateInput 创建模板输入
type CreateTemplateInput struct {
	TriggerCode string             `json:"trigger_code"`
	Language    string             `json:"language"`
	Name        string             `json:"name"`
	Subject     string             `json:"subject"`
	BodyHTML    string             `json:"body_html"`
	BodyText    string             `json:"body_text"`
	Status      model.TemplateStatus `json:"status"`
	Cc          string             `json:"cc"`
	Bcc         string             `json:"bcc"`
	ReplyTo     string             `json:"reply_to"`
}

// UpdateTemplateInput 更新模板输入
type UpdateTemplateInput struct {
	Name     *string             `json:"name"`
	Subject  *string             `json:"subject"`
	BodyHTML *string             `json:"body_html"`
	BodyText *string             `json:"body_text"`
	Status   *model.TemplateStatus `json:"status"`
	Cc       *string             `json:"cc"`
	Bcc      *string             `json:"bcc"`
	ReplyTo  *string             `json:"reply_to"`
}

// SendInput 发送输入
type SendInput struct {
	TriggerCode string         // 触发点代码（必填）
	Recipient   string         // 收件人（必填，可逗号分隔多个）
	Language    string         // 语言（可选，默认 zh-CN）
	Params      map[string]any // 参数（通用+Trigger 专属）

	// 以下字段可覆盖模板配置
	Cc          []string     // 抄送（追加到模板配置）
	Bcc         []string     // 密送（追加到模板配置）
	ReplyTo     string       // 回复地址（覆盖模板配置）
	From        string       // 发件人（覆盖默认配置）
	FromName    string       // 发件人名称
	Subject     string       // 主题（覆盖模板，用于特殊场景）
	Attachments []Attachment // 附件
}

// Attachment 附件
type Attachment struct {
	Filename    string // 文件名
	Content     []byte // 文件内容
	ContentType string // MIME 类型
}

// PreviewResult 预览结果
type PreviewResult struct {
	Subject  string `json:"subject"`
	BodyHTML string `json:"body_html"`
	BodyText string `json:"body_text"`
}
