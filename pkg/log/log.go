package log

import (
	"io"
	golog "log"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	consoleWriter := newConsoleWriter(os.Stderr)
	log.Logger = log.With().Caller().Logger().Output(consoleWriter)
}

func AddFileOutput(logPath string) {
	if logPath == "" {
		return
	}

	consoleWriter := newConsoleWriter(os.Stderr)
	// output log to file
	logFileWriter := createRollingLogFile(logPath)
	if logFileWriter != nil {
		formatWriter := newFileFormatWriter(logFileWriter)
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, formatWriter)
		log.Logger = log.With().Caller().Logger().Output(multiWriter)
	} else {
		log.Logger = log.With().Caller().Logger().Output(consoleWriter)
	}
}

func SetDebugLevel() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func SetTraceLevel() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func Print(v ...interface{}) {
	log.Print(v)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v)
}

func createRollingLogFile(logPath string) io.Writer {
	logDir := path.Dir(logPath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		golog.Panicf("Can't create log directory. path: %s, err: %s", logPath, err)
	}

	log.Info().Msgf("Log file: %s", logPath)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //daysdays
	}
}

func createAccessRollingLogFile(logPath string) io.Writer {
	logDir := path.Dir(logPath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		golog.Panicf("Can't create log directory. path: %s, err: %s", logPath, err)
	}

	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 31,
		MaxAge:     7, //daysdays
	}
}
