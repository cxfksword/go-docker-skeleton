package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/shiena/ansicolor"
)

const (
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	reset  = "\033[0m"
)

type consoleWriter struct {
	l io.Writer
}

func newConsoleWriter(out io.Writer) consoleWriter {
	return consoleWriter{l: ansicolor.NewAnsiColorWriter(out)}
}

func (w consoleWriter) Write(p []byte) (n int, err error) {
	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	var msgInfo interface{} = ""
	if evt[zerolog.MessageFieldName] != nil {
		msgInfo = evt[zerolog.MessageFieldName]
	}

	var errInfo interface{} = ""
	if evt[zerolog.ErrorFieldName] != nil {
		errInfo = evt[zerolog.ErrorFieldName]
	}
	level := strings.ToUpper(evt[zerolog.LevelFieldName].(string))
	newformat := fmt.Sprintf("%-24s %-10s %-25s %s%s\n", evt[zerolog.TimestampFieldName], w.formatLevel(level), w.formatCaller(evt[zerolog.CallerFieldName]), msgInfo, errInfo)
	if level == "TRACE" || level == "DEBUG" {
		newformat = fmt.Sprintf("%s%-24s %-10s %-25s %s%s%s\n", gray, evt[zerolog.TimestampFieldName], w.formatLevel(level), w.formatCaller(evt[zerolog.CallerFieldName]), msgInfo, errInfo, reset)
	}
	_, err = w.l.Write([]byte(newformat))

	return len(p), err
}

func (w consoleWriter) formatCaller(i interface{}) string {
	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, c); err == nil {
				c = rel
			}
		}
	}
	return fmt.Sprintf("%s>", c)
}

func (w consoleWriter) formatLevel(level string) string {
	switch level {
	case "WARN":
		return fmt.Sprintf("%s[%s]%s", yellow, level, reset)
	case "ERROR":
		fallthrough
	case "FATAL":
		return fmt.Sprintf("%s[%s]%s", red, level, reset)
	default:
		return fmt.Sprintf("[%s]", level)
	}
}
