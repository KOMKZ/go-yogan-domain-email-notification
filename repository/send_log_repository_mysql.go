package repository

import (
	"context"
	"errors"

	domainerrors "github.com/KOMKZ/go-yogan-domain-email-notification/errors"
	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"github.com/KOMKZ/go-yogan-framework/database"
	"gorm.io/gorm"
)

type SendLogMySQLRepository struct {
	base *database.BaseRepository[model.SendLog]
	db   *gorm.DB
}

func NewSendLogMySQLRepository(db *gorm.DB) *SendLogMySQLRepository {
	return &SendLogMySQLRepository{
		base: database.NewBaseRepository[model.SendLog](db),
		db:   db,
	}
}

func (r *SendLogMySQLRepository) Create(ctx context.Context, l *model.SendLog) error {
	return r.base.Create(ctx, l)
}

func (r *SendLogMySQLRepository) Update(ctx context.Context, l *model.SendLog) error {
	return r.base.Update(ctx, l)
}

func (r *SendLogMySQLRepository) GetByID(ctx context.Context, id uint) (*model.SendLog, error) {
	var l model.SendLog
	err := r.db.WithContext(ctx).First(&l, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrSendLogNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *SendLogMySQLRepository) List(ctx context.Context, filter LogFilter) (*PageResult[model.SendLog], error) {
	query := r.db.WithContext(ctx).Model(&model.SendLog{})
	if filter.TriggerCode != "" {
		query = query.Where("trigger_code = ?", filter.TriggerCode)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartTime != "" {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != "" {
		query = query.Where("created_at <= ?", filter.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page, pageSize := filter.Page, filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var items []model.SendLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult[model.SendLog]{
		Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages,
	}, nil
}
