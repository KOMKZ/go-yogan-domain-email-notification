package email_notification

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-email-notification/model"
)

// TemplateFilter 模板筛选条件
type TemplateFilter struct {
	TriggerCode string
	Language    string
	Status      model.TemplateStatus
	Page        int
	PageSize    int
}

// LogFilter 日志筛选条件
type LogFilter struct {
	TriggerCode string
	Status      model.SendStatus
	StartTime   string
	EndTime     string
	Page        int
	PageSize    int
}

// PageResult 分页结果
type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// TemplateRepository 模板仓储接口
type TemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *model.Template) error

	// Update 更新模板
	Update(ctx context.Context, template *model.Template) error

	// Delete 删除模板（软删除）
	Delete(ctx context.Context, id uint) error

	// GetByID 根据 ID 获取模板
	GetByID(ctx context.Context, id uint) (*model.Template, error)

	// GetActiveTemplate 获取指定触发点和语言的启用模板
	GetActiveTemplate(ctx context.Context, triggerCode, language string) (*model.Template, error)

	// List 列表查询
	List(ctx context.Context, filter TemplateFilter) (*PageResult[model.Template], error)

	// ExistsByTriggerAndLanguage 检查是否存在相同触发点和语言的模板
	ExistsByTriggerAndLanguage(ctx context.Context, triggerCode, language string, excludeID uint) (bool, error)
}

// SendLogRepository 发送日志仓储接口
type SendLogRepository interface {
	// Create 创建日志
	Create(ctx context.Context, log *model.SendLog) error

	// Update 更新日志
	Update(ctx context.Context, log *model.SendLog) error

	// GetByID 根据 ID 获取日志
	GetByID(ctx context.Context, id uint) (*model.SendLog, error)

	// List 列表查询
	List(ctx context.Context, filter LogFilter) (*PageResult[model.SendLog], error)
}
