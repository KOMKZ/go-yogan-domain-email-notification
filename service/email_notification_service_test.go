package service

import (
	"context"
	"errors"
	"testing"
	"time"

	emailnotification "github.com/KOMKZ/go-yogan-domain-email-notification"
	domainerrors "github.com/KOMKZ/go-yogan-domain-email-notification/errors"
	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"github.com/KOMKZ/go-yogan-domain-email-notification/repository"
	"github.com/KOMKZ/go-yogan-framework/logger"
)

// ==================== Mocks ====================

type mockTemplateRepo struct {
	templates map[uint]*model.Template
	nextID    uint
}

func newMockTemplateRepo() *mockTemplateRepo {
	return &mockTemplateRepo{templates: make(map[uint]*model.Template), nextID: 1}
}

func (m *mockTemplateRepo) Create(_ context.Context, t *model.Template) error {
	t.ID = m.nextID
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	m.nextID++
	m.templates[t.ID] = t
	return nil
}

func (m *mockTemplateRepo) Update(_ context.Context, t *model.Template) error {
	if _, ok := m.templates[t.ID]; !ok {
		return domainerrors.ErrTemplateNotFound
	}
	t.UpdatedAt = time.Now()
	m.templates[t.ID] = t
	return nil
}

func (m *mockTemplateRepo) Delete(_ context.Context, id uint) error {
	delete(m.templates, id)
	return nil
}

func (m *mockTemplateRepo) GetByID(_ context.Context, id uint) (*model.Template, error) {
	t, ok := m.templates[id]
	if !ok {
		return nil, domainerrors.ErrTemplateNotFound
	}
	return t, nil
}

func (m *mockTemplateRepo) GetActiveTemplate(_ context.Context, triggerCode, language string) (*model.Template, error) {
	for _, t := range m.templates {
		if t.TriggerCode == triggerCode && t.Language == language && t.Status == model.TemplateStatusEnabled {
			return t, nil
		}
	}
	return nil, domainerrors.ErrTemplateNotFound
}

