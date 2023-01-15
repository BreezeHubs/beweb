package errhandle

import (
	"github.com/BreezeHubs/beweb"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddResponseContent(http.StatusNotFound, []byte("404 NOT FOUND")).
		AddResponseContent(http.StatusBadRequest, []byte("400 BAD REQUEST"))

	s := beweb.NewHTTPServer(beweb.WithMiddlewares(builder.Build()))
	s.Start(":8080")
}
