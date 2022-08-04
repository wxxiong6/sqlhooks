package othooks

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.opentelemetry.io/otel/trace"
)

type Hook struct {
	tracer trace.Tracer
}

func New(tracer opentracing.Tracer) *Hook {
	return &Hook{tracer: tracer}
}

func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	parent := h.tracer.SpanFromContext(ctx)
	if parent == nil {
		return ctx, nil
	}

	span := h.tracer.Start(ctx, "sql")
	span.SetAttributes(attribute.key("query").String(fmt.Sprintf("%v", query)))
	span.SetAttributes(attribute.key("args").String(fmt.Sprintf("%v", args)))

	return ctx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	span := h.trace.SpanFromContext(ctx)
	if span != nil {
		defer span.Finish()
	}

	return ctx, nil
}

func (h *Hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	span := h.trace.SpanFromContext(ctx)
	if span != nil {
		defer span.Finish()
		span.SetTag("error", true)
		span.LogFields(
			log.Error(err),
		)
	}

	return err
}