func (m *mockTemplateRepo) List(_ context.Context, filter repository.TemplateFilter) (*repository.PageResult[model.Template], error) {
	var items []model.Template
	for _, t := range m.templates {
		if filter.TriggerCode != "" && t.TriggerCode != filter.TriggerCode {
			continue
		}
		if filter.Language != "" && t.Language != filter.Language {
			continue
		}
		if filter.Status != "" && t.Status != filter.Status {
			continue
		}
		items = append(items, *t)
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &repository.PageResult[model.Template]{
		Items:      items,
		Total:      int64(len(items)),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: 1,
	}, nil
}

func (m *mockTemplateRepo) ExistsByTriggerAndLanguage(_ context.Context, triggerCode, language string, excludeID uint) (bool, error) {
	for _, t := range m.templates {
		if t.TriggerCode == triggerCode && t.Language == language && t.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}

type mockSendLogRepo struct {
	logs   map[uint]*model.SendLog
	nextID uint
}

func newMockSendLogRepo() *mockSendLogRepo {
	return &mockSendLogRepo{logs: make(map[uint]*model.SendLog), nextID: 1}
}

func (m *mockSendLogRepo) Create(_ context.Context, l *model.SendLog) error {
	l.ID = m.nextID
	l.CreatedAt = time.Now()
	m.nextID++
	m.logs[l.ID] = l
	return nil
}

func (m *mockSendLogRepo) Update(_ context.Context, l *model.SendLog) error {
	m.logs[l.ID] = l
	return nil
}

func (m *mockSendLogRepo) GetByID(_ context.Context, id uint) (*model.SendLog, error) {
	l, ok := m.logs[id]
	if !ok {
		return nil, domainerrors.ErrSendLogNotFound
	}
	return l, nil
}

func (m *mockSendLogRepo) List(_ context.Context, filter repository.LogFilter) (*repository.PageResult[model.SendLog], error) {
	var items []model.SendLog
	for _, l := range m.logs {
		if filter.TriggerCode != "" && l.TriggerCode != filter.TriggerCode {
			continue
		}
		if filter.Status != "" && l.Status != filter.Status {
			continue
		}
		items = append(items, *l)
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &repository.PageResult[model.SendLog]{
		Items:      items,
		Total:      int64(len(items)),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: 1,
	}, nil
}

type mockEmailSender struct {
	sent    []emailnotification.EmailSendInput
	failErr error
}

func newMockEmailSender() *mockEmailSender {
	return &mockEmailSender{}
}

func (m *mockEmailSender) Send(_ context.Context, input emailnotification.EmailSendInput) error {
	if m.failErr != nil {
		return m.failErr
	}
	m.sent = append(m.sent, input)
	return nil
}

// ==================== Helper ====================

func setupService() (*EmailNotificationService, *mockTemplateRepo, *mockSendLogRepo, *mockEmailSender) {
	templateRepo := newMockTemplateRepo()
	logRepo := newMockSendLogRepo()
	sender := newMockEmailSender()
	registry := emailnotification.NewTriggerRegistry()
	registry.SetCommonParams([]emailnotification.Param{
		{Name: "AppName", Type: "string", Example: "TestApp"},
	})
	registry.Register("user.registered", "用户注册", "新用户注册", []emailnotification.Param{
		{Name: "UserName", Type: "string", Required: true, Example: "张三"},
	})
	registry.Register("order.paid", "订单支付", "订单支付成功", []emailnotification.Param{
		{Name: "OrderNo", Type: "string", Required: true, Example: "ORD001"},
	})

	svc := NewEmailNotificationService(templateRepo, logRepo, sender, registry, map[string]any{"AppName": "TestApp"}, logger.GetLogger("email_test"))
	return svc, templateRepo, logRepo, sender
}

func createEnabledTemplate(t *testing.T, svc *EmailNotificationService) *model.Template {
	t.Helper()
	ctx := context.Background()
	tmpl, err := svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "user.registered",
		Language:    "zh-CN",
		Name:        "注册邮件",
		Subject:     "欢迎 {{.UserName}}",
		BodyHTML:    "<p>Hello {{.UserName}}</p>",
		BodyText:    "Hello {{.UserName}}",
		Status:      model.TemplateStatusEnabled,
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}
	return tmpl
}

// ==================== Tests ====================

func TestListTriggers(t *testing.T) {
	svc, _, _, _ := setupService()
	triggers := svc.ListTriggers(context.Background())
	if len(triggers) != 2 {
		t.Errorf("len = %d, want 2", len(triggers))
	}
}

func TestGetTrigger(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()

	trigger, err := svc.GetTrigger(ctx, "user.registered")
	if err != nil {
		t.Fatal(err)
	}
	if trigger.Code != "user.registered" {
		t.Errorf("code = %s", trigger.Code)
	}

	_, err = svc.GetTrigger(ctx, "nonexistent")
	if !errors.Is(err, domainerrors.ErrTriggerNotFound) {
		t.Errorf("expected ErrTriggerNotFound, got %v", err)
	}
}

func TestGetTriggerParams(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()

	params, err := svc.GetTriggerParams(ctx, "user.registered")
	if err != nil {
		t.Fatal(err)
	}
	if len(params) != 2 {
		t.Errorf("len = %d, want 2 (1 common + 1 trigger)", len(params))
	}

	_, err = svc.GetTriggerParams(ctx, "nonexistent")
	if !errors.Is(err, domainerrors.ErrTriggerNotFound) {
		t.Errorf("expected ErrTriggerNotFound")
	}
}

func TestGetRegistry(t *testing.T) {
	svc, _, _, _ := setupService()
	reg := svc.GetRegistry()
	if reg == nil {
		t.Fatal("registry is nil")
	}
}

func TestCreateTemplate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()

	tmpl, err := svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "user.registered",
		Name:        "Test",
		Subject:     "Subject",
		BodyHTML:    "<p>Body</p>",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if tmpl.Language != "zh-CN" {
		t.Errorf("language = %s, want zh-CN (default)", tmpl.Language)
	}
	if tmpl.Status != model.TemplateStatusDraft {
		t.Errorf("status = %s, want draft (default)", tmpl.Status)
	}
}

func TestCreateTemplate_TriggerNotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	_, err := svc.CreateTemplate(context.Background(), CreateTemplateInput{
		TriggerCode: "nonexistent",
		Name:        "Test",
		Subject:     "Sub",
		BodyHTML:    "<p></p>",
	})
	if !errors.Is(err, domainerrors.ErrTriggerNotFound) {
		t.Errorf("expected ErrTriggerNotFound, got %v", err)
	}
}

func TestCreateTemplate_Duplicate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	_, _ = svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "user.registered",
		Language:    "zh-CN",
		Name:        "First",
		Subject:     "Sub",
		BodyHTML:    "<p></p>",
	})
	_, err := svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "user.registered",
		Language:    "zh-CN",
		Name:        "Second",
		Subject:     "Sub",
		BodyHTML:    "<p></p>",
	})
	if !errors.Is(err, domainerrors.ErrTemplateExists) {
		t.Errorf("expected ErrTemplateExists, got %v", err)
	}
}

