package beweb

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// templateName模板名，按名索引
	// data用于渲染页面的数据
	Render(ctx context.Context, templateName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	Engine *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, templateName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.Engine.ExecuteTemplate(bs, templateName, data)
	return bs.Bytes(), err
}
