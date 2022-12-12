package jsonschema

import (
	"reflect"
	"strings"
)

// ReflectFromType 根据 Go Struct 生成 JSON Schema。
func ReflectFromType(t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // re-assign from pointer
	}

	s := new(Schema)
	switch t.Kind() {
	case reflect.Struct:
		r.reflectStruct(definitions, t, st)

	case reflect.Slice, reflect.Array:
		r.reflectSliceOrArray(definitions, t, st)

	case reflect.Map:
		r.reflectMap(definitions, t, st)

	case reflect.Interface:
		// empty

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = "integer"

	case reflect.Float32, reflect.Float64:
		s.Type = "number"

	case reflect.Bool:
		s.Type = "boolean"

	case reflect.String:
		s.Type = "string"

	default:
		panic("unsupported type " + t.String())
	}

	return s
}

// reflectStructFields 处理结构体中的字段
func reflectStructFields(s *Schema, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		reflectStructField(s, f)
	}
}

// reflectStructField 处理结构体字段
func reflectStructField(s *Schema, f reflect.StructField) {
	required = strings.Contains(f.Tag, "required")
	jsonTagString, _ := f.Tag.Lookup("json")

}

func appendUniqueString(base []string, value string) []string {
	for _, v := range base {
		if v == value {
			return base
		}
	}
	return append(base, value)
}
