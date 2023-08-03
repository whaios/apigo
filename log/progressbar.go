package log

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

// ProgressWidth 进度条长度
const ProgressWidth = 40

// DrawProgressBar 显示进度条
func DrawProgressBar(prefix string, val, max int) {
	proportion := float32(val) / float32(max)
	pos := int(proportion * ProgressWidth)
	s := fmt.Sprintf("%s [%s%*s] %6.2f%% \t[%d/%d]",
		prefix, strings.Repeat("■", pos), ProgressWidth-pos, "", proportion*100, val, max)
	fmt.Print(color.CyanString("\r" + s))
	if proportion >= 1 {
		fmt.Print("\n")
	}
}
