package logs

import (
	"fmt"
	"os"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	//blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
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
