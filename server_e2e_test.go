//go:build e2e

package beweb

import (
	"net/http"
	"testing"
)

func TestE2eServer(t *testing.T) {
	h := HTTPServer{}

	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {

	})

	h.Start(":8080")
}
