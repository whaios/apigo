package jsonschema

type Schema struct {
	Type        string   `json:"type,omitempty"`        // 对象类型
	Description string   `json:"description,omitempty"` // 说明
	Properties  Schemas  `json:"properties,omitempty"`  // 属性
	Required    []string `json:"required,omitempty"`    // 必须属性
}

type Schemas map[string]*Schema
