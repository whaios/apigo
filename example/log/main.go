package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/whaios/apigo/log"
	"time"
)

func main() {
	ctx := context.Background()
	log.StartSpinner(ctx)

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("测试日志%d", i+1)
		// 输出的日志不能比spinner的日志长度短，否则无法覆盖完上次的日志（上次的日志会有残留）
		color.Cyan(log.SpinnerString(msg))

		log.UpdateSpinner(msg)
		time.Sleep(1 * time.Second)
	}

	log.StopSpinner()
}
