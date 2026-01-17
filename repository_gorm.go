package email_notification

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
	"gorm.io/gorm"
)

// ============ Template Repository GORM 实现 ============

type gormTemplateRepository struct {
	db *gorm.DB
}

// NewGormTemplateRepository 创建 GORM 模板仓储
func NewGormTemplateRepository(db *gorm.DB) TemplateRepository {
	return &gormTemplateRepository{db: db}
}

func (r *gormTemplateRepository) Create(ctx context.Context, template *model.Template) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *gormTemplateRepository) Update(ctx context.Context, template *model.Template) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *gormTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Template{}, id).Error
}

func (r *gormTemplateRepository) GetByID(ctx context.Context, id uint) (*model.Template, error) {
	var template model.Template
	err := r.db.WithContext(ctx).First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, ErrDatabaseError.Wrap(err)
	}
	return &template, nil
}

func (r *gormTemplateRepository) GetActiveTemplate(ctx context.Context, triggerCode, language string) (*model.Template, error) {
	var template model.Template
	err := r.db.WithContext(ctx).
		Where("trigger_code = ? AND language = ? AND status = ?", triggerCode, language, model.TemplateStatusEnabled).
		First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, ErrDatabaseError.Wrap(err)
	}
	return &template, nil
}

func (r *gormTemplateRepository) List(ctx context.Context, filter TemplateFilter) (*PageResult[model.Template], error) {
	query := r.db.WithContext(ctx).Model(&model.Template{})

	if filter.TriggerCode != "" {
		query = query.Where("trigger_code = ?", filter.TriggerCode)
	}
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []model.Template
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult[model.Template]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *gormTemplateRepository) ExistsByTriggerAndLanguage(ctx context.Context, triggerCode, language string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Template{}).
		Where("trigger_code = ? AND language = ?", triggerCode, language)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, ErrDatabaseError.Wrap(err)
	}
	return count > 0, nil
}

// ============ SendLog Repository GORM 实现 ============

type gormSendLogRepository struct {
	db *gorm.DB
}

// NewGormSendLogRepository 创建 GORM 发送日志仓储
func NewGormSendLogRepository(db *gorm.DB) SendLogRepository {
	return &gormSendLogRepository{db: db}
}

func (r *gormSendLogRepository) Create(ctx context.Context, log *model.SendLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *gormSendLogRepository) Update(ctx context.Context, log *model.SendLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *gormSendLogRepository) GetByID(ctx context.Context, id uint) (*model.SendLog, error) {
	var log model.SendLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrSendLogNotFound
		}
		return nil, ErrDatabaseError.Wrap(err)
	}
	return &log, nil
}

func (r *gormSendLogRepository) List(ctx context.Context, filter LogFilter) (*PageResult[model.SendLog], error) {
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
		return nil, ErrDatabaseError.Wrap(err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []model.SendLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult[model.SendLog]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