func TestUpdateTemplate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	newName := "Updated"
	newStatus := model.TemplateStatusDisabled
	updated, err := svc.UpdateTemplate(ctx, tmpl.ID, UpdateTemplateInput{
		Name:   &newName,
		Status: &newStatus,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Updated" {
		t.Errorf("name = %s", updated.Name)
	}
	if updated.Status != model.TemplateStatusDisabled {
		t.Errorf("status = %s", updated.Status)
	}
}

func TestUpdateTemplate_NotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	name := "x"
	_, err := svc.UpdateTemplate(context.Background(), 9999, UpdateTemplateInput{Name: &name})
	if !errors.Is(err, domainerrors.ErrTemplateNotFound) {
		t.Errorf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestUpdateTemplate_AllFields(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	subject := "new subject"
	bodyHTML := "<p>new</p>"
	bodyText := "new text"
	cc := "cc@test.com"
	bcc := "bcc@test.com"
	replyTo := "reply@test.com"

	updated, err := svc.UpdateTemplate(ctx, tmpl.ID, UpdateTemplateInput{
		Subject:  &subject,
		BodyHTML: &bodyHTML,
		BodyText: &bodyText,
		Cc:       &cc,
		Bcc:      &bcc,
		ReplyTo:  &replyTo,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Subject != subject || updated.BodyHTML != bodyHTML || updated.BodyText != bodyText {
		t.Error("fields not updated")
	}
	if updated.Cc != cc || updated.Bcc != bcc || updated.ReplyTo != replyTo {
		t.Error("cc/bcc/reply_to not updated")
	}
}

func TestDeleteTemplate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	if err := svc.DeleteTemplate(ctx, tmpl.ID); err != nil {
		t.Fatal(err)
	}
	_, err := svc.GetTemplate(ctx, tmpl.ID)
	if !errors.Is(err, domainerrors.ErrTemplateNotFound) {
		t.Error("expected not found after delete")
	}
}

func TestGetTemplate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	got, err := svc.GetTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != tmpl.ID {
		t.Errorf("id = %d, want %d", got.ID, tmpl.ID)
	}
}

func TestListTemplates(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	result, err := svc.ListTemplates(ctx, ListTemplatesInput{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
}

func TestListTemplates_WithFilter(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	result, err := svc.ListTemplates(ctx, ListTemplatesInput{TriggerCode: "user.registered"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d", result.Total)
	}

	result, err = svc.ListTemplates(ctx, ListTemplatesInput{TriggerCode: "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 0 {
		t.Errorf("total = %d, want 0", result.Total)
	}
}

func TestGetTemplateByTrigger(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	tmpl, err := svc.GetTemplateByTrigger(ctx, "user.registered", "zh-CN")
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.TriggerCode != "user.registered" {
		t.Errorf("trigger = %s", tmpl.TriggerCode)
	}
}

func TestGetTemplateByTrigger_NotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	_, err := svc.GetTemplateByTrigger(context.Background(), "nonexistent", "zh-CN")
	if !errors.Is(err, domainerrors.ErrTemplateNotFound) {
		t.Errorf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestPreviewTemplate(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	preview, err := svc.PreviewTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatal(err)
	}
	if preview.Subject == "" || preview.BodyHTML == "" {
		t.Error("expected non-empty preview")
	}
}

func TestPreviewTemplate_NotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	_, err := svc.PreviewTemplate(context.Background(), 9999)
	if !errors.Is(err, domainerrors.ErrTemplateNotFound) {
		t.Errorf("expected ErrTemplateNotFound")
	}
}

func TestTestSend(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	err := svc.TestSend(ctx, tmpl.ID, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(sender.sent) != 1 {
		t.Errorf("sent count = %d", len(sender.sent))
	}
}

func TestTestSend_NoRecipient(t *testing.T) {
	svc, _, _, _ := setupService()
	tmpl := createEnabledTemplate(t, svc)

	err := svc.TestSend(context.Background(), tmpl.ID, "")
	if !errors.Is(err, domainerrors.ErrNoRecipient) {
		t.Errorf("expected ErrNoRecipient, got %v", err)
	}
}

func TestTestSend_TemplateNotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	err := svc.TestSend(context.Background(), 9999, "test@example.com")
	if !errors.Is(err, domainerrors.ErrTemplateNotFound) {
		t.Errorf("expected ErrTemplateNotFound")
	}
}

func TestSend(t *testing.T) {
	svc, _, logRepo, sender := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	err := svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "user@example.com",
		Params:      map[string]any{"UserName": "TestUser"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sender.sent) != 1 {
		t.Errorf("sent count = %d", len(sender.sent))
	}
	if len(logRepo.logs) != 1 {
		t.Errorf("log count = %d", len(logRepo.logs))
	}
	for _, l := range logRepo.logs {
		if l.Status != model.SendStatusSent {
			t.Errorf("log status = %s, want sent", l.Status)
		}
	}
}

func TestSend_NoTriggerCode(t *testing.T) {
	svc, _, _, _ := setupService()
	err := svc.Send(context.Background(), SendInput{Recipient: "a@b.com"})
	if !errors.Is(err, domainerrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestSend_NoRecipient(t *testing.T) {
	svc, _, _, _ := setupService()
	err := svc.Send(context.Background(), SendInput{TriggerCode: "user.registered"})
	if !errors.Is(err, domainerrors.ErrNoRecipient) {
		t.Errorf("expected ErrNoRecipient, got %v", err)
	}
}

func TestSend_TriggerNotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	err := svc.Send(context.Background(), SendInput{TriggerCode: "nope", Recipient: "a@b.com"})
	if !errors.Is(err, domainerrors.ErrTriggerNotFound) {
		t.Errorf("expected ErrTriggerNotFound, got %v", err)
	}
}

func TestSend_LanguageFallback(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc) // zh-CN

	err := svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Language:    "en-US",
		Params:      map[string]any{"UserName": "Test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sender.sent) != 1 {
		t.Error("expected fallback to zh-CN")
	}
}

func TestSend_Failure(t *testing.T) {
	svc, _, logRepo, sender := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	sender.failErr = errors.New("smtp timeout")

	err := svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Params:      map[string]any{"UserName": "Test"},
	})
	if !errors.Is(err, domainerrors.ErrSendFailed) {
		t.Errorf("expected ErrSendFailed, got %v", err)
	}
	for _, l := range logRepo.logs {
		if l.Status != model.SendStatusFailed {
			t.Errorf("log status = %s, want failed", l.Status)
		}
	}
}

func TestSend_WithCcBccReplyTo(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()

	svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "order.paid",
		Language:    "zh-CN",
		Name:        "订单邮件",
		Subject:     "Order {{.OrderNo}}",
		BodyHTML:    "<p>{{.OrderNo}}</p>",
		Status:      model.TemplateStatusEnabled,
		Cc:          "cc1@test.com, cc2@test.com",
		Bcc:         "bcc1@test.com",
		ReplyTo:     "reply@test.com",
	})

	err := svc.Send(ctx, SendInput{
		TriggerCode: "order.paid",
		Recipient:   "buyer@test.com",
		Params:      map[string]any{"OrderNo": "ORD123"},
		Cc:          []string{"extra_cc@test.com"},
		Bcc:         []string{"extra_bcc@test.com"},
		ReplyTo:     "override@test.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sender.sent) != 1 {
		t.Fatal("expected 1 sent")
	}
	sent := sender.sent[0]
	if len(sent.Cc) != 3 {
		t.Errorf("cc count = %d, want 3", len(sent.Cc))
	}
	if len(sent.Bcc) != 2 {
		t.Errorf("bcc count = %d, want 2", len(sent.Bcc))
	}
	if sent.ReplyTo != "override@test.com" {
		t.Errorf("reply_to = %s", sent.ReplyTo)
	}
}

func TestSend_WithFromOverride(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	err := svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Params:      map[string]any{"UserName": "Test"},
		From:        "custom@sender.com",
		FromName:    "Custom Sender",
	})
	if err != nil {
		t.Fatal(err)
	}
	if sender.sent[0].From != "custom@sender.com" {
		t.Errorf("from = %s", sender.sent[0].From)
	}
	if sender.sent[0].FromName != "Custom Sender" {
		t.Errorf("from_name = %s", sender.sent[0].FromName)
	}
}

