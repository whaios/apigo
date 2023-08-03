package log

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var IsDebug = false

// Debug 输出调试信息
func Debug(format string, a ...interface{}) {
	if IsDebug {
		msg := fmt.Sprintf(format, a...)
		color.White(SpinnerString("[debug] " + msg))
	}
}

func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	color.Cyan(SpinnerString(msg))
}

func Warn(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	color.Yellow(SpinnerString(msg))
}

func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	color.Red(SpinnerString(msg))
}

func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	color.Green(SpinnerString(msg))
}

func Fatal(err error) {
	color.Red(SpinnerString(err.Error()))
	StopSpinner()
	os.Exit(1)
}
