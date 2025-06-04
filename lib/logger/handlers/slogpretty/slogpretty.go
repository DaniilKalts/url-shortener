package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	logger *log.Logger
	attrs  []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	output io.Writer,
) *PrettyHandler {
	base := slog.NewTextHandler(output, opts.SlogOpts)

	return &PrettyHandler{
		opts:    opts,
		Handler: base,
		logger:  log.New(output, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(
		func(a slog.Attr) bool {
			fields[a.Key] = a.Value.Any()

			return true
		},
	)

	for _, attr := range h.attrs {
		fields[attr.Key] = attr.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.logger.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{Handler: h.Handler, logger: h.logger, attrs: attrs}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name), logger: h.logger,
	}
}
