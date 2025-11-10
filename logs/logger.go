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

var now = time.Now().Format("[15:04:05.000]")

func Debug(msg string, args ...any) {
	fmt.Println(cyan + now + " [DEBUG] " + fmt.Sprintf(msg, args...) + reset)
}
func Info(msg string, args ...any) {
	fmt.Println(green + now + " [INFO] " + fmt.Sprintf(msg, args...) + reset)
}
func Warn(msg string, args ...any) {
	fmt.Println(yellow + now + " [WARN] " + fmt.Sprintf(msg, args...) + reset)
}
func Error(msg string, args ...any) {
	fmt.Println(purple + now + " [ERROR] " + fmt.Sprintf(msg, args...) + reset)
}
func Fatal(msg string, args ...any) {
	fmt.Println(red + now + " [FATAL] " + fmt.Sprintf(msg, args...) + reset)
	os.Exit(0)
}
