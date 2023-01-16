package beweb

import "context"

type TemplateEngine interface {
	// Render 渲染页面
	// templateName模板名，按名索引
	// data用于渲染页面的数据
	Render(ctx context.Context, templateName string, data any) ([]byte, error)
}
