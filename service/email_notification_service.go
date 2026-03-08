package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	emailnotification "github.com/KOMKZ/go-yogan-domain-email-notification"
	domainerrors "github.com/KOMKZ/go-yogan-domain-email-notification/errors"
	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"github.com/KOMKZ/go-yogan-domain-email-notification/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

type CreateTemplateInput struct {
	TriggerCode string
	Language    string
	Name        string
	Subject     string
	BodyHTML    string
	BodyText    string
	Status      model.TemplateStatus
	Cc          string
	Bcc         string
	ReplyTo     string
}

type UpdateTemplateInput struct {
	Name     *string
	Subject  *string
	BodyHTML *string
	BodyText *string
	Status   *model.TemplateStatus
	Cc       *string
	Bcc      *string
	ReplyTo  *string
}

type SendInput struct {
	TriggerCode string
	Recipient   string
	Language    string
	Params      map[string]any
	Cc          []string
	Bcc         []string
	ReplyTo     string
	From        string
	FromName    string
	Subject     string
}

type PreviewResult struct {
	Subject  string `json:"subject"`
	BodyHTML string `json:"body_html"`
	BodyText string `json:"body_text"`
}

type EmailNotificationService struct {
	templateRepo repository.TemplateRepository
	logRepo      repository.SendLogRepository
	emailSender  emailnotification.EmailSender
	registry     *emailnotification.TriggerRegistry
	engine       *emailnotification.TemplateEngine
	commonParams map[string]any
	logger       *logger.CtxZapLogger
}

func NewEmailNotificationService(
	templateRepo repository.TemplateRepository,
	logRepo repository.SendLogRepository,
	emailSender emailnotification.EmailSender,
	registry *emailnotification.TriggerRegistry,
	commonParams map[string]any,
	log *logger.CtxZapLogger,
) *EmailNotificationService {
	return &EmailNotificationService{
		templateRepo: templateRepo,
		logRepo:      logRepo,
		emailSender:  emailSender,
		registry:     registry,
		engine:       emailnotification.NewTemplateEngine(),
		commonParams: commonParams,
		logger:       log,
	}
}

func (s *EmailNotificationService) ListTriggers(ctx context.Context) []*emailnotification.TriggerDefinition {
	return s.registry.GetAll()
}

func (s *EmailNotificationService) GetTrigger(ctx context.Context, code string) (*emailnotification.TriggerDefinition, error) {
	trigger, ok := s.registry.Get(code)
	if !ok {
		return nil, domainerrors.ErrTriggerNotFound
	}
	return trigger, nil
}

func (s *EmailNotificationService) GetTriggerParams(ctx context.Context, code string) ([]emailnotification.Param, error) {
	if !s.registry.Exists(code) {
		return nil, domainerrors.ErrTriggerNotFound
	}
	return s.registry.GetAllParams(code), nil
}

func (s *EmailNotificationService) GetRegistry() *emailnotification.TriggerRegistry {
	return s.registry
}

