package email_notification

import (
	"testing"
)

func TestTriggerRegistry_Register(t *testing.T) {
	registry := NewTriggerRegistry()

	// 注册触发点
	trigger := registry.Register("user:registered", "用户注册", "用户完成注册后发送", []Param{
		{Name: "UserName", Type: "string", Required: true},
		{Name: "Email", Type: "string", Required: true},
	})

	if trigger == nil {
		t.Fatal("expected trigger to be created")
	}
	if trigger.Code != "user:registered" {
		t.Errorf("expected code 'user:registered', got '%s'", trigger.Code)
	}
	if len(trigger.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(trigger.Params))
	}
}

func TestTriggerRegistry_Get(t *testing.T) {
	registry := NewTriggerRegistry()
	registry.Register("test:trigger", "测试", "测试触发点", nil)

	// 获取存在的触发点
	trigger, ok := registry.Get("test:trigger")
	if !ok {
		t.Fatal("expected trigger to exist")
	}
	if trigger.Name != "测试" {
		t.Errorf("expected name '测试', got '%s'", trigger.Name)
	}

	// 获取不存在的触发点
	_, ok = registry.Get("not:exists")
	if ok {
		t.Error("expected trigger not to exist")
	}
}

func TestTriggerRegistry_GetAll(t *testing.T) {
	registry := NewTriggerRegistry()
	registry.Register("trigger:a", "A", "", nil)
	registry.Register("trigger:b", "B", "", nil)

	triggers := registry.GetAll()
	if len(triggers) != 2 {
		t.Errorf("expected 2 triggers, got %d", len(triggers))
	}
}

func TestTriggerRegistry_Exists(t *testing.T) {
	registry := NewTriggerRegistry()
	registry.Register("exists", "Exists", "", nil)

	if !registry.Exists("exists") {
		t.Error("expected trigger to exist")
	}
	if registry.Exists("not:exists") {
		t.Error("expected trigger not to exist")
	}
}

func TestTriggerRegistry_Codes(t *testing.T) {
	registry := NewTriggerRegistry()
	registry.Register("code:a", "A", "", nil)
	registry.Register("code:b", "B", "", nil)

	codes := registry.Codes()
	if len(codes) != 2 {
		t.Errorf("expected 2 codes, got %d", len(codes))
	}
}

func TestTriggerRegistry_CommonParams(t *testing.T) {
	registry := NewTriggerRegistry()

	// 初始无通用参数
	if len(registry.GetCommonParams()) != 0 {
		t.Error("expected no common params initially")
	}

	// 设置通用参数
	registry.SetCommonParams([]Param{
		{Name: "AppName", Type: "string"},
		{Name: "AppLogo", Type: "url"},
	})

	params := registry.GetCommonParams()
	if len(params) != 2 {
		t.Errorf("expected 2 common params, got %d", len(params))
	}
}

func TestTriggerRegistry_GetAllParams(t *testing.T) {
	registry := NewTriggerRegistry()

	// 设置通用参数
	registry.SetCommonParams([]Param{
		{Name: "AppName", Type: "string"},
	})

	// 注册触发点
	registry.Register("user:registered", "用户注册", "", []Param{
		{Name: "UserName", Type: "string"},
		{Name: "Email", Type: "string"},
	})

	// 获取全部参数
	params := registry.GetAllParams("user:registered")
	if len(params) != 3 {
		t.Errorf("expected 3 params (1 common + 2 trigger), got %d", len(params))
	}

	// 第一个应该是通用参数
	if params[0].Name != "AppName" {
		t.Errorf("expected first param 'AppName', got '%s'", params[0].Name)
	}

	// 不存在的触发点应该只返回通用参数
	params = registry.GetAllParams("not:exists")
	if len(params) != 1 {
		t.Errorf("expected 1 common param for non-existent trigger, got %d", len(params))
	}
}
