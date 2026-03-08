package model

import "testing"

func TestTemplate_TableName(t *testing.T) {
	tmpl := Template{}
	if tmpl.TableName() != "email_templates" {
		t.Errorf("table name = %s", tmpl.TableName())
	}
}

func TestTemplate_IsEnabled(t *testing.T) {
	tests := []struct {
		status TemplateStatus
		want   bool
	}{
		{TemplateStatusEnabled, true},
		{TemplateStatusDraft, false},
		{TemplateStatusDisabled, false},
	}
	for _, tt := range tests {
		tmpl := Template{Status: tt.status}
		if tmpl.IsEnabled() != tt.want {
			t.Errorf("status=%s: IsEnabled()=%v, want %v", tt.status, tmpl.IsEnabled(), tt.want)
		}
	}
}
