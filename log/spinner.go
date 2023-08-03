package log

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/tj/go-spin"
	"sync"
	"time"
)

// VT100 终端控制码
const (
	eraseLine = "\033[2K" // 清除当前行内容
)

var instance = &Spinner{
	spinner:     spin.New(),
	placeholder: "处理中",
}

func StartSpinner(ctx context.Context) {
	instance.Start(ctx)
}

func StopSpinner() {
	instance.Stop()
}

func UpdateSpinner(format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)
	instance.Update(msg)
	return msg
}

// SpinnerString 输出并且要保留的日志不能比spinner的日志长度短，否则无法覆盖完上次的日志（上次的日志会有残留）
func SpinnerString(msg string) string {
	return fmt.Sprintf("%s\r%s\n", eraseLine, msg)
}

type Spinner struct {
	spinner *spin.Spinner

	mu          sync.Mutex
	stop        func()
	placeholder string
}

func (s *Spinner) print() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Print(color.BlueString("\r%s %s", s.spinner.Next(), s.placeholder))
}

func (s *Spinner) Start(ctx context.Context) {
	// 使用WithCancel派生一个可被取消的ctx，用来控制后台协程。
	ctx, s.stop = context.WithCancel(ctx)
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.print()
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if s.stop != nil {
		s.stop()
		// 用空白字符把动画清除掉
		fmt.Print(eraseLine)
	}
}

func (s *Spinner) Update(placeholder string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.placeholder = placeholder
}
