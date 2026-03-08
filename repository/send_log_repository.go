package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
)

type LogFilter struct {
	TriggerCode string
	Status      model.SendStatus
	StartTime   string
	EndTime     string
	Page        int
	PageSize    int
}

type SendLogRepository interface {
	Create(ctx context.Context, log *model.SendLog) error
	Update(ctx context.Context, log *model.SendLog) error
	GetByID(ctx context.Context, id uint) (*model.SendLog, error)
	List(ctx context.Context, filter LogFilter) (*PageResult[model.SendLog], error)
}
