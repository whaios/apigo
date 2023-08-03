// Package apifox 开放 API（https://apifox-openapi.apifox.cn/）
package apifox

import (
	"fmt"
	"github.com/levigross/grequests"
	"strconv"
)

const (
	BaseUrl = "https://api.apifox.cn" // 开放接口地址
	Version = "2022-11-16"            // 接口版本号，当前统一填写：2022-11-16
)

var (
	ProjectId   string // 项目 ID，打开 Apifox 进入项目里的“项目设置”查看
	AccessToken string // 身份认证，个人访问令牌，查看[如何获取个人访问令牌](https://www.apifox.cn/help/openapi/)

	// ApiOverwriteMode 匹配到相同接口时的覆盖模式，不传表示忽略。
	//	枚举值: methodAndPath=覆盖，both=保留两者，merge=智能合并，ignore=不导入
	ApiOverwriteMode string
	// SchemaOverwriteMode 匹配到相同数据模型时的覆盖模式，不传表示忽略。
	//	枚举值: name=覆盖，both=保留两者，merge=智能合并，ignore=不导入
	SchemaOverwriteMode string
	// SyncApiFolder 是否同步更新接口所在目录（默认值: false）
	SyncApiFolder bool
)

const (
	// OpenApi 导入数据格式，目前只支持openapi
	OpenApi = "openapi" // 表示 Swagger 或 OpenAPI 格式

	// Apifox 扩展支持

	XFolder = "x-apifox-folder" // 接口所属目录，多级目录使用斜杠/分隔。其中\和/为特殊字符，需要转义，\/表示字符/，\\表示字符\。
	XStatus = "x-apifox-status" // 接口状态
	XOrders = "x-apifox-orders" // 属性排序
)

// PostImportData 导入接口数据
//   - importFormat 导入数据格式，目前只支持openapi
//   - jsonData Swagger（OpenAPI） 格式 json 字符串，如 ../example/openapi_data.json
func PostImportData(jsonData string) (*ImportDataResult, error) {
	if ProjectId == "" {
		return nil, fmt.Errorf("[项目ID]不能为空")
	}
	if AccessToken == "" {
		return nil, fmt.Errorf("[访问令牌]不能为空")
	}

	url := fmt.Sprintf("%s/api/v1/projects/%s/import-data", BaseUrl, ProjectId)
	headers := map[string]string{
		"X-Apifox-Version": Version,
		"Authorization":    "Bearer " + AccessToken,
	}
	data := map[string]string{
		"importFormat":  OpenApi,
		"data":          jsonData,
		"syncApiFolder": strconv.FormatBool(SyncApiFolder),
	}
	{
		if ApiOverwriteMode != "" {
			data["apiOverwriteMode"] = ApiOverwriteMode
		}
		if SchemaOverwriteMode != "" {
			data["schemaOverwriteMode"] = SchemaOverwriteMode
		}
	}
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:    data,
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}

	result := &ImportDataResult{}
	return result, resp.JSON(result)
}

// ImportDataResult 导入接口数据返回结果
type ImportDataResult struct {
	Success bool `json:"success"` // 接口状态
	Data    struct {
		ApiCollection    ImportDataResultCollection `json:"apiCollection"`    // 导入接口情况
		SchemaCollection ImportDataResultCollection `json:"schemaCollection"` // 导入数据模型情况
	} `json:"data"` // 导入结果
}

func (p *ImportDataResult) String() string {
	return fmt.Sprintf(`新增接口 %d，修改接口 %d，出错接口 %d，忽略接口 %d
新增模型 %d，修改模型 %d，出错模型 %d，忽略模型 %d`,
		p.Data.ApiCollection.Item.CreateCount,
		p.Data.ApiCollection.Item.UpdateCount,
		p.Data.ApiCollection.Item.ErrorCount,
		p.Data.ApiCollection.Item.IgnoreCount,
		p.Data.SchemaCollection.Item.CreateCount,
		p.Data.SchemaCollection.Item.UpdateCount,
		p.Data.SchemaCollection.Item.ErrorCount,
		p.Data.SchemaCollection.Item.IgnoreCount,
	)
}

type ImportDataResultCollection struct {
	Item   ImportDataResultItem `json:"item"`   // 接口
	Folder ImportDataResultItem `json:"folder"` // 接口目录
}

type ImportDataResultItem struct {
	CreateCount int `json:"createCount"` // 新增的接口数
	UpdateCount int `json:"updateCount"` // 修改的接口数
	ErrorCount  int `json:"errorCount"`  // 导入出错接口数
	IgnoreCount int `json:"ignoreCount"` // 忽略的接口数
}
