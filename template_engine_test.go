package emailnotification

import (
	"errors"
	"strings"
	"testing"

	domainerrors "github.com/KOMKZ/go-yogan-domain-email-notification/errors"
)

func TestTemplateEngine_Render(t *testing.T) {
	e := NewTemplateEngine()

	tests := []struct {
		name     string
		tmpl     string
		params   map[string]any
		expected string
		wantErr  bool
	}{
		{
			name:     "simple",
			tmpl:     "Hello {{.Name}}",
			params:   map[string]any{"Name": "World"},
			expected: "Hello World",
		},
		{
			name:     "multiple params",
			tmpl:     "{{.Greeting}}, {{.Name}}!",
			params:   map[string]any{"Greeting": "Hi", "Name": "Test"},
			expected: "Hi, Test!",
		},
		{
			name:     "number param",
			tmpl:     "Year: {{.Year}}",
			params:   map[string]any{"Year": 2026},
			expected: "Year: 2026",
		},
		{
			name:     "html content",
			tmpl:     "<h1>{{.Title}}</h1>",
			params:   map[string]any{"Title": "Hello"},
			expected: "<h1>Hello</h1>",
		},
		{
			name:    "invalid syntax",
			tmpl:    "{{.Name",
			params:  map[string]any{},
			wantErr: true,
		},
		{
			name:     "empty params",
			tmpl:     "static text",
			params:   map[string]any{},
			expected: "static text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Render(tt.tmpl, tt.params)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				if !errors.Is(err, domainerrors.ErrTemplateRender) {
					t.Errorf("expected ErrTemplateRender, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("result = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTemplateEngine_Preview(t *testing.T) {
	e := NewTemplateEngine()

	params := []Param{
		{Name: "UserName", Example: "张三"},
		{Name: "Code", Example: ""},
	}

	result, err := e.Preview("Hello {{.UserName}}, code: {{.Code}}", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "张三") {
		t.Error("expected example value for UserName")
	}
	if !strings.Contains(result, "{{.Code}}") {
		t.Error("expected placeholder for Code without example")
	}
}

func TestTemplateEngine_Preview_AllExamples(t *testing.T) {
	e := NewTemplateEngine()
	params := []Param{
		{Name: "A", Example: "val_a"},
		{Name: "B", Example: "val_b"},
	}
	result, err := e.Preview("{{.A}} {{.B}}", params)
	if err != nil {
		t.Fatal(err)
	}
	if result != "val_a val_b" {
		t.Errorf("result = %q", result)
	}
}

func TestTemplateEngine_Preview_Empty(t *testing.T) {
	e := NewTemplateEngine()
	result, err := e.Preview("static", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != "static" {
		t.Errorf("result = %q", result)
	}
}

func TestTemplateEngine_Render_ExecuteError(t *testing.T) {
	e := NewTemplateEngine()
	// {{.Name}} with a map that has no Name key will succeed with <no value>
	// But calling a method on nil will fail
	_, err := e.Render("{{.Name.Invalid}}", map[string]any{"Name": "text"})
	if err == nil {
		t.Error("expected error for invalid field access")
	}
	if !errors.Is(err, domainerrors.ErrTemplateRender) {
		t.Errorf("expected ErrTemplateRender, got %v", err)
	}
}
