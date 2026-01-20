package email_notification

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	email "github.com/KOMKZ/go-yogan-component-email"
	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"gorm.io/gorm"
)

// Service 邮件通知服务
type Service struct {
	db           *gorm.DB
	templateRepo TemplateRepository
	logRepo      SendLogRepository
	emailMgr     *email.Manager
	registry     *TriggerRegistry
	engine       *TemplateEngine
	commonParams map[string]any // 通用参数（应用级注入，Send 时自动合并）
}

// NewService 创建服务
func NewService(db *gorm.DB, emailMgr *email.Manager, registry *TriggerRegistry, commonParams map[string]any) *Service {
	return &Service{
		db:           db,
		templateRepo: NewGormTemplateRepository(db),
		logRepo:      NewGormSendLogRepository(db),
		emailMgr:     emailMgr,
		registry:     registry,
		engine:       NewTemplateEngine(),
		commonParams: commonParams,
	}
}

// ========== Trigger 查询 ==========

// ListTriggers 获取所有触发点
func (s *Service) ListTriggers(ctx context.Context) []*TriggerDefinition {
	return s.registry.GetAll()
}

// GetTrigger 获取触发点详情
func (s *Service) GetTrigger(ctx context.Context, code string) (*TriggerDefinition, error) {
	trigger, ok := s.registry.Get(code)
	if !ok {
		return nil, ErrTriggerNotFound
	}
	return trigger, nil
}

// GetTriggerParams 获取触发点的全部参数
func (s *Service) GetTriggerParams(ctx context.Context, code string) []Param {
	return s.registry.GetAllParams(code)
}

// GetRegistry 获取注册表
func (s *Service) GetRegistry() *TriggerRegistry {
	return s.registry
}

// ========== Template CRUD ==========

