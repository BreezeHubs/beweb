//go:build e2e

package template

import (
	"github.com/BreezeHubs/beweb"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)

	s := beweb.NewHTTPServer(
		beweb.WithTemplateEngine(
			&beweb.GoTemplateEngine{Engine: tpl},
		),
	)

	s.Get("/login", func(ctx *beweb.Context) {
		ctx.Render("login.gohtml", nil)
	})
	
	s.Start(":8080")
}