func TestMergeParams(t *testing.T) {
	svc, _, _, _ := setupService()
	merged := svc.mergeParams(map[string]any{"UserName": "Test", "AppName": "Override"})

	if merged["AppName"] != "Override" {
		t.Error("user params should override common")
	}
	if merged["UserName"] != "Test" {
		t.Error("user params missing")
	}
	if _, ok := merged["CurrentYear"]; !ok {
		t.Error("CurrentYear should be injected")
	}
}

func TestGetSendLogs(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	_ = svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Params:      map[string]any{"UserName": "Test"},
	})

	logs, err := svc.GetSendLogs(ctx, ListSendLogsInput{})
	if err != nil {
		t.Fatal(err)
	}
	if logs.Total != 1 {
		t.Errorf("total = %d", logs.Total)
	}
}

func TestGetSendLog(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	_ = svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Params:      map[string]any{"UserName": "Test"},
	})

	log, err := svc.GetSendLog(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if log.Recipient != "a@b.com" {
		t.Errorf("recipient = %s", log.Recipient)
	}
}

func TestGetSendLog_NotFound(t *testing.T) {
	svc, _, _, _ := setupService()
	_, err := svc.GetSendLog(context.Background(), 9999)
	if !errors.Is(err, domainerrors.ErrSendLogNotFound) {
		t.Errorf("expected ErrSendLogNotFound, got %v", err)
	}
}