func (s *EmailNotificationService) CreateTemplate(ctx context.Context, input CreateTemplateInput) (*model.Template, error) {
	if !s.registry.Exists(input.TriggerCode) {
		return nil, fmt.Errorf("%w: %s", domainerrors.ErrTriggerNotFound, input.TriggerCode)
	}

	if input.Language == "" {
		input.Language = "zh-CN"
	}

	exists, err := s.templateRepo.ExistsByTriggerAndLanguage(ctx, input.TriggerCode, input.Language, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrTemplateExists
	}

	if input.Status == "" {
		input.Status = model.TemplateStatusDraft
	}

	tmpl := &model.Template{
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

	if err := s.templateRepo.Create(ctx, tmpl); err != nil {
		s.logger.ErrorCtx(ctx, "create email template failed", zap.String("trigger", input.TriggerCode), zap.Error(err))
		return nil, err
	}

	s.logger.InfoCtx(ctx, "email template created", zap.Uint("template_id", tmpl.ID), zap.String("trigger", tmpl.TriggerCode))
	return tmpl, nil
}

func (s *EmailNotificationService) UpdateTemplate(ctx context.Context, id uint, input UpdateTemplateInput) (*model.Template, error) {
	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		tmpl.Name = *input.Name
	}
	if input.Subject != nil {
		tmpl.Subject = *input.Subject
	}
	if input.BodyHTML != nil {
		tmpl.BodyHTML = *input.BodyHTML
	}
	if input.BodyText != nil {
		tmpl.BodyText = *input.BodyText
	}
	if input.Status != nil {
		tmpl.Status = *input.Status
	}
	if input.Cc != nil {
		tmpl.Cc = *input.Cc
	}
	if input.Bcc != nil {
		tmpl.Bcc = *input.Bcc
	}
	if input.ReplyTo != nil {
		tmpl.ReplyTo = *input.ReplyTo
	}

	if err := s.templateRepo.Update(ctx, tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *EmailNotificationService) DeleteTemplate(ctx context.Context, id uint) error {
	return s.templateRepo.Delete(ctx, id)
}

func (s *EmailNotificationService) GetTemplate(ctx context.Context, id uint) (*model.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

type ListTemplatesInput struct {
	TriggerCode string
	Language    string
	Status      model.TemplateStatus
	Page        int
	PageSize    int
}

func (s *EmailNotificationService) ListTemplates(ctx context.Context, input ListTemplatesInput) (*repository.PageResult[model.Template], error) {
	filter := repository.TemplateFilter{
		TriggerCode: input.TriggerCode,
		Language:    input.Language,
		Status:      input.Status,
		Page:        input.Page,
		PageSize:    input.PageSize,
	}
	return s.templateRepo.List(ctx, filter)
}

func (s *EmailNotificationService) GetTemplateByTrigger(ctx context.Context, triggerCode, language string) (*model.Template, error) {
	return s.templateRepo.GetActiveTemplate(ctx, triggerCode, language)
}

func (s *EmailNotificationService) PreviewTemplate(ctx context.Context, id uint) (*PreviewResult, error) {
	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	params := s.registry.GetAllParams(tmpl.TriggerCode)

	subject, err := s.engine.Preview(tmpl.Subject, params)
	if err != nil {
		return nil, err
	}

	bodyHTML, err := s.engine.Preview(tmpl.BodyHTML, params)
	if err != nil {
		return nil, err
	}

	bodyText := ""
	if tmpl.BodyText != "" {
		bodyText, err = s.engine.Preview(tmpl.BodyText, params)
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

func (s *EmailNotificationService) TestSend(ctx context.Context, id uint, recipient string) error {
	if recipient == "" {
		return domainerrors.ErrNoRecipient
	}

	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	params := s.registry.GetAllParams(tmpl.TriggerCode)
	exampleParams := make(map[string]any)
	for _, p := range params {
		if p.Example != "" {
			exampleParams[p.Name] = p.Example
		} else {
			exampleParams[p.Name] = "{{." + p.Name + "}}"
		}
	}
	exampleParams["CurrentYear"] = time.Now().Year()

	return s.sendWithTemplate(ctx, tmpl, recipient, exampleParams, nil)
}

func (s *EmailNotificationService) Send(ctx context.Context, input SendInput) error {
	if input.TriggerCode == "" {
		return fmt.Errorf("%w: trigger_code is required", domainerrors.ErrInvalidInput)
	}
	if input.Recipient == "" {
		return domainerrors.ErrNoRecipient
	}
	if !s.registry.Exists(input.TriggerCode) {
		return fmt.Errorf("%w: %s", domainerrors.ErrTriggerNotFound, input.TriggerCode)
	}

	language := input.Language
	if language == "" {
		language = "zh-CN"
	}

	tmpl, err := s.templateRepo.GetActiveTemplate(ctx, input.TriggerCode, language)
	if err != nil {
		if language != "zh-CN" {
			tmpl, err = s.templateRepo.GetActiveTemplate(ctx, input.TriggerCode, "zh-CN")
		}
		if err != nil {
			return err
		}
	}

	params := s.mergeParams(input.Params)
	return s.sendWithTemplate(ctx, tmpl, input.Recipient, params, &input)
}

func (s *EmailNotificationService) sendWithTemplate(ctx context.Context, tmpl *model.Template, recipient string, params map[string]any, input *SendInput) error {
	subject, err := s.engine.Render(tmpl.Subject, params)
	if err != nil {
		return err
	}

	body, err := s.engine.Render(tmpl.BodyHTML, params)
	if err != nil {
		return err
	}

	paramsJSON, _ := json.Marshal(params)
	sendLog := &model.SendLog{
		TemplateID:  &tmpl.ID,
		TriggerCode: tmpl.TriggerCode,
		Language:    tmpl.Language,
		Recipient:   recipient,
		Subject:     subject,
		Params:      string(paramsJSON),
		Status:      model.SendStatusPending,
	}
	if err := s.logRepo.Create(ctx, sendLog); err != nil {
		return err
	}

	emailInput := emailnotification.EmailSendInput{
		To:       []string{recipient},
		Subject:  subject,
		HTMLBody: body,
	}

	if input != nil && input.From != "" {
		emailInput.From = input.From
	}
	if input != nil && input.FromName != "" {
		emailInput.FromName = input.FromName
	}

	var ccList []string
	if tmpl.Cc != "" {
		for _, cc := range strings.Split(tmpl.Cc, ",") {
			cc = strings.TrimSpace(cc)
			if cc != "" {
				ccList = append(ccList, cc)
			}
		}
	}
	if input != nil {
		ccList = append(ccList, input.Cc...)
	}
	emailInput.Cc = ccList

	var bccList []string
	if tmpl.Bcc != "" {
		for _, bcc := range strings.Split(tmpl.Bcc, ",") {
			bcc = strings.TrimSpace(bcc)
			if bcc != "" {
				bccList = append(bccList, bcc)
			}
		}
	}
	if input != nil {
		bccList = append(bccList, input.Bcc...)
	}
	emailInput.Bcc = bccList

	replyTo := ""
	if input != nil && input.ReplyTo != "" {
		replyTo = input.ReplyTo
	} else if tmpl.ReplyTo != "" {
		replyTo = tmpl.ReplyTo
	}
	emailInput.ReplyTo = replyTo

	sendErr := s.emailSender.Send(ctx, emailInput)
	if sendErr != nil {
		sendLog.MarkFailed(sendErr.Error())
		s.logger.ErrorCtx(ctx, "email send failed", zap.String("recipient", recipient), zap.String("trigger", tmpl.TriggerCode), zap.Error(sendErr))
	} else {
		sendLog.MarkSent()
		s.logger.InfoCtx(ctx, "email sent", zap.String("recipient", recipient), zap.String("trigger", tmpl.TriggerCode))
	}
	_ = s.logRepo.Update(ctx, sendLog)

	if sendErr != nil {
		return fmt.Errorf("%w: %v", domainerrors.ErrSendFailed, sendErr)
	}

	return nil
}

func (s *EmailNotificationService) mergeParams(params map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range s.commonParams {
		merged[k] = v
	}
	if _, ok := merged["CurrentYear"]; !ok {
		merged["CurrentYear"] = time.Now().Year()
	}
	for k, v := range params {
		merged[k] = v
	}
	return merged
}

type ListSendLogsInput struct {
	TriggerCode string
	Status      model.SendStatus
	StartTime   string
	EndTime     string
	Page        int
	PageSize    int
}

func (s *EmailNotificationService) GetSendLogs(ctx context.Context, input ListSendLogsInput) (*repository.PageResult[model.SendLog], error) {
	filter := repository.LogFilter{
		TriggerCode: input.TriggerCode,
		Status:      input.Status,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Page:        input.Page,
		PageSize:    input.PageSize,
	}
	return s.logRepo.List(ctx, filter)
}

func (s *EmailNotificationService) GetSendLog(ctx context.Context, id uint) (*model.SendLog, error) {
	return s.logRepo.GetByID(ctx, id)
}
