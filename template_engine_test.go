package email_notification

import (
	"strings"
	"testing"
)

func TestTemplateEngine_Render(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		params   map[string]any
		expected string
		wantErr  bool
	}{
		{
			name:     "simple string param",
			template: "Hello, {{.UserName}}!",
			params:   map[string]any{"UserName": "张三"},
			expected: "Hello, 张三!",
		},
		{
			name:     "multiple params",
			template: "Hi {{.UserName}}, your email is {{.Email}}",
			params:   map[string]any{"UserName": "张三", "Email": "test@example.com"},
			expected: "Hi 张三, your email is test@example.com",
		},
		{
			name:     "number param",
			template: "Year: {{.Year}}",
			params:   map[string]any{"Year": 2026},
			expected: "Year: 2026",
		},
		{
			name:     "html content",
			template: "<h1>Welcome {{.UserName}}</h1>",
			params:   map[string]any{"UserName": "李四"},
			expected: "<h1>Welcome 李四</h1>",
		},
		{
			name:     "invalid template syntax",
			template: "Hello, {{.UserName",
			params:   map[string]any{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Render(tt.template, tt.params)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTemplateEngine_Preview(t *testing.T) {
	engine := NewTemplateEngine()

	params := []Param{
		{Name: "UserName", Type: "string", Example: "张三"},
		{Name: "Email", Type: "string", Example: "test@example.com"},
		{Name: "NoExample", Type: "string"}, // 无示例值
	}

	template := "Hi {{.UserName}}, email: {{.Email}}, other: {{.NoExample}}"

	result, err := engine.Preview(template, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "张三") {
		t.Error("expected result to contain '张三'")
	}
	if !strings.Contains(result, "test@example.com") {
		t.Error("expected result to contain 'test@example.com'")
	}
	if !strings.Contains(result, "{{.NoExample}}") {
		t.Error("expected result to contain placeholder for NoExample")
	}
}
