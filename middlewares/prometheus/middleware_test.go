package prometheus

import (
	"github.com/BreezeHubs/beweb"
	"github.com/BreezeHubs/beweb/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder(
		"breeze",
		"web",
		"http_response",
		nil,
	).Build()

	s := beweb.NewHTTPServer(
		beweb.WithMiddlewares(builder),
	)

	s.Get("/user", func(ctx *beweb.Context) {
		//暂停随机时间，查看监控
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)

		util.ResponseJSONSuccess(ctx, struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   1,
			Name: "breeze",
		})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()

	s.Start(":8080")
}
