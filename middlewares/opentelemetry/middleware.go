package opentelemetry

import (
	"github.com/BreezeHubs/beweb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/BreezeHubs/beweb/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	trace trace.Tracer
}

func NewMiddlewareBuilder(trace trace.Tracer) *MiddlewareBuilder {
	if trace == nil {
		return &MiddlewareBuilder{
			trace: otel.GetTracerProvider().Tracer(instrumentationName),
		}
	}
	return &MiddlewareBuilder{trace: trace}
}

func (m MiddlewareBuilder) Build() beweb.Middleware {
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			reqCtx := ctx.Req.Context()
			//考虑上下游span
			//尝试和客户端的trace结合，上游trace id放在http header => propagation.HeaderCarrier
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			//
			_, span := m.trace.Start(reqCtx, "unknown")
			defer span.End()

			//数据
			//next执行前能拿到的数据
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))

			next(ctx)

			//next执行后能拿到的数据
			span.SetName(ctx.MatchedRoute)
		}
	}
}
