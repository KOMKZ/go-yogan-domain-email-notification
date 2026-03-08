package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
)

type TemplateFilter struct {
	TriggerCode string
	Language    string
	Status      model.TemplateStatus
	Page        int
	PageSize    int
}

type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type TemplateRepository interface {
	Create(ctx context.Context, template *model.Template) error
	Update(ctx context.Context, template *model.Template) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Template, error)
	GetActiveTemplate(ctx context.Context, triggerCode, language string) (*model.Template, error)
	List(ctx context.Context, filter TemplateFilter) (*PageResult[model.Template], error)
	ExistsByTriggerAndLanguage(ctx context.Context, triggerCode, language string, excludeID uint) (bool, error)
}
