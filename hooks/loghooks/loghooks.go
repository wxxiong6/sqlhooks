package loghooks

import (
	"context"
	"time"

	"github.com/wxxiong6/kratos-pkg/zap_log"
	"go.uber.org/zap"
)

type sqlDurationKey struct{}

type Hook struct {
	log                *zap_log.ZapLogger
	IsPrintSQLDuration bool
}

func buildQueryArgsFields(query string, args ...interface{}) []zap.Field {
	if len(args) == 0 {
		return []zap.Field{zap.String("query", query)}
	}
	return []zap.Field{zap.String("query", query), zap.Any("args", args)}
}

func New() *Hook {
	return &Hook{
		log: zap_log.Logger(),
	}
}

func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if h == nil || h.log == nil {
		return ctx, nil
	}
	//h.log.ZapLog.Info("log before sql exec", buildQueryArgsFields(query, args...)...)

	if h.IsPrintSQLDuration {
		ctx = context.WithValue(ctx, (*sqlDurationKey)(nil), time.Now())
	}
	return ctx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if h == nil || h.log == nil {
		return ctx, nil
	}

	var durationField = zap.Skip()
	if v, ok := ctx.Value((*sqlDurationKey)(nil)).(time.Time); ok {
		durationField = zap.Duration("duration", time.Now().Sub(v))
	}

	h.log.ZapLog.With(durationField).Info("log after sql exec", buildQueryArgsFields(query, args...)...)
	return ctx, nil
}

func (h *Hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	if h == nil || h.log == nil {
		return nil
	}
	h.log.ZapLog.With(zap.Error(err)).Error("log after err happened", buildQueryArgsFields(query, args...)...)
	return nil
}
