package state

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFile struct {
	slog.Handler
}

func (l LogFile) Enabled(c context.Context, lvl slog.Level) bool {
	return lvl == slog.LevelInfo
}

func (l LogFile) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := l.Handler.WithAttrs(attrs)
	return LogFile{newHandler}
}

func (l LogFile) WithGroup(name string) slog.Handler {
	newHandler := l.Handler.WithGroup(name)
	return LogFile{newHandler}
}

func NewMultilogger(ErrorFile, LogInfoFile *os.File) *slog.Logger {
	logRotator := &lumberjack.Logger{
		Filename:   LogInfoFile.Name(),
		MaxSize:    200,
		MaxBackups: 2,
		MaxAge:     15,
		Compress:   true,
	}

	LF := LogFile{
		slog.NewJSONHandler(logRotator, nil),
	}

	logErrorRotator := &lumberjack.Logger{
		Filename:   ErrorFile.Name(),
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     15,
		Compress:   true,
	}

	ErrorHandler := NewHandler(logErrorRotator, slog.LevelWarn)
	return slog.New(slog.NewMultiHandler(LF, ErrorHandler))
}

func NewHandler(FileToLog io.Writer, level slog.Level) *slog.JSONHandler {
	return slog.NewJSONHandler(FileToLog, &slog.HandlerOptions{Level: level})
}

func NewLogError(logger *slog.Logger, t time.Time, message string, err error) {
	logger.Error(message, "Date", t, "Error", err)
}

func NewLogClosedInfo(logger *slog.Logger, t time.Time, message string, err error) {
	logger.Info(message, "Time_of_close", t, "Cause", err)
}
