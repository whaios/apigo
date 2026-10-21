// Package apifox 开放 API（https://apifox-openapi.apifox.cn/）
package apifox

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/levigross/grequests"
)

const (
	BaseUrl = "https://api.apifox.cn" // 开放接口地址
	Version = "2024-03-28"            // 接口版本号，当前统一填写：2024-03-28
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
	// Apifox 扩展支持

	XFolder = "x-apifox-folder" // 接口所属目录，多级目录使用斜杠/分隔。其中\和/为特殊字符，需要转义，\/表示字符/，\\表示字符\。
	XStatus = "x-apifox-status" // 接口状态
	XOrders = "x-apifox-orders" // 属性排序
)

// PostImportData 导入接口数据
//   - importFormat 导入数据格式，目前只支持openapi
//   - apiJson Swagger（OpenAPI） 格式 json 字符串，如 ../example/openapi_data.json
func PostImportData(apiJson string) (*ImportDataResult, error) {
	if ProjectId == "" {
		return nil, fmt.Errorf("[项目ID]不能为空")
	}
	if AccessToken == "" {
		return nil, fmt.Errorf("[访问令牌]不能为空")
	}

	url := fmt.Sprintf("%s/v1/projects/%s/import-openapi", BaseUrl, ProjectId)
	headers := map[string]string{
		"X-Apifox-Api-Version": Version,
		"Authorization":        "Bearer " + AccessToken,
	}
	options := map[string]interface{}{
		//"targetEndpointFolderId": 19591819,
	}
	if ApiOverwriteMode != "" {
		options["endpointOverwriteBehavior"] = ApiOverwriteMode
	}
	if SchemaOverwriteMode != "" {
		options["schemaOverwriteBehavior"] = SchemaOverwriteMode
	}
	data := map[string]interface{}{
		"input":   apiJson,
		"options": options,
	}
	dataJson, _ := json.Marshal(data)
	ro := &grequests.RequestOptions{
		JSON:    dataJson,
		Headers: headers,
	}
	resp, err := grequests.Post(url, grequests.FromRequestOptions(ro))
	if err != nil {
		return nil, err
	}
	result := &ImportDataResult{}
	return result, resp.JSON(result)
}

// ImportDataResult 导入接口数据返回结果
type ImportDataResult struct {
	Data struct {
		Errors   []ImportDataResultError `json:"errors"` // 错误
		Counters struct {
			EndpointCreated       int `json:"endpointCreated"`       // 新增的接口数
			EndpointUpdated       int `json:"endpointUpdated"`       // 修改的接口数
			EndpointFailed        int `json:"endpointFailed"`        // 导入出错接口数
			EndpointIgnored       int `json:"endpointIgnored"`       // 忽略的接口数
			SchemaCreated         int `json:"schemaCreated"`         // 新增的数据模型数
			SchemaUpdated         int `json:"schemaUpdated"`         // 修改的数据模型数
			SchemaFailed          int `json:"schemaFailed"`          // 导入出错数据模型数
			SchemaIgnored         int `json:"schemaIgnored"`         // 忽略的数据模型数
			EndpointFolderCreated int `json:"endpointFolderCreated"` // 新增的接口目录数
			EndpointFolderUpdated int `json:"endpointFolderUpdated"` // 修改的接口目录数
			EndpointFolderFailed  int `json:"endpointFolderFailed"`  // 导入出错接口目录数
			EndpointFolderIgnored int `json:"endpointFolderIgnored"` // 忽略的接口目录数
			SchemaFolderCreated   int `json:"schemaFolderCreated"`   // 新增的数据模型目录数
			SchemaFolderUpdated   int `json:"schemaFolderUpdated"`   // 修改的数据模型目录数
			SchemaFolderFailed    int `json:"schemaFolderFailed"`    // 导入出错数据模型目录数
			SchemaFolderIgnored   int `json:"schemaFolderIgnored"`   // 忽略的数据模型目录数
		} `json:"counters"` // 导入的接口情况统计
	} `json:"data"`
}

func (p *ImportDataResult) String() string {
	return fmt.Sprintf(`新增接口 %d，修改接口 %d，出错接口 %d，忽略接口 %d
新增模型 %d，修改模型 %d，出错模型 %d，忽略模型 %d`,
		p.Data.Counters.EndpointCreated,
		p.Data.Counters.EndpointUpdated,
		p.Data.Counters.EndpointFailed,
		p.Data.Counters.EndpointIgnored,
		p.Data.Counters.SchemaCreated,
		p.Data.Counters.SchemaUpdated,
		p.Data.Counters.SchemaFailed,
		p.Data.Counters.SchemaIgnored,
	)
}

func (p *ImportDataResult) HasError() bool {
	return len(p.Data.Errors) > 0
}

func (p *ImportDataResult) ErrorMessage() string {
	msgs := make([]string, 0)
	for _, e := range p.Data.Errors {
		msgs = append(msgs, fmt.Sprintf("[%d]%s", e.Code, e.Message))
	}
	return strings.Join(msgs, "；")
}

type ImportDataResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
