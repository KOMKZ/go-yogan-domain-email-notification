package email_notification

import (
	"bytes"
	"text/template"
)

// TemplateEngine 模板渲染引擎
type TemplateEngine struct{}

// NewTemplateEngine 创建模板引擎
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Render 渲染模板
func (e *TemplateEngine) Render(templateStr string, params map[string]any) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", ErrTemplateRender.Wrap(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", ErrTemplateRender.Wrap(err)
	}

	return buf.String(), nil
}

// Preview 预览模板（使用示例值）
func (e *TemplateEngine) Preview(templateStr string, params []Param) (string, error) {
	exampleParams := make(map[string]any)
	for _, p := range params {
		if p.Example != "" {
			exampleParams[p.Name] = p.Example
		} else {
			exampleParams[p.Name] = "{{." + p.Name + "}}"
		}
	}
	return e.Render(templateStr, exampleParams)
}
