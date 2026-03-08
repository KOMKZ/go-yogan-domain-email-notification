package emailnotification

import "sync"

type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example"`
}

type TriggerDefinition struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
}

type TriggerRegistry struct {
	mu           sync.RWMutex
	triggers     map[string]*TriggerDefinition
	commonParams []Param
}

func NewTriggerRegistry() *TriggerRegistry {
	return &TriggerRegistry{
		triggers: make(map[string]*TriggerDefinition),
	}
}

func (r *TriggerRegistry) Register(code, name, description string, params []Param) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.triggers[code] = &TriggerDefinition{
		Code:        code,
		Name:        name,
		Description: description,
		Params:      params,
	}
}

func (r *TriggerRegistry) Get(code string) (*TriggerDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.triggers[code]
	return t, ok
}

func (r *TriggerRegistry) GetAll() []*TriggerDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*TriggerDefinition, 0, len(r.triggers))
	for _, t := range r.triggers {
		result = append(result, t)
	}
	return result
}

func (r *TriggerRegistry) Exists(code string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.triggers[code]
	return ok
}

func (r *TriggerRegistry) Codes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	codes := make([]string, 0, len(r.triggers))
	for code := range r.triggers {
		codes = append(codes, code)
	}
	return codes
}

func (r *TriggerRegistry) SetCommonParams(params []Param) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commonParams = params
}

func (r *TriggerRegistry) GetCommonParams() []Param {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.commonParams
}

func (r *TriggerRegistry) GetAllParams(code string) []Param {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Param
	result = append(result, r.commonParams...)
	if t, ok := r.triggers[code]; ok {
		result = append(result, t.Params...)
	}
	return result
}
