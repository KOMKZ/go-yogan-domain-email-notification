package do

import (
	"context"
	"testing"

	emailnotification "github.com/KOMKZ/go-yogan-domain-email-notification"
	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"github.com/KOMKZ/go-yogan-domain-email-notification/repository"
	"github.com/KOMKZ/go-yogan-domain-email-notification/service"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"github.com/samber/do/v2"
)

type stubTemplateRepo struct{}

func (s *stubTemplateRepo) Create(_ context.Context, _ *model.Template) error          { return nil }
func (s *stubTemplateRepo) Update(_ context.Context, _ *model.Template) error          { return nil }
func (s *stubTemplateRepo) Delete(_ context.Context, _ uint) error                     { return nil }
func (s *stubTemplateRepo) GetByID(_ context.Context, _ uint) (*model.Template, error) { return nil, nil }
func (s *stubTemplateRepo) GetActiveTemplate(_ context.Context, _, _ string) (*model.Template, error) {
	return nil, nil
}
func (s *stubTemplateRepo) List(_ context.Context, _ repository.TemplateFilter) (*repository.PageResult[model.Template], error) {
	return nil, nil
}
func (s *stubTemplateRepo) ExistsByTriggerAndLanguage(_ context.Context, _, _ string, _ uint) (bool, error) {
	return false, nil
}

type stubSendLogRepo struct{}

func (s *stubSendLogRepo) Create(_ context.Context, _ *model.SendLog) error          { return nil }
func (s *stubSendLogRepo) Update(_ context.Context, _ *model.SendLog) error          { return nil }
func (s *stubSendLogRepo) GetByID(_ context.Context, _ uint) (*model.SendLog, error) { return nil, nil }
func (s *stubSendLogRepo) List(_ context.Context, _ repository.LogFilter) (*repository.PageResult[model.SendLog], error) {
	return nil, nil
}

type stubEmailSender struct{}

func (s *stubEmailSender) Send(_ context.Context, _ emailnotification.EmailSendInput) error {
	return nil
}

func TestProvideEmailNotificationService(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.TemplateRepository](injector, &stubTemplateRepo{})
	do.ProvideValue[repository.SendLogRepository](injector, &stubSendLogRepo{})
	do.ProvideValue[emailnotification.EmailSender](injector, &stubEmailSender{})
	do.ProvideValue(injector, emailnotification.NewTriggerRegistry())
	do.ProvideNamedValue(injector, "email_common_params", map[string]any{"AppName": "Test"})

	do.Provide(injector, ProvideEmailNotificationService)

	svc, err := do.Invoke[*service.EmailNotificationService](injector)
	if err != nil {
		t.Fatalf("invoke service: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
	}
}

func TestProvideEmailNotificationService_WithoutCommonParams(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.TemplateRepository](injector, &stubTemplateRepo{})
	do.ProvideValue[repository.SendLogRepository](injector, &stubSendLogRepo{})
	do.ProvideValue[emailnotification.EmailSender](injector, &stubEmailSender{})
	do.ProvideValue(injector, emailnotification.NewTriggerRegistry())

	do.Provide(injector, ProvideEmailNotificationService)

	svc, err := do.Invoke[*service.EmailNotificationService](injector)
	if err != nil {
		t.Fatalf("invoke service: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
	}
}

func TestProvideEmailNotificationService_MissingTemplateRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.SendLogRepository](injector, &stubSendLogRepo{})
	do.ProvideValue[emailnotification.EmailSender](injector, &stubEmailSender{})
	do.ProvideValue(injector, emailnotification.NewTriggerRegistry())

	do.Provide(injector, ProvideEmailNotificationService)

	_, err := do.Invoke[*service.EmailNotificationService](injector)
	if err == nil {
		t.Fatal("expected error when template repo missing")
	}
}

func TestProvideEmailNotificationService_MissingLogRepo(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.TemplateRepository](injector, &stubTemplateRepo{})
	do.ProvideValue[emailnotification.EmailSender](injector, &stubEmailSender{})
	do.ProvideValue(injector, emailnotification.NewTriggerRegistry())

	do.Provide(injector, ProvideEmailNotificationService)

	_, err := do.Invoke[*service.EmailNotificationService](injector)
	if err == nil {
		t.Fatal("expected error when log repo missing")
	}
}

func TestProvideEmailNotificationService_MissingEmailSender(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.TemplateRepository](injector, &stubTemplateRepo{})
	do.ProvideValue[repository.SendLogRepository](injector, &stubSendLogRepo{})
	do.ProvideValue(injector, emailnotification.NewTriggerRegistry())

	do.Provide(injector, ProvideEmailNotificationService)

	_, err := do.Invoke[*service.EmailNotificationService](injector)
	if err == nil {
		t.Fatal("expected error when email sender missing")
	}
}

func TestProvideEmailNotificationService_MissingRegistry(t *testing.T) {
	injector := do.New()
	do.ProvideNamedValue(injector, "email", logger.GetLogger("email_test"))
	do.ProvideValue[repository.TemplateRepository](injector, &stubTemplateRepo{})
	do.ProvideValue[repository.SendLogRepository](injector, &stubSendLogRepo{})
	do.ProvideValue[emailnotification.EmailSender](injector, &stubEmailSender{})

	do.Provide(injector, ProvideEmailNotificationService)

	_, err := do.Invoke[*service.EmailNotificationService](injector)
	if err == nil {
		t.Fatal("expected error when registry missing")
	}
}
