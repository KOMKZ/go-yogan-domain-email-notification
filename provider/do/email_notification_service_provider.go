package do

import (
	emailnotification "github.com/KOMKZ/go-yogan-domain-email-notification"
	"github.com/KOMKZ/go-yogan-domain-email-notification/repository"
	"github.com/KOMKZ/go-yogan-domain-email-notification/sender"
	"github.com/KOMKZ/go-yogan-domain-email-notification/service"
	email "github.com/KOMKZ/go-yogan-component-email"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

// ---- Repository Providers ----

func ProvideTemplateRepository(i do.Injector) (repository.TemplateRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewTemplateMySQLRepository(db), nil
}

func ProvideSendLogRepository(i do.Injector) (repository.SendLogRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewSendLogMySQLRepository(db), nil
}

// ---- Sender Providers ----

func ProvideEmailSender(i do.Injector) (emailnotification.EmailSender, error) {
	manager, err := do.Invoke[*email.Manager](i)
	if err != nil {
		return nil, err
	}
	return sender.NewComponentEmailSender(manager), nil
}

// ---- Service Providers ----

func ProvideEmailNotificationService(i do.Injector) (*service.EmailNotificationService, error) {
	templateRepo, err := do.Invoke[repository.TemplateRepository](i)
	if err != nil {
		return nil, err
	}
	logRepo, err := do.Invoke[repository.SendLogRepository](i)
	if err != nil {
		return nil, err
	}
	emailSender, err := do.Invoke[emailnotification.EmailSender](i)
	if err != nil {
		return nil, err
	}
	registry, err := do.Invoke[*emailnotification.TriggerRegistry](i)
	if err != nil {
		return nil, err
	}
	commonParams, _ := do.InvokeNamed[map[string]any](i, "email_common_params")
	log := do.MustInvokeNamed[*logger.CtxZapLogger](i, "email")
	return service.NewEmailNotificationService(templateRepo, logRepo, emailSender, registry, commonParams, log), nil
}
