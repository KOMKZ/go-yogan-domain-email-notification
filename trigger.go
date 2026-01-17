package email_notification

import "sync"

// Param 参数定义
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`        // string, number, url, datetime, array
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example"` // 示例值（用于预览）
}

// TriggerDefinition 触发点定义
type TriggerDefinition struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
}

// TriggerRegistry 触发点注册表
type TriggerRegistry struct {
	mu           sync.RWMutex
	triggers     map[string]*TriggerDefinition
	commonParams []Param // 通用参数（应用层注入）
}

// NewTriggerRegistry 创建注册表（应用层自己保存）
func NewTriggerRegistry() *TriggerRegistry {
	return &TriggerRegistry{
		triggers:     make(map[string]*TriggerDefinition),
		commonParams: nil,
	}
}

// Register 注册触发点
func (r *TriggerRegistry) Register(code, name, description string, params []Param) *TriggerDefinition {
	def := &TriggerDefinition{
		Code:        code,
		Name:        name,
		Description: description,
		Params:      params,
	}
	r.mu.Lock()
	r.triggers[code] = def
	r.mu.Unlock()
	return def
}

// Get 获取触发点
func (r *TriggerRegistry) Get(code string) (*TriggerDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.triggers[code]
	return t, ok
}

// GetAll 获取所有触发点
func (r *TriggerRegistry) GetAll() []*TriggerDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*TriggerDefinition, 0, len(r.triggers))
	for _, t := range r.triggers {
		result = append(result, t)
	}
	return result
}

// Exists 检查触发点是否存在
func (r *TriggerRegistry) Exists(code string) bool {
	_, ok := r.Get(code)
	return ok
}

// Codes 获取所有触发点代码
func (r *TriggerRegistry) Codes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	codes := make([]string, 0, len(r.triggers))
	for code := range r.triggers {
		codes = append(codes, code)
	}
	return codes
}

// SetCommonParams 设置通用参数（应用层注入）
func (r *TriggerRegistry) SetCommonParams(params []Param) {
	r.mu.Lock()
	r.commonParams = params
	r.mu.Unlock()
}

// GetCommonParams 获取通用参数
func (r *TriggerRegistry) GetCommonParams() []Param {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.commonParams
}

// GetAllParams 获取触发点的全部参数（通用 + 专属）
func (r *TriggerRegistry) GetAllParams(triggerCode string) []Param {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trigger, ok := r.triggers[triggerCode]
	if !ok {
		return r.commonParams
	}

	result := make([]Param, 0, len(r.commonParams)+len(trigger.Params))
	result = append(result, r.commonParams...)
	result = append(result, trigger.Params...)
	return result
}
