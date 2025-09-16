package utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	ErrorLevel LogLevel = iota
	WarningLevel
	InfoLevel
	SuccessLevel
	CommandLevel
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.log(SuccessLevel, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WarningLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Command displays the command that will be or was executed
func (l *Logger) Command(format string, args ...interface{}) {
	l.log(CommandLevel, format, args...)
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	timeStr := time.Now().Format("15:04:05")
	var levelStr string

	switch level {
	case CommandLevel:
		levelStr = "\033[35mCMD\033[0m" // Magenta for commands
	case InfoLevel:
		levelStr = "\033[34mINFO\033[0m" // Blue
	case SuccessLevel:
		levelStr = "\033[32mOK\033[0m" // Green
	case WarningLevel:
		levelStr = "\033[33mWARN\033[0m" // Yellow
	case ErrorLevel:
		levelStr = "\033[31mERROR\033[0m" // Red
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s: %s\n", timeStr, levelStr, message)
}
