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

func Log(level Level, format string, args ...interface{}) {
  // fmt.Fprintf(os.Stderr, "[%s] ", levelNames[level])
  // fmt.Fprintln(os.Stderr, args...)
  fmt.Fprintf(os.Stderr, "[%s] ", levelNames[level])
  fmt.Fprint(os.Stderr, fmt.Sprintf(format, args...))
  fmt.Fprintln(os.Stderr)
}

func LogDebug(format string, args ...interface{}) {
  Log(Debug, format, args...)
}

func LogInfo(format string, args ...interface{}) {
  Log(Info, format, args...)
}

func LogWarn(format string, args ...interface{}) {
  Log(Warn, format, args...)
}

func LogError(format string, args ...interface{}) {
  Log(Error, format, args...)
}