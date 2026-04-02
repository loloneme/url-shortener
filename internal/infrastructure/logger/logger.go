package logger

import (
	"context"
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Log = slog.New(handler)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "traceID", traceID)
}

func FromContext(ctx context.Context) *slog.Logger {
	if traceID, ok := ctx.Value("traceID").(string); ok {
		return Log.With("trace.id", traceID)
	}
	return Log
}
