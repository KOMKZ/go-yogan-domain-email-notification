package email_notification

import "github.com/KOMKZ/go-yogan-framework/errcode"

const (
	ModuleCode = 28 // 邮件通知模块码（20:option, 21:storage, 22:article, 23:auth, 24:admin, 25:member, 26:rbac, 27:docs, 28:email_notification）
)

var (
	ErrTriggerNotFound   = errcode.Register(errcode.New(ModuleCode, 1001, "email_notification", "trigger.not_found", "触发点不存在", 404))
	ErrTemplateNotFound  = errcode.Register(errcode.New(ModuleCode, 1002, "email_notification", "template.not_found", "邮件模板不存在", 404))
	ErrTemplateExists    = errcode.Register(errcode.New(ModuleCode, 1003, "email_notification", "template.exists", "该触发点和语言的模板已存在", 400))
	ErrTemplateDisabled  = errcode.Register(errcode.New(ModuleCode, 1004, "email_notification", "template.disabled", "邮件模板已禁用", 400))
	ErrTemplateRender    = errcode.Register(errcode.New(ModuleCode, 1005, "email_notification", "template.render_failed", "模板渲染失败", 500))
	ErrNoRecipient       = errcode.Register(errcode.New(ModuleCode, 1006, "email_notification", "no_recipient", "收件人不能为空", 400))
	ErrSendFailed        = errcode.Register(errcode.New(ModuleCode, 1007, "email_notification", "send_failed", "邮件发送失败", 500))
	ErrInvalidInput      = errcode.Register(errcode.New(ModuleCode, 1008, "email_notification", "invalid_input", "输入参数无效", 400))
	ErrDatabaseError     = errcode.Register(errcode.New(ModuleCode, 1009, "email_notification", "database_error", "数据库操作失败", 500))
	ErrNotImplemented    = errcode.Register(errcode.New(ModuleCode, 1010, "email_notification", "not_implemented", "功能暂未实现", 501))
	ErrSendLogNotFound       = errcode.Register(errcode.New(ModuleCode, 1011, "email_notification", "send_log.not_found", "发送日志不存在", 404))
	ErrServiceNotAvailable   = errcode.Register(errcode.New(ModuleCode, 1012, "email_notification", "service.not_available", "邮件通知服务不可用", 503))
)