// CreateTemplate 创建模板
func (s *Service) CreateTemplate(ctx context.Context, input CreateTemplateInput) (*model.Template, error) {
	// 验证触发点存在
	if !s.registry.Exists(input.TriggerCode) {
		return nil, ErrTriggerNotFound.WithMsg("触发点不存在: " + input.TriggerCode)
	}

	// 设置默认语言
	if input.Language == "" {
		input.Language = "zh-CN"
	}

	// 检查是否已存在
	exists, err := s.templateRepo.ExistsByTriggerAndLanguage(ctx, input.TriggerCode, input.Language, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTemplateExists
	}

	// 设置默认状态
	if input.Status == "" {
		input.Status = model.TemplateStatusDraft
	}

	template := &model.Template{
		TriggerCode: input.TriggerCode,
		Language:    input.Language,
		Name:        input.Name,
		Subject:     input.Subject,
		BodyHTML:    input.BodyHTML,
		BodyText:    input.BodyText,
		Status:      input.Status,
		Cc:          input.Cc,
		Bcc:         input.Bcc,
		ReplyTo:     input.ReplyTo,
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	return template, nil
}

// UpdateTemplate 更新模板
func (s *Service) UpdateTemplate(ctx context.Context, id uint, input UpdateTemplateInput) (*model.Template, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		template.Name = *input.Name
	}
	if input.Subject != nil {
		template.Subject = *input.Subject
	}
	if input.BodyHTML != nil {
		template.BodyHTML = *input.BodyHTML
	}
	if input.BodyText != nil {
		template.BodyText = *input.BodyText
	}
	if input.Status != nil {
		template.Status = *input.Status
	}
	if input.Cc != nil {
		template.Cc = *input.Cc
	}
	if input.Bcc != nil {
		template.Bcc = *input.Bcc
	}
	if input.ReplyTo != nil {
		template.ReplyTo = *input.ReplyTo
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	return template, nil
}

// DeleteTemplate 删除模板
func (s *Service) DeleteTemplate(ctx context.Context, id uint) error {
	return s.templateRepo.Delete(ctx, id)
}

// GetTemplate 获取模板详情
func (s *Service) GetTemplate(ctx context.Context, id uint) (*model.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// ListTemplates 模板列表
func (s *Service) ListTemplates(ctx context.Context, filter TemplateFilter) (*PageResult[model.Template], error) {
	return s.templateRepo.List(ctx, filter)
}

// GetTemplateByTrigger 获取指定触发点和语言的模板
func (s *Service) GetTemplateByTrigger(ctx context.Context, triggerCode, language string) (*model.Template, error) {
	return s.templateRepo.GetActiveTemplate(ctx, triggerCode, language)
}

// ========== 预览与测试 ==========

// PreviewTemplate 预览模板
func (s *Service) PreviewTemplate(ctx context.Context, id uint) (*PreviewResult, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	params := s.registry.GetAllParams(template.TriggerCode)

	subject, err := s.engine.Preview(template.Subject, params)
	if err != nil {
		return nil, err
	}

	bodyHTML, err := s.engine.Preview(template.BodyHTML, params)
	if err != nil {
		return nil, err
	}

	bodyText := ""
	if template.BodyText != "" {
		bodyText, err = s.engine.Preview(template.BodyText, params)
		if err != nil {
			return nil, err
		}
	}

	return &PreviewResult{
		Subject:  subject,
		BodyHTML: bodyHTML,
		BodyText: bodyText,
	}, nil
}

// TestSend 测试发送
func (s *Service) TestSend(ctx context.Context, id uint, recipient string) error {
	if recipient == "" {
		return ErrNoRecipient
	}

	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	params := s.registry.GetAllParams(template.TriggerCode)

	// 使用示例值渲染
	exampleParams := make(map[string]any)
	for _, p := range params {
		if p.Example != "" {
			exampleParams[p.Name] = p.Example
		} else {
			exampleParams[p.Name] = "{{." + p.Name + "}}"
		}
	}
	exampleParams["CurrentYear"] = time.Now().Year()

	return s.sendWithTemplate(ctx, template, recipient, exampleParams, nil)
}

// ========== 发送 ==========

// Send 同步发送邮件
func (s *Service) Send(ctx context.Context, input SendInput) error {
	// 验证输入
	if input.TriggerCode == "" {
		return ErrInvalidInput.WithMsg("触发点代码不能为空")
	}
	if input.Recipient == "" {
		return ErrNoRecipient
	}

	// 验证触发点存在
	if !s.registry.Exists(input.TriggerCode) {
		return ErrTriggerNotFound.WithMsg("触发点不存在: " + input.TriggerCode)
	}

	// 确定语言
	language := input.Language
	if language == "" {
		language = "zh-CN"
	}

	// 获取启用的模板
	template, err := s.templateRepo.GetActiveTemplate(ctx, input.TriggerCode, language)
	if err != nil {
		// 尝试回退到默认语言
		if language != "zh-CN" {
			template, err = s.templateRepo.GetActiveTemplate(ctx, input.TriggerCode, "zh-CN")
		}
		if err != nil {
			return err
		}
	}

	// 合并参数
	params := s.mergeParams(input.Params)

	return s.sendWithTemplate(ctx, template, input.Recipient, params, &input)
}

// SendAsync 异步发送邮件（预留，暂不实现）
func (s *Service) SendAsync(ctx context.Context, input SendInput) error {
	return ErrNotImplemented.WithMsg("异步发送暂未实现")
}

// sendWithTemplate 使用模板发送邮件
func (s *Service) sendWithTemplate(ctx context.Context, template *model.Template, recipient string, params map[string]any, input *SendInput) error {
	// 渲染主题
	subject, err := s.engine.Render(template.Subject, params)
	if err != nil {
		return err
	}

	// 渲染正文
	body, err := s.engine.Render(template.BodyHTML, params)
	if err != nil {
		return err
	}

	// 创建发送日志
	paramsJSON, _ := json.Marshal(params)
	sendLog := &model.SendLog{
		TemplateID:  &template.ID,
		TriggerCode: template.TriggerCode,
		Language:    template.Language,
		Recipient:   recipient,
		Subject:     subject,
		Params:      string(paramsJSON),
		Status:      model.SendStatusPending,
	}
	if err := s.logRepo.Create(ctx, sendLog); err != nil {
		return ErrDatabaseError.Wrap(err)
	}

	// 构建邮件
	builder := s.emailMgr.New().
		To(recipient).
		Subject(subject).
		Body(body)

	// 发件人
	if input != nil && input.From != "" {
		builder.From(input.From)
	}
	if input != nil && input.FromName != "" {
		builder.FromName(input.FromName)
	}

	// 抄送：模板配置 + input 追加
	if template.Cc != "" {
		for _, cc := range strings.Split(template.Cc, ",") {
			cc = strings.TrimSpace(cc)
			if cc != "" {
				builder.Cc(cc)
			}
		}
	}
	if input != nil {
		for _, cc := range input.Cc {
			builder.Cc(cc)
		}
	}

	// 密送：模板配置 + input 追加
	if template.Bcc != "" {
		for _, bcc := range strings.Split(template.Bcc, ",") {
			bcc = strings.TrimSpace(bcc)
			if bcc != "" {
				builder.Bcc(bcc)
			}
		}
	}
	if input != nil {
		for _, bcc := range input.Bcc {
			builder.Bcc(bcc)
		}
	}

	// 回复地址
	replyTo := ""
	if input != nil && input.ReplyTo != "" {
		replyTo = input.ReplyTo
	} else if template.ReplyTo != "" {
		replyTo = template.ReplyTo
	}
	if replyTo != "" {
		builder.ReplyTo(replyTo)
	}

	// 附件
	if input != nil {
		for _, att := range input.Attachments {
			builder.AttachWithType(att.Filename, att.Content, att.ContentType)
		}
	}

	// 发送
	_, sendErr := builder.Send(ctx)

	// 更新日志
	if sendErr != nil {
		sendLog.MarkFailed(sendErr.Error())
	} else {
		sendLog.MarkSent()
	}
	s.logRepo.Update(ctx, sendLog)

	if sendErr != nil {
		return ErrSendFailed.Wrap(sendErr)
	}

	return nil
}

// mergeParams 合并参数
func (s *Service) mergeParams(params map[string]any) map[string]any {
	merged := make(map[string]any)

	// 1. 注入通用参数（应用级）
	for k, v := range s.commonParams {
		merged[k] = v
	}

	// 2. 自动注入 CurrentYear（如未在 commonParams 中定义）
	if _, ok := merged["CurrentYear"]; !ok {
		merged["CurrentYear"] = time.Now().Year()
	}

	// 3. 用户传入的参数（优先级最高）
	for k, v := range params {
		merged[k] = v
	}

	return merged
}

// ========== 日志查询 ==========

// GetSendLogs 获取发送日志
func (s *Service) GetSendLogs(ctx context.Context, filter LogFilter) (*PageResult[model.SendLog], error) {
	return s.logRepo.List(ctx, filter)
}

// GetSendLog 获取日志详情
func (s *Service) GetSendLog(ctx context.Context, id uint) (*model.SendLog, error) {
	return s.logRepo.GetByID(ctx, id)
}
