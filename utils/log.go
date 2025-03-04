package utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelPanic
	InfoColor  = "\033[34m"
	WarnColor  = "\033[33m"
	ErrorColor = "\033[31m"
	ResetColor = "\033[0m"
)

func myLog(level LogLevel, format string, v ...any) {
	var levelStr string
	var color string
	switch level {
	case LogLevelInfo:
		levelStr = "INFO"
		color = InfoColor
	case LogLevelWarn:
		levelStr = "WARN"
		color = WarnColor
	case LogLevelError:
		levelStr = "ERRO"
		color = ErrorColor
	case LogLevelFatal:
		levelStr = "FATA"
		color = ErrorColor
	default:
		levelStr = "UNKN"
	}
	fmt.Printf("%s%s%s [%s] %s \n", color, levelStr, ResetColor, time.Now().Format("2006-01-02T15:04:05"), fmt.Sprintf(format, v...))
}

func LogInfo(format string, v ...any) {
	myLog(LogLevelInfo, format, v...)
}
func LogWarn(format string, v ...any) {
	myLog(LogLevelWarn, format, v...)
}
func LogError(format string, v ...any) {
	myLog(LogLevelError, format, v...)
}
func LogFatal(format string, v ...any) {
	myLog(LogLevelFatal, format, v...)
}
