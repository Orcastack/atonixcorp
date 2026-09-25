package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger interface {
	Info(...any)
	Error(...any)
	Debug(...any)
}

type logger struct {
	level string
	out   *log.Logger
}

func New(level string) Logger {
	return &logger{
		level: level,
		out:   log.New(os.Stdout, "", 0),
	}
}

func (l *logger) Info(v ...any) {
	l.out.Println(colorBlue("[INFO]"), timestamp(), fmt.Sprint(v...))
}

func (l *logger) Error(v ...any) {
	l.out.Println(colorRed("[ERROR]"), timestamp(), fmt.Sprint(v...))
}

func (l *logger) Debug(v ...any) {
	if l.level == "debug" {
		l.out.Println(colorYellow("[DEBUG]"), timestamp(), fmt.Sprint(v...))
	}
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ─────────────────────────────────────────────
// Colors (ANSI)
// ─────────────────────────────────────────────

func colorRed(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func colorBlue(s string) string {
	return "\033[34m" + s + "\033[0m"
}

func colorYellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}
