package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Enabled implements slog.Handler.
func (h *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

// Handle implements slog.Handler.
func (h *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

// WithAttrs implements slog.Handler.
func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup implements slog.Handler.
func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	return h
}