func TestCreateTemplate_WithExplicitStatus(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl, err := svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "user.registered",
		Language:    "en-US",
		Name:        "English",
		Subject:     "Welcome",
		BodyHTML:    "<p>Hi</p>",
		Status:      model.TemplateStatusEnabled,
	})
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.Status != model.TemplateStatusEnabled {
		t.Errorf("status = %s, want enabled", tmpl.Status)
	}
	if tmpl.Language != "en-US" {
		t.Errorf("language = %s, want en-US", tmpl.Language)
	}
}

func TestPreviewTemplate_WithBodyText(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl := createEnabledTemplate(t, svc)

	preview, err := svc.PreviewTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatal(err)
	}
	if preview.BodyText == "" {
		t.Error("expected non-empty body_text preview")
	}
}

func TestPreviewTemplate_NoBodyText(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	tmpl, _ := svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "order.paid",
		Language:    "zh-CN",
		Name:        "No Text",
		Subject:     "Order {{.OrderNo}}",
		BodyHTML:    "<p>{{.OrderNo}}</p>",
		BodyText:    "",
		Status:      model.TemplateStatusEnabled,
	})
	preview, err := svc.PreviewTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatal(err)
	}
	if preview.BodyText != "" {
		t.Errorf("expected empty body_text, got %q", preview.BodyText)
	}
}

func TestSend_DefaultLanguage(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	err := svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Language:    "",
		Params:      map[string]any{"UserName": "Test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sender.sent) != 1 {
		t.Error("should have sent with default language")
	}
}

func TestSend_TemplateReplyToFallback(t *testing.T) {
	svc, _, _, sender := setupService()
	ctx := context.Background()
	svc.CreateTemplate(ctx, CreateTemplateInput{
		TriggerCode: "order.paid",
		Language:    "zh-CN",
		Name:        "With ReplyTo",
		Subject:     "Order {{.OrderNo}}",
		BodyHTML:    "<p>{{.OrderNo}}</p>",
		Status:      model.TemplateStatusEnabled,
		ReplyTo:     "tmpl-reply@test.com",
	})

	err := svc.Send(ctx, SendInput{
		TriggerCode: "order.paid",
		Recipient:   "a@b.com",
		Params:      map[string]any{"OrderNo": "ORD1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if sender.sent[0].ReplyTo != "tmpl-reply@test.com" {
		t.Errorf("reply_to = %s, want tmpl-reply@test.com", sender.sent[0].ReplyTo)
	}
}

func TestGetSendLogs_WithFilter(t *testing.T) {
	svc, _, _, _ := setupService()
	ctx := context.Background()
	createEnabledTemplate(t, svc)

	_ = svc.Send(ctx, SendInput{
		TriggerCode: "user.registered",
		Recipient:   "a@b.com",
		Params:      map[string]any{"UserName": "Test"},
	})

	logs, err := svc.GetSendLogs(ctx, ListSendLogsInput{TriggerCode: "user.registered"})
	if err != nil {
		t.Fatal(err)
	}
	if logs.Total != 1 {
		t.Errorf("total = %d", logs.Total)
	}

	logs, err = svc.GetSendLogs(ctx, ListSendLogsInput{TriggerCode: "nonexistent"})
	if err != nil {
		t.Fatal(err)
	}
	if logs.Total != 0 {
		t.Errorf("total = %d, want 0", logs.Total)
	}
}
