package logger

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

var levelColors = map[slog.Level]string{
	slog.LevelDebug: "\033[36m", // голубой
	slog.LevelInfo:  "\033[32m", // зелёный
	slog.LevelWarn:  "\033[33m", // жёлтый
	slog.LevelError: "\033[31m", // красный
}

const reset = "\033[0m"

type ColorHandler struct {
	Level slog.Level
	mx    sync.Mutex
}

func NewColorHandler(level slog.Level) *ColorHandler {
	return &ColorHandler{
		Level: level,
		mx:    sync.Mutex{},
	}
}

func (h *ColorHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.Level
}

func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	color := levelColors[r.Level]
	timeStr := r.Time.Local().Format("2006-01-02 15:04:05")

	h.mx.Lock()
	defer h.mx.Unlock()

	fmt.Printf("%s[%s] [%s] %s%s\n", color, timeStr, r.Level, strings.ToUpper(r.Message), reset)

	var num int
	var prevLen int
	r.Attrs(func(a slog.Attr) bool {
		end := "\n"
		pref := "    "
		if num%2 == 0 {
			end = ""
		} else {
			for range 40 - prevLen {
				pref += " "
			}
		}
		prevLen, _ = fmt.Printf("%s%s: %v%s", pref, a.Key, a.Value, end)
		num += 1
		return true
	})
	if r.NumAttrs()%2 != 0 {
		fmt.Println()
	}
	fmt.Println()
	return nil
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return h
}
