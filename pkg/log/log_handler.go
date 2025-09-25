package log

import (
	"context"
	"log/slog"
)

type customKey int

const (
	LogDataKey customKey = iota
)

type MyJSONLogHandler struct {
	handler slog.Handler
}

type LogData struct {
	UserID    string
	Service   string
	ProductID string
}

func NewMyJSONLogHandler(h slog.Handler) *MyJSONLogHandler {
	return &MyJSONLogHandler{handler: h}
}

// В секции ниже добавляю методы к моей структуре, чтобы она удовлетворяла
// интерфейсу slog.Handler

func (h *MyJSONLogHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

func (h *MyJSONLogHandler) Handle(ctx context.Context, rec slog.Record) error {
	if ld, ok := ctx.Value(LogDataKey).(LogData); ok {
		if ld.UserID != "" {
			rec.Add("user_id", ld.UserID)
		}
		if ld.Service != "" {
			rec.Add("phone", ld.Service)
		}
		if ld.ProductID != "" {
			rec.Add("message", ld.ProductID)
		}
	}
	return h.handler.Handle(ctx, rec)
}

func (h *MyJSONLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *MyJSONLogHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

// Ниже находятся функции для удобного добавления данных
// в контекст логгера

func WithUserID(ctx context.Context, userID string) context.Context {
	if ld, ok := ctx.Value(LogDataKey).(LogData); ok {
		ld.UserID = userID
		return context.WithValue(ctx, LogDataKey, ld)
	}
	return context.WithValue(ctx, LogDataKey, LogData{UserID: userID})
}

func WithService(ctx context.Context, service string) context.Context {
	if ld, ok := ctx.Value(LogDataKey).(LogData); ok {
		ld.Service = service
		return context.WithValue(ctx, LogDataKey, ld)
	}
	return context.WithValue(ctx, LogDataKey, LogData{Service: service})
}

func WithProductID(ctx context.Context, productID string) context.Context {
	if ld, ok := ctx.Value(LogDataKey).(LogData); ok {
		ld.ProductID = productID
		return context.WithValue(ctx, LogDataKey, ld)
	}
	return context.WithValue(ctx, LogDataKey, LogData{ProductID: productID})
}
