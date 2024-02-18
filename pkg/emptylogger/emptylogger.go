package emptylogger

import (
	"context"
	"log/slog"
)

func New() *slog.Logger {
	return slog.New(NewEmptyHandler())
}

type EmptyHandler struct{}

func NewEmptyHandler() *EmptyHandler {
	return &EmptyHandler{}
}

func (h *EmptyHandler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *EmptyHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

func (h *EmptyHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

func (h *EmptyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}
