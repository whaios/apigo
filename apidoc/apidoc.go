package apidoc

import (
	"path"
)

func NewApiDoc(commdoc *ApiDoc) *ApiDoc {
	doc := &ApiDoc{}

	if commdoc != nil {
		doc.Catalog = commdoc.Catalog
		doc.Remark = commdoc.Remark
		for _, header := range commdoc.Parameters.Header {
			doc.Parameters.Header = append(doc.Parameters.Header, header)
		}
	}
	return doc
}

// ApiDoc 接口文档
type ApiDoc struct {
	Title       string `json:"title"`       // 接口名称
	Catalog     string `json:"catalog"`     // 接口目录，多级目录用 / 隔开，例如 “一层/二层/三层”
	Description string `json:"description"` // 接口说明
	Remark      string `json:"remark"`      // 备注信息

	Method string `json:"method"` // http 请求方式
	Path   string `json:"path"`   // http 请求路径

	Parameters      Parameters  `json:"parameters,omitempty"`       // 请求参数
	RequestBody     RequestBody `json:"requestBody,omitempty"`      // 请求 Body 参数
	Responses       []*Response `json:"responses,omitempty"`        // 返回响应
	ResponseExample string      `json:"responseExamples,omitempty"` // 响应示例
}

// Name 文档分类+标题
func (p *ApiDoc) Name() string {
	return path.Join(p.Catalog, p.Title)
}

// Invalid 没有标题或Url，不是有效的API文档
func (p *ApiDoc) Invalid() bool {
	return p.Title == "" || p.Method == "" || p.Path == ""
}

// Parameters 请求参数
type Parameters struct {
	Path   []*Parameter `json:"path,omitempty"`   // Path 参数
	Query  []*Parameter `json:"query,omitempty"`  // Query 参数
	Header []*Parameter `json:"header,omitempty"` // Header 参数
	Cookie []*Parameter `json:"cookie,omitempty"` // Cookie 参数
}

func NewParameter(name, tpe, required, example, desc string) *Parameter {
	return &Parameter{
		Name:        name,
		Type:        tpe,
		Required:    StrToBool(required),
		Example:     example,
		Description: desc,
	}
}

func NewRespParameter(name, tpe, desc string) *Parameter {
	return &Parameter{
		Name:        name,
		Type:        tpe,
		Description: desc,
	}
}

// Parameter 请求参数
type Parameter struct {
	Name        string `json:"name,omitempty"`        // 参数名
	Type        string `json:"type,omitempty"`        // 类型
	Required    bool   `json:"required,omitempty"`    // 必填
	Example     string `json:"example,omitempty"`     // 示例值
	Description string `json:"description,omitempty"` // 说明
}

// RequestBody 请求 Body 参数
type RequestBody struct {
	Type       string       `json:"type,omitempty"`       // 类型
	Parameters []*Parameter `json:"parameters,omitempty"` // 参数（类型为 form-data 或 x-www-form-urlencoded 时的参数）
	JsonSchema *JSONSchema  `json:"jsonSchema,omitempty"` // json数据结构
	Example    string       `json:"example,omitempty"`    // 示例值
}

func NewResponse(name string) *Response {
	return &Response{
		Name: name,
	}
}

// Response 返回响应
type Response struct {
	Name       string      `json:"name"` // 成功 或 失败，可自定义
	JsonSchema *JSONSchema `json:"jsonSchema,omitempty"`
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

const (
	RequestBodyTypeNone       = "none"
	RequestBodyTypeJSON       = "application/json"
	RequestBodyTypeFormData   = "multipart/form-data"
	RequestBodyTypeUrlEncoded = "application/x-www-form-urlencoded"
)

const (
	JSONSchemaTypeString  = "string"
	JSONSchemaTypeNumber  = "number"
	JSONSchemaTypeInteger = "integer"
	JSONSchemaTypeBoolean = "boolean"
	JSONSchemaTypeArray   = "array"
	JSONSchemaTypeObject  = "object"
)

type JSONSchema struct {
	Type        string   `json:"type,omitempty"`        // 对象类型
	Description string   `json:"description,omitempty"` // 说明
	Properties  Schemas  `json:"properties,omitempty"`  // 属性
	Required    []string `json:"required,omitempty"`    // 必须属性
}

type Schemas map[string]*JSONSchema
