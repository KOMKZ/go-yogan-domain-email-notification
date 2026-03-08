package errors

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

const ModuleEmailNotification = 28

var (
	ErrTriggerNotFound = errcode.Register(errcode.New(
		ModuleEmailNotification, 1001, "email_notification",
		"error.email.trigger_not_found", "触发点不存在",
		http.StatusNotFound,
	))
	ErrTemplateNotFound = errcode.Register(errcode.New(
		ModuleEmailNotification, 1002, "email_notification",
		"error.email.template_not_found", "邮件模板不存在",
		http.StatusNotFound,
	))
	ErrTemplateExists = errcode.Register(errcode.New(
		ModuleEmailNotification, 1003, "email_notification",
		"error.email.template_exists", "该触发点和语言的模板已存在",
		http.StatusConflict,
	))
	ErrTemplateDisabled = errcode.Register(errcode.New(
		ModuleEmailNotification, 1004, "email_notification",
		"error.email.template_disabled", "模板已禁用",
		http.StatusBadRequest,
	))
	ErrTemplateRender = errcode.Register(errcode.New(
		ModuleEmailNotification, 1005, "email_notification",
		"error.email.template_render", "模板渲染失败",
		http.StatusInternalServerError,
	))
	ErrNoRecipient = errcode.Register(errcode.New(
		ModuleEmailNotification, 1006, "email_notification",
		"error.email.no_recipient", "收件人不能为空",
		http.StatusBadRequest,
	))
	ErrSendFailed = errcode.Register(errcode.New(
		ModuleEmailNotification, 1007, "email_notification",
		"error.email.send_failed", "邮件发送失败",
		http.StatusInternalServerError,
	))
	ErrInvalidInput = errcode.Register(errcode.New(
		ModuleEmailNotification, 1008, "email_notification",
		"error.email.invalid_input", "输入参数无效",
		http.StatusBadRequest,
	))
	ErrSendLogNotFound = errcode.Register(errcode.New(
		ModuleEmailNotification, 1011, "email_notification",
		"error.email.send_log_not_found", "发送日志不存在",
		http.StatusNotFound,
	))
)
