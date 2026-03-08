package repository

import "testing"

func TestPageResult_Structure(t *testing.T) {
	result := PageResult[string]{
		Items:      []string{"a", "b"},
		Total:      2,
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}
	if len(result.Items) != 2 {
		t.Errorf("items len = %d", len(result.Items))
	}
	if result.Total != 2 {
		t.Errorf("total = %d", result.Total)
	}
}

func TestTemplateFilter_Defaults(t *testing.T) {
	f := TemplateFilter{}
	if f.TriggerCode != "" || f.Language != "" {
		t.Error("default filter should have empty fields")
	}
}

func TestLogFilter_Defaults(t *testing.T) {
	f := LogFilter{}
	if f.TriggerCode != "" || f.StartTime != "" {
		t.Error("default filter should have empty fields")
	}
}
