package prometheus

import (
	"github.com/BreezeHubs/beweb"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	namespace  string
	subsystem  string
	name       string
	objectives map[float64]float64
}

func NewMiddlewareBuilder(namespace, subsystem, name string, objectives map[float64]float64) *MiddlewareBuilder {
	m := &MiddlewareBuilder{
		namespace: "namespace",
		subsystem: "subsystem",
		name:      "name",
		objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}
	if namespace != "" {
		m.namespace = namespace
	}
	if subsystem != "" {
		m.namespace = subsystem
	}
	if name != "" {
		m.namespace = name
	}
	if objectives != nil && len(objectives) > 0 {
		m.objectives = objectives
	}
	return m
}

func (m MiddlewareBuilder) Build() beweb.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  m.namespace,
		Subsystem:  m.subsystem,
		Name:       m.name,
		Objectives: m.objectives,
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(vector) //注册
	
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Microseconds()

				pattern := ctx.MatchedRoute
				if pattern == "" {
					pattern = "unknown"
				}

				vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(ctx.ResponseStatus)).Observe(float64(duration))
			}()

			next(ctx)
		}
	}
}
