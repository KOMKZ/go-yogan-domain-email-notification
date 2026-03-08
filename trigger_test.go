package emailnotification

import (
	"sort"
	"testing"
)

func TestTriggerRegistry_Register(t *testing.T) {
	r := NewTriggerRegistry()
	r.Register("user.registered", "用户注册", "新用户注册时发送邮件", []Param{
		{Name: "UserName", Type: "string", Required: true, Example: "张三"},
	})

	trigger, ok := r.Get("user.registered")
	if !ok {
		t.Fatal("expected trigger to exist")
	}
	if trigger.Code != "user.registered" {
		t.Errorf("code = %s, want user.registered", trigger.Code)
	}
	if trigger.Name != "用户注册" {
		t.Errorf("name = %s, want 用户注册", trigger.Name)
	}
	if len(trigger.Params) != 1 {
		t.Errorf("params len = %d, want 1", len(trigger.Params))
	}
}

func TestTriggerRegistry_Get_NotFound(t *testing.T) {
	r := NewTriggerRegistry()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("expected trigger not found")
	}
}

func TestTriggerRegistry_GetAll(t *testing.T) {
	r := NewTriggerRegistry()
	r.Register("a", "A", "desc A", nil)
	r.Register("b", "B", "desc B", nil)

	all := r.GetAll()
	if len(all) != 2 {
		t.Errorf("len = %d, want 2", len(all))
	}
}

func TestTriggerRegistry_Exists(t *testing.T) {
	r := NewTriggerRegistry()
	r.Register("test", "Test", "desc", nil)

	if !r.Exists("test") {
		t.Error("expected exists")
	}
	if r.Exists("nope") {
		t.Error("expected not exists")
	}
}

func TestTriggerRegistry_Codes(t *testing.T) {
	r := NewTriggerRegistry()
	r.Register("a", "A", "desc", nil)
	r.Register("b", "B", "desc", nil)

	codes := r.Codes()
	sort.Strings(codes)
	if len(codes) != 2 || codes[0] != "a" || codes[1] != "b" {
		t.Errorf("codes = %v, want [a b]", codes)
	}
}

func TestTriggerRegistry_CommonParams(t *testing.T) {
	r := NewTriggerRegistry()
	r.SetCommonParams([]Param{
		{Name: "AppName", Type: "string", Example: "Yogan"},
	})

	params := r.GetCommonParams()
	if len(params) != 1 || params[0].Name != "AppName" {
		t.Errorf("common params = %v", params)
	}
}

func TestTriggerRegistry_GetAllParams(t *testing.T) {
	r := NewTriggerRegistry()
	r.SetCommonParams([]Param{
		{Name: "AppName", Type: "string"},
	})
	r.Register("test", "Test", "desc", []Param{
		{Name: "UserName", Type: "string"},
	})

	params := r.GetAllParams("test")
	if len(params) != 2 {
		t.Errorf("len = %d, want 2", len(params))
	}

	paramsUnknown := r.GetAllParams("unknown")
	if len(paramsUnknown) != 1 {
		t.Errorf("unknown trigger params len = %d, want 1 (common only)", len(paramsUnknown))
	}
}

func TestTriggerRegistry_Overwrite(t *testing.T) {
	r := NewTriggerRegistry()
	r.Register("a", "Old", "old desc", nil)
	r.Register("a", "New", "new desc", nil)

	trigger, _ := r.Get("a")
	if trigger.Name != "New" {
		t.Errorf("name = %s, want New", trigger.Name)
	}
}
