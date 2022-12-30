package apidoc

// Generate 生成API文档
func Generate(searchDir string) error {

	return nil
}

type ApiDoc struct {
	Catalog     string `json:"catalog"`     // 目录，例如 “一层/二层/三层”
	Name        string `json:"name"`        // 文档名称
	Description string `json:"description"` // 说明
	Remark      string `json:"remark"`      // 备注

	Method string `json:"method"`
	Path   string `json:"path"`

	Parameters      Parameters  `json:"parameters,omitempty"`
	RequestBody     RequestBody `json:"requestBody,omitempty"`
	Responses       []*Response `json:"responses,omitempty"`
	ResponseExample string      `json:"responseExamples,omitempty"`
}

type Parameters struct {
	Path   []*Parameter `json:"path,omitempty"`
	Query  []*Parameter `json:"query,omitempty"`
	Header []*Parameter `json:"header,omitempty"`
	Cookie []*Parameter `json:"cookie,omitempty"`
}

type Parameter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Type        string `json:"type,omitempty"`
	Example     string `json:"example,omitempty"`
}

type RequestBody struct {
	Type       string       `json:"type,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
	JsonSchema *JSONSchema  `json:"jsonSchema,omitempty"`
	Example    string       `json:"example,omitempty"`
}

const (
	RequestBodyTypeJSON = "application/json"
	RequestBodyTypeForm = "multipart/form-data"
)

const (
	JSONSchemaTypeString  = "string"
	JSONSchemaTypeNumber  = "number"
	JSONSchemaTypeInteger = "integer"
	JSONSchemaTypeBoolean = "boolean"
	JSONSchemaTypeObject  = "object"
)

type JSONSchema struct {
	Type        string   `json:"type,omitempty"`        // 对象类型
	Description string   `json:"description,omitempty"` // 说明
	Properties  Schemas  `json:"properties,omitempty"`  // 属性
	Required    []string `json:"required,omitempty"`    // 必须属性
}

type Schemas map[string]*JSONSchema

type Response struct {
	JsonSchema *JSONSchema `json:"jsonSchema,omitempty"`
}
