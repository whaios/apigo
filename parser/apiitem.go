package parser

import (
	"github.com/go-openapi/spec"
	"path"
	"strconv"
	"strings"
)

// ApiItem 接口文档
type ApiItem struct {
	Title       string `json:"title"`       // 必填，接口名称
	Folder      string `json:"folder"`      // 接口目录，多级目录用 / 隔开，例如 “一层/二层/三层”
	Status      string `json:"status"`      // 接口状态
	Description string `json:"description"` // 接口说明
	Remark      string `json:"remark"`      // 备注信息

	Method string `json:"method"` // 必填，http 请求方式
	Path   string `json:"path"`   // 必填，http 请求路径

	Parameters  Parameters  `json:"parameters,omitempty"` // 请求参数
	ContentType string      `json:"content_type"`         // 响应类型
	Responses   []*Response `json:"responses,omitempty"`  // 返回响应
}

// Name 文档分类+标题
func (p *ApiItem) Name() string {
	return path.Join(p.Folder, p.Title)
}

// Invalid 没有标题或Url，不是有效的API文档
func (p *ApiItem) Invalid() bool {
	return p.Title == "" || p.Method == "" || p.Path == ""
}

func (p *ApiItem) UseCommon(comm *ApiItem) {
	p.AddFolder(comm.Folder)
	p.AddRemark(comm.Remark)
	for _, header := range comm.Parameters.Header {
		p.Parameters.Header = append(p.Parameters.Header, header)
	}
	if len(comm.Responses) > 0 {
		if commSchema := comm.Responses[0].JsonSchema; commSchema != nil {
			if len(p.Responses) == 0 {
				jsonSchema := *commSchema
				p.Responses = append(p.Responses, &Response{
					Code:       comm.Responses[0].Code,
					Name:       comm.Responses[0].Name,
					JsonSchema: &jsonSchema,
				})
			} else {
				cfKey := SchemaGetComposedFieldKey(commSchema)
				if cfKey == "" {
					cfKey = "data"
				}
				for _, resp := range p.Responses {
					if resp.JsonSchema != nil {
						resp.JsonSchema = spec.ComposedSchema(*commSchema, spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type:       []string{OBJECT},
								Properties: map[string]spec.Schema{cfKey: *resp.JsonSchema},
							},
						})
					}
				}
			}
		}
	}
}

func (p *ApiItem) AddFolder(folder string) {
	p.Folder = path.Join(folder, p.Folder)
}

func (p *ApiItem) AddRemark(remark string) {
	if remark == "" {
		return
	}
	if p.Remark != "" {
		p.Remark += "\n"
	}
	p.Remark += remark
}

// Parameters 请求参数
type Parameters struct {
	Path       []Parameter  `json:"path,omitempty"`       // Path 参数
	Query      []Parameter  `json:"query,omitempty"`      // Query 参数
	Header     []Parameter  `json:"header,omitempty"`     // Header 参数
	Cookie     []Parameter  `json:"cookie,omitempty"`     // Cookie 参数
	BodyType   string       `json:"bodyType,omitempty"`   // Body 类型
	FormData   []Parameter  `json:"formData,omitempty"`   // 参数（类型为 form-data 或 x-www-form-urlencoded 时的参数）
	JsonSchema *spec.Schema `json:"jsonSchema,omitempty"` // json数据结构
}

func NewParameter(name, tpe, required, example, desc string) Parameter {
	return Parameter{
		Name:        name,
		Type:        tpe,
		Required:    strToBool(required),
		Example:     example,
		Description: desc,
	}
}

// SchemaToParameters 将对象的属性拆解为二维的参数数组
func SchemaToParameters(schema *spec.Schema) []Parameter {
	params := make([]Parameter, 0)
	if schema == nil {
		return params
	}

	required := schema.SchemaProps.Required
	newParam := func(propName string, prop spec.Schema) Parameter {
		return Parameter{
			Name:        propName,
			Type:        prop.SchemaProps.Type[0],
			Required:    strSliceContains(required, propName),
			Example:     "",
			Description: prop.Description,
		}
	}
	// 如果有排序则按排序
	if orders := SchemaGetPropertiesOrders(schema); len(orders) > 0 {
		for _, propName := range orders {
			prop, ok := schema.SchemaProps.Properties[propName]
			if !ok {
				continue
			}
			param := newParam(propName, prop)
			params = append(params, param)
		}
	} else {
		for propName, prop := range schema.SchemaProps.Properties {
			param := newParam(propName, prop)
			params = append(params, param)
		}
	}
	return params
}

// Parameter 请求参数
type Parameter struct {
	Name        string `json:"name,omitempty"`        // 参数名
	Type        string `json:"type,omitempty"`        // 数据类型
	Required    bool   `json:"required,omitempty"`    // 必填
	Example     string `json:"example,omitempty"`     // 示例值
	Description string `json:"description,omitempty"` // 说明
}

// Response 返回响应
type Response struct {
	Code       int          `json:"code"`                 // HTTP 状态码
	Name       string       `json:"name"`                 // 成功 或 失败，可自定义
	JsonSchema *spec.Schema `json:"jsonSchema,omitempty"` // 响应数据
}

const (
	MethodGet      = "get"
	MethodPost     = "post"
	MethodPut      = "put"
	MethodDelete   = "delete"
	MethodOptions  = "options"
	MethodHead     = "head"
	MethodPatch    = "patch"
	MethodTrace    = "trace"
	MethodConnect  = "connect"
	MethodCopy     = "copy"
	MethodLink     = "link"
	MethodUnlink   = "unlink"
	MethodPurge    = "purge"
	MethodLock     = "lock"
	MethodUnlock   = "unlock"
	MethodMkcol    = "mkcol"
	MethodMove     = "move"
	MethodPropfind = "propfind"
	MethodReport   = "report"
	MethodView     = "view"
)

func strToInt(s string) int {
	s = strings.TrimSpace(s)
	v, _ := strconv.ParseInt(s, 10, 64)
	return int(v)
}

func strToBool(val string) bool {
	val = strings.TrimSpace(val)
	if val == "是" {
		return true
	} else if val == "否" {
		return false
	}
	b, _ := strconv.ParseBool(val)
	return b
}

func strSliceContains(opts []string, val string) bool {
	val = strings.TrimSpace(val)
	for _, opt := range opts {
		if opt == val {
			return true
		}
	}
	return false
}
