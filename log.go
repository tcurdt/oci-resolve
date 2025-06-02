package main

import (
	"fmt"
	"os"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

var levelNames = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

var currentLevel Level = Info

func SetLogLevel(level Level) {
	currentLevel = level
}

func Log(level Level, format string, args ...any) {
	if level < currentLevel {
		return
	}
	fmt.Fprintf(os.Stderr, "[%s] ", levelNames[level])
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}

func LogDebug(format string, args ...any) {
	Log(Debug, format, args...)
}

func LogInfo(format string, args ...any) {
	Log(Info, format, args...)
}

func LogWarn(format string, args ...any) {
	Log(Warn, format, args...)
}

func LogError(format string, args ...any) {
	Log(Error, format, args...)
}
