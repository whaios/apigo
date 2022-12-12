package main

import (
	"github.com/urfave/cli/v2"
	"github.com/whaios/apigo/log"
	"os"
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "显示帮助",
	}
	app := cli.NewApp()
	app.Name = "Apigo"
	app.Usage = "Go API 文档工具"
	app.Description = `解析 Go 代码文件中的注释生成 API 文档。`
	app.Version = "1.0.0"

	app.Commands = []*cli.Command{}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
