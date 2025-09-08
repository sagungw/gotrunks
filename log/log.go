package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

type handler struct {
	next              slog.Handler
	contextAttributes []*contextAttribute
}

func InitLogger(configOpts ...HandlerConfigFunc) *slog.Logger {
	handlerConfig := &handlerConfig{
		writer: os.Stdout,
		level:  slog.LevelInfo,
	}

	for _, opt := range configOpts {
		opt(handlerConfig)
	}

	logger := slog.New(&handler{
		next: slog.NewJSONHandler(handlerConfig.writer, &slog.HandlerOptions{
			Level: handlerConfig.level,
		}),
		contextAttributes: handlerConfig.contextAttributes,
	})

	slog.SetDefault(logger)

	return logger
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.next.Enabled(ctx, l)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	frs := runtime.CallersFrames([]uintptr{r.PC})
	fr, _ := frs.Next()
	if fr.Function != "" {
		r.AddAttrs(slog.String("func", fr.Function))
		r.AddAttrs(slog.String("file", fmt.Sprintf("%s:%d", fr.File, fr.Line)))
	}

	h.addContextAttributes(ctx, &r)

	return h.next.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{next: h.next.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{next: h.next.WithGroup(name)}
}

func (h *handler) addContextAttributes(ctx context.Context, r *slog.Record) {
	for _, ca := range h.contextAttributes {
		val := ctx.Value(ca.contextKey)
		switch v := val.(type) {
		case string:
			if v != "" {
				r.AddAttrs(slog.String(ca.logKey, v))
			}
		case int64:
			r.AddAttrs(slog.Int64(ca.logKey, v))
		case int:
			r.AddAttrs(slog.Int(ca.logKey, v))
		case float64:
			r.AddAttrs(slog.Float64(ca.logKey, v))
		}
	}
}

type handlerConfig struct {
	writer            io.Writer
	level             slog.Level
	contextAttributes []*contextAttribute
}

type contextAttribute struct {
	logKey     string
	contextKey any
}

type HandlerConfigFunc func(cfg *handlerConfig)

func WithWriter(w io.Writer) HandlerConfigFunc {
	return func(cfg *handlerConfig) {
		cfg.writer = w
	}
}

func WithLevel(level string) HandlerConfigFunc {
	return func(cfg *handlerConfig) {
		logLevel := cfg.level
		err := logLevel.UnmarshalText([]byte(level))
		if err != nil {
			return
		}

		cfg.level = logLevel
	}
}

func WithContextAttribute(logKey string, contextKey any) HandlerConfigFunc {
	return func(cfg *handlerConfig) {
		if cfg.contextAttributes == nil {
			cfg.contextAttributes = []*contextAttribute{}
		}

		cfg.contextAttributes = append(cfg.contextAttributes, &contextAttribute{
			logKey:     logKey,
			contextKey: contextKey,
		})
	}
}
