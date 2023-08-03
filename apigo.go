package main

import (
	"context"
	"encoding/json"
	"github.com/urfave/cli/v2"
	"github.com/whaios/apigo/apifox"
	"github.com/whaios/apigo/log"
	"github.com/whaios/apigo/parser"
	"io/fs"
	"os"
)

// 环境变量
const (
	EnvApifoxProjectId   = "ApifoxProjectId"   // 项目 ID
	EnvApifoxAccessToken = "ApifoxAccessToken" // 个人访问令牌
)

const (
	flagDir     = "dir"
	flagOutFile = "outfile"
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "显示帮助",
	}
	app := cli.NewApp()
	app.Name = "Apigo"
	app.Usage = "快速生成 Go (Golang) API 文档"
	app.Description = `Apigo 主要用于解析 Go (Golang) 代码注释，快速生成 API 文档，并同步到 Apifox，实现代码零入侵。`
	app.Version = "v1.0.0"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "开启调试模式。",
			Value:       log.IsDebug,
			Destination: &log.IsDebug,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "apifox",
			Aliases: []string{"af"},
			Usage:   "快速生成 API 文档，并同步到 Apifox。",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "project",
					Aliases:     []string{"p"},
					Usage:       "项目 ID",
					Value:       apifox.ProjectId,
					Destination: &apifox.ProjectId,
					EnvVars:     []string{EnvApifoxProjectId},
				},
				&cli.StringFlag{
					Name:        "token",
					Aliases:     []string{"t"},
					Usage:       "个人访问令牌",
					Value:       apifox.AccessToken,
					Destination: &apifox.AccessToken,
					EnvVars:     []string{EnvApifoxAccessToken},
				},
				&cli.StringFlag{
					Name:    flagDir,
					Aliases: []string{"d"},
					Value:   "",
					Usage:   "要解析的 Go 源码文件的目录，该目录下必须有 Go 源码文件。",
					//Required: true,
				},
				&cli.StringFlag{
					Name:    flagOutFile,
					Aliases: []string{"of"},
					Value:   "",
					Usage:   "将生成的文档数据导出到指定文件，不上传到 Apifox。",
				},
				&cli.StringFlag{
					Name:        "apiOverwriteMode",
					Value:       "methodAndPath",
					Usage:       "匹配到相同接口时的覆盖模式，不传表示忽略。枚举值: methodAndPath=覆盖，both=保留两者，merge=智能合并，ignore=不导入",
					Destination: &apifox.ApiOverwriteMode,
				},
				&cli.StringFlag{
					Name:        "schemaOverwriteMode",
					Value:       "",
					Usage:       "匹配到相同数据模型时的覆盖模式，不传表示忽略。枚举值: name=覆盖，both=保留两者，merge=智能合并，ignore=不导入",
					Destination: &apifox.SchemaOverwriteMode,
				},
				&cli.BoolFlag{
					Name:        "syncApiFolder",
					Value:       false,
					Usage:       "是否同步更新接口所在目录（默认值: false）",
					Destination: &apifox.SyncApiFolder,
				},
			},
			Action: func(c *cli.Context) error {
				apifoxImptData(c.Context, c.String(flagDir), c.String(flagOutFile))
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:  "flags",
					Usage: "查询相关参数。",
					Action: func(c *cli.Context) error {
						log.Info("baseUrl=%s", apifox.BaseUrl)
						log.Info("projectId=%s", apifox.ProjectId)
						log.Info("accessToken=%s", apifox.AccessToken)
						log.Info("isDebug=%v", log.IsDebug)
						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// 解析 go 注释生成 API 文档，并导入到 apifox。
func apifoxImptData(ctx context.Context, dir, outFile string) {
	log.StartSpinner(ctx)
	defer log.StopSpinner()

	log.Info("扫描目录 %s", dir)

	parser.SetSchemaExtraPropertiesOrdersKey(apifox.XOrders)
	goParser := parser.NewParser()
	fileCount, err := goParser.Scan(dir)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("采集到%d个Go代码文件", fileCount)

	// 解析每个文件并生成接口文档
	items, err := goParser.Parse()
	if err != nil {
		log.Error(err.Error())
		return
	}

	api2 := apifox.NewOpenApi2()
	apifox.OpenApi2AddPaths(api2, items)

	api2JsonData, err := json.MarshalIndent(api2, "", "    ")
	if err != nil {
		log.Error(err.Error())
		return
	}

	// 导出到文件
	if outFile != "" {
		log.Debug(log.UpdateSpinner("导出到文件 %s", outFile))

		if err = os.WriteFile(outFile, api2JsonData, fs.ModePerm); err != nil {
			log.Error(err.Error())
			return
		}
		log.Success("导出文件成功 %s", outFile)
		return
	}

	log.Debug(log.UpdateSpinner("同步到 Apifox"))

	// 上传到 Apifox 服务器
	result, err := apifox.PostImportData(string(api2JsonData))
	if err != nil {
		log.Error(err.Error())
		return
	}
	if result.Success {
		log.Info(result.String())
		log.Success("同步 Apifox 成功")
	} else {
		log.Error("同步 Apifox 失败")
	}
	return
}
