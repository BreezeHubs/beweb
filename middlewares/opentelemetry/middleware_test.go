//go:build e2e

package opentelemetry

import (
	"github.com/BreezeHubs/beweb"
	"github.com/BreezeHubs/beweb/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"os"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := NewMiddlewareBuilder().SetTrace(tracer).Build()

	s := beweb.NewHTTPServer(
		beweb.WithMiddlewares(builder),
	)

	s.Get("/user", func(ctx *beweb.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		secondC, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)

		_, third1 := tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		_, third2 := tracer.Start(secondC, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()

		second.End()

		_, first := tracer.Start(ctx.Req.Context(), "first_layer_1")
		defer first.End()

		time.Sleep(100 * time.Millisecond)

		util.ResponseJSONSuccess(ctx, struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   1,
			Name: "breeze",
		})
	})

	initZipkin(t)
	//访问：http://localhost:9411/zipkin/

	s.Start(":8080")
}

func initZipkin(t *testing.T) {
	url := "http://localhost:9411/api/v2/spans"
	serviceName := "opentelemetry-demo"

	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(
			log.New(os.Stderr, serviceName, log.Ldate|log.Ltime|log.Llongfile),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(
			sdktrace.NewBatchSpanProcessor(exporter),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)
	otel.SetTracerProvider(tp)
}

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	serviceName := "opentelemetry-demo"

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)
	otel.SetTracerProvider(tp)
}
