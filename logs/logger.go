package logs

import (
	"fmt"
	"os"
	"time"
)

const (
	reset  = "\033[0m"
	purple = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	cyan   = "\033[36m"
	red    = "\033[91m"
)

func log(level, color, msg string, args ...any) {
	now := time.Now().Format("[15:04:05.000]")
	fmt.Println(color + now + " [" + level + "] " + fmt.Sprintf(msg, args...) + reset)
}

func Debug(msg string, args ...any) { log("DEBUG", cyan, msg, args...) }
func Info(msg string, args ...any)  { log("INFO", green, msg, args...) }
func Warn(msg string, args ...any)  { log("WARN", yellow, msg, args...) }
func Error(msg string, args ...any) { log("ERROR", purple, msg, args...) }
func Fatal(msg string, args ...any) { log("FATAL", red, msg, args...); os.Exit(0) }
