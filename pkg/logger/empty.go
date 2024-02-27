package logger

import (
	"context"
	"log/slog"
)

func NewEmpty() *slog.Logger {
	return slog.New(newEmptyHandler())
}

type handler struct{}

func newEmptyHandler() *handler {
	return &handler{}
}

func (h *handler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *handler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

func (h *handler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

func (h *handler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}
