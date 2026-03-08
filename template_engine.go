package emailnotification

import (
	"bytes"
	"fmt"
	"text/template"

	domainerrors "github.com/KOMKZ/go-yogan-domain-email-notification/errors"
)

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

func (e *TemplateEngine) Render(templateStr string, params map[string]any) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domainerrors.ErrTemplateRender, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("%w: %v", domainerrors.ErrTemplateRender, err)
	}

	return buf.String(), nil
}

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
