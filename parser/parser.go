package parser

import (
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/whaios/apigo/goscanner"
	"github.com/whaios/apigo/log"
	"go/ast"
	"regexp"
	"sort"
	"strings"
)

// 参数类型
const (
	ParamTypePath   = "path"
	ParamTypeQuery  = "query"
	ParamTypeHeader = "header"
	ParamTypeCookie = "cookie"
	ParamTypeForm   = "form"
	ParamTypeBody   = "body"
)

// 数据格式
const (
	MimeAliasFormData       = "form-data"             // multipart/form-data
	MimeAliasFormUrlencoded = "x-www-form-urlencoded" // application/x-www-form-urlencoded
	MimeAliasJson           = "json"                  // application/json
	MimeAliasXml            = "xml"                   // text/xml
	MimeAliasHtml           = "html"                  // text/html
	MimeAliasRaw            = "raw"                   // text/plain
	MimeAliasBinary         = "binary"                // application/octet-stream

	BodyTypeNone           = "none"
	BodyTypeFormData       = "multipart/form-data"
	BodyTypeFormUrlEncoded = "application/x-www-form-urlencoded"
	BodyTypeJSON           = "application/json"
	BodyTypeXML            = "application/xml"
	BodyTypeHTML           = "text/html"
	BodyTypePlain          = "text/plain"
	BodyTypeOctetStream    = "application/octet-stream"
)

// API 接口注释支持的标签
const (
	TagTitle  = "@title"  // 接口名称，默认取方法名称
	TagFolder = "@folder" // 接口目录，多级目录用 / 隔开
	TagStatus = "@status" // 可选，接口状态
	TagDesc   = "@desc"   // 可选，接口说明
	TagRemark = "@remark" // 可选，备注信息

	TagUrl         = "@url"         // 接口URL，格式为：[method] [url]
	TagBodyType    = "@bodytype"    // 可选，Body 类型，仅影响具有请求正文的操作，例如 POST、PUT 和 PATCH。
	TagParam       = "@param"       // 可选，请求参数。支持结构体（如：[参数类型] [Struct{}]，一对大括号结尾） 或 单个参数（如：[参数类型] [参数名] [数据类型] [必填] ["值"] ["备注"]）两种方式。
	TagContentType = "@contenttype" // 可选，响应类型，默认 JSON。
	TagSuccess     = "@success"     // 可选，成功(200)响应，例如 Struct{}。
	TagResp        = "@resp"        // 可选，返回内容。支持结构体（如：[http 状态码] [名称] 结构体{}）。
)

func NewParser() *Parser {
	return &Parser{
		parsedSchemas: make(map[*goscanner.AstTypeSpec]*spec.Schema),
		scanner:       goscanner.New(),
	}
}

type Parser struct {
	parsedSchemas map[*goscanner.AstTypeSpec]*spec.Schema
	scanner       *goscanner.Scanner
}

func (p *Parser) SetScanner(scanner *goscanner.Scanner) {
	p.scanner = scanner
}

// Scan 扫描指定目录中的 go 代码，返回文件个数。
func (p *Parser) Scan(dir string) (int, error) {
	// 扫描 go 代码
	err := p.scanner.Scan(dir)
	return p.scanner.FileCount(), err
}

// Parse 解析 go 代码注释为 API接口文档
func (p *Parser) Parse() ([]ApiItem, error) {
	apiItems := make([]ApiItem, 0)
	// 循环解析每个 go 代码文件中的注释
	for _, file := range p.scanner.Files() {
		log.Debug(log.UpdateSpinner("解析文件 %s", file.Path()))
		if items, err := p.parseGoFile(file); err != nil {
			return apiItems, err
		} else {
			apiItems = append(apiItems, items...)
		}
	}
	return apiItems, nil
}

// parseGoFile 解析 go 代码文件中的注释
func (p *Parser) parseGoFile(file *goscanner.AstFile) ([]ApiItem, error) {
	commItem := &ApiItem{}
	apiItems := make([]ApiItem, 0)
	order := 1
	for _, astDescription := range file.File().Decls {
		switch astDescription.(type) {
		case *ast.GenDecl:
			// 解析类型上的通用注释
			astDecl := astDescription.(*ast.GenDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析通用注释: %s", file.Path())
				for _, comment := range astDecl.Doc.List {
					if err := p.parseGoComment(commItem, file, "", comment.Text); err != nil {
						return nil, fmt.Errorf("解析通用注释出错 %s :%+v", file.Path(), err)
					}
				}
			}
		case *ast.FuncDecl:
			// 解析方法上的注释
			astDecl := astDescription.(*ast.FuncDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析方法注释: %s %s()", file.Path(), astDecl.Name.Name)
				apiItem := &ApiItem{}
				// 逐行解析方法上的注释块
				for _, comment := range astDecl.Doc.List {
					log.Debug("	> 注释: %s", comment.Text)
					if err := p.parseGoComment(apiItem, file, astDecl.Name.Name, comment.Text); err != nil {
						return nil, fmt.Errorf("解析方法注释出错 %s %s():%+v", file.Path(), astDecl.Name.Name, err)
					}
				}
				// 检查是否合法的API文档
				if apiItem.Invalid() {
					log.Debug("忽略方法注释（没有 title 或 url）: %s()", astDecl.Name.Name)
					continue
				}

				apiItem.UseCommon(commItem)
				log.Info("生成接口文档(%d) %s", order, apiItem.Name())
				apiItems = append(apiItems, *apiItem)
				order++
			}
		}
	}
	return apiItems, nil
}

// parseGoComment 解析单行注释
func (p *Parser) parseGoComment(apiItem *ApiItem, file *goscanner.AstFile, funcName, commentLine string) error {
	// 移除注释开头的 // 和空格
	comment := strings.TrimSpace(strings.TrimLeft(commentLine, "//"))
	if comment == "" {
		// 没有注释内容
		return nil
	}

	funcName = strings.ToLower(strings.TrimSpace(funcName))
	tagName := strings.ToLower(strings.TrimSpace(strings.Fields(comment)[0]))
	lineRemainder := strings.TrimSpace(comment[len(tagName):])

	var err error
	switch tagName {
	case funcName, TagTitle:
		apiItem.Title = lineRemainder
	case TagFolder:
		apiItem.AddFolder(lineRemainder)
	case TagStatus:
		apiItem.Status = lineRemainder
	case TagUrl:
		err = p.parseUrlComment(apiItem, lineRemainder)
	case TagDesc:
		if apiItem.Description != "" && lineRemainder != "" {
			apiItem.Description += "\n"
		}
		apiItem.Description += lineRemainder
	case TagBodyType:
		err = p.parseBodyTypeComment(apiItem, lineRemainder)
	case TagParam:
		err = p.parseParamComment(apiItem, lineRemainder, file)
	case TagContentType:
		err = p.parseContentTypeComment(apiItem, lineRemainder)
	case TagSuccess:
		err = p.parseSuccessComment(apiItem, lineRemainder, file)
	case TagResp:
		err = p.parseRespComment(apiItem, lineRemainder, file)
	case TagRemark:
		apiItem.AddRemark(lineRemainder)
	}

	return err
}

// parseUrlComment 解析URL，格式为：[method] [url]
//
// 如：GET /pet/{petId}
func (p *Parser) parseUrlComment(apiItem *ApiItem, comment string) error {
	fields := strings.Fields(comment)
	if len(fields) != 2 {
		return fmt.Errorf("无法解析 url 注释 \"%s\"", comment)
	}
	apiItem.Method = strings.ToLower(strings.TrimSpace(fields[0])) // 使用小写
	apiItem.Path = strings.TrimSpace(fields[1])

	// POST请求默认Body类型为JSON
	if apiItem.Parameters.BodyType == "" &&
		(apiItem.Method == MethodPost || apiItem.Method == MethodPut || apiItem.Method == MethodPatch) {
		apiItem.Parameters.BodyType = BodyTypeJSON
	}
	return nil
}

// parseParamModeComment 解析 Body 参数方式。
//   - none
//   - json、application/json
//   - formdata、form-data、multipart/form-data
//   - urlencoded、x-www-form-urlencoded、application/x-www-form-urlencoded
func (p *Parser) parseBodyTypeComment(apiItem *ApiItem, comment string) error {
	comment = strings.ToLower(comment)
	switch comment {
	case MimeAliasFormData, BodyTypeFormData:
		apiItem.Parameters.BodyType = BodyTypeFormData
	case MimeAliasFormUrlencoded, BodyTypeFormUrlEncoded:
		apiItem.Parameters.BodyType = BodyTypeFormUrlEncoded
	case MimeAliasJson, BodyTypeJSON:
		apiItem.Parameters.BodyType = BodyTypeJSON
	case MimeAliasXml, BodyTypeXML:
		apiItem.Parameters.BodyType = BodyTypeXML
	case MimeAliasHtml, BodyTypeHTML:
		apiItem.Parameters.BodyType = BodyTypeHTML
	case MimeAliasRaw, BodyTypePlain:
		apiItem.Parameters.BodyType = BodyTypePlain
	case MimeAliasBinary, BodyTypeOctetStream:
		apiItem.Parameters.BodyType = BodyTypeOctetStream
	default:
		return fmt.Errorf("不支持 %s 请求参数类型", comment)
	}
	return nil
}

var paramPattern = regexp.MustCompile(`(\w+)\s+(\S+)\s+([\w\-.\\{}=,\[\s\]]+)\s+(\w+)\s+"([^"]*)"\s+"([^"]*)"`)

// parseParamComment 解析参数
//
//	格式：	[字段名]		[参数类型]	[数据类型]		[必填]	[值]	[备注]
//
// @param	page		query		int				true	"1"		"第几页"
// @param				body		model.Pet{}
func (p *Parser) parseParamComment(apiItem *ApiItem, comment string, file *goscanner.AstFile) error {
	parseParamErr := func(errMsg string) error {
		return fmt.Errorf("无法解析 param 注释 \"%s\"\n%s", comment, errMsg)
	}

	var paramType, name, dataType, required, example, desc string

	// 结构体
	if strings.HasSuffix(comment, "{}") {
		if fields := strings.Fields(comment); len(fields) == 2 {
			paramType = strings.TrimSpace(fields[0])
			dataType = strings.TrimSpace(fields[1])
		}
	} else {
		// 属性的解析
		matches := paramPattern.FindStringSubmatch(comment)
		if len(matches) != 7 {
			return parseParamErr("属性格式错误")
		}

		paramType = strings.TrimSpace(matches[1])
		name = strings.TrimSpace(matches[2])
		dataType = strings.TrimSpace(matches[3])
		required = strings.TrimSpace(matches[4])
		example = strings.TrimSpace(matches[5])
		desc = strings.TrimSpace(matches[6])
	}

	parameters := make([]Parameter, 0)
	// 结构体类型
	if strings.HasSuffix(dataType, "{}") {
		objectType := strings.TrimRight(dataType, "{}")
		schema, err := p.ParseType(objectType, file)
		if err != nil {
			return parseParamErr(err.Error())
		}
		if paramType == ParamTypeBody {
			if apiItem.Parameters.BodyType == "" {
				apiItem.Parameters.BodyType = BodyTypeJSON
			}
			apiItem.Parameters.JsonSchema = schema
			return nil
		}

		// 将结构体转为参数数组
		parameters = SchemaToParameters(schema)
	} else {
		param := NewParameter(name, dataType, required, example, desc)
		parameters = append(parameters, param)
	}

	// 属性
	switch paramType {
	case ParamTypePath:
		apiItem.Parameters.Path = append(apiItem.Parameters.Path, parameters...)
	case ParamTypeQuery:
		apiItem.Parameters.Query = append(apiItem.Parameters.Query, parameters...)
	case ParamTypeHeader:
		apiItem.Parameters.Header = append(apiItem.Parameters.Header, parameters...)
	case ParamTypeCookie:
		apiItem.Parameters.Cookie = append(apiItem.Parameters.Cookie, parameters...)
	case ParamTypeForm:
		apiItem.Parameters.FormData = append(apiItem.Parameters.FormData, parameters...)
	}
	return nil
}

// parseContentTypeComment 解析响应内容格式。
//   - json、application/json
//   - text/xml
//   - text/html
//   - text/plain
//   - application/octet-stream
func (p *Parser) parseContentTypeComment(apiItem *ApiItem, comment string) error {
	comment = strings.ToLower(comment)
	switch comment {
	case MimeAliasJson, BodyTypeJSON:
		apiItem.ContentType = BodyTypeJSON
	case MimeAliasXml, BodyTypeXML:
		apiItem.ContentType = BodyTypeXML
	case MimeAliasHtml, BodyTypeHTML:
		apiItem.ContentType = BodyTypeHTML
	case MimeAliasRaw, BodyTypePlain:
		apiItem.ContentType = BodyTypePlain
	case MimeAliasBinary, BodyTypeOctetStream:
		apiItem.ContentType = BodyTypeOctetStream
	default:
		return fmt.Errorf("不支持 %s 响应内容格式", comment)
	}
	return nil
}

var respPattern = regexp.MustCompile(`(\d+)\s+"([^"]+)"\s+([\w\-.\\{}=,\[\s\]]+)`)

// ResponseType{data1=Type1,data2=Type2}.
var combinedPattern = regexp.MustCompile(`^([\w\-./\[\]]+){(.*)}$`)

// parseRespComment 解析返回样例
//
//	格式		[http 状态码] 	[名称] 	结构体{}
//
// @resp	200				"成功"	model.Pet{}
func (p *Parser) parseRespComment(doc *ApiItem, comment string, file *goscanner.AstFile) error {
	matches := respPattern.FindStringSubmatch(comment)
	if len(matches) != 4 {
		return fmt.Errorf("无法解析 response 注释 \"%s\"\n不符合格式", comment)
	}

	return p.parseResp(doc,
		strToInt(matches[1]),
		strings.TrimSpace(matches[2]),
		strings.TrimSpace(matches[3]), file)
}

// parseSuccessComment 解析返回成功
//
// @success	model.Pet{}
func (p *Parser) parseSuccessComment(doc *ApiItem, comment string, file *goscanner.AstFile) error {
	return p.parseResp(doc, 200, "成功", comment, file)
}

func (p *Parser) parseResp(doc *ApiItem, code int, name, dataType string, file *goscanner.AstFile) error {
	matches := combinedPattern.FindStringSubmatch(dataType)
	if len(matches) != 3 {
		return fmt.Errorf("无法解析 response 数据类型 \"%s\"", dataType)
	}

	objectType := matches[1]
	schema, err := p.ParseType(objectType, file)
	if err != nil {
		return err
	}

	// 组装多个类型
	fields, props := parseFields(matches[2]), map[string]spec.Schema{}
	for _, field := range fields {
		keyVal := strings.SplitN(field, "=", 2)
		if len(keyVal) == 1 {
			SchemaSetComposedFieldKey(schema, keyVal[0])
			continue
		}
		if len(keyVal) == 2 {
			subSchema, err := p.ParseType(keyVal[1], file)
			if err != nil {
				return err
			}
			props[keyVal[0]] = *subSchema
		}
	}
	if len(props) > 0 {
		schema = spec.ComposedSchema(*schema, spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:       []string{OBJECT},
				Properties: props,
			},
		})
	}

	resp := &Response{
		Code:       code,
		Name:       name,
		JsonSchema: schema,
	}
	doc.Responses = append(doc.Responses, resp)
	return nil
}

// ParseType 解析指定类型
func (p *Parser) ParseType(typeName string, astFile *goscanner.AstFile) (*spec.Schema, error) {
	log.Debug("解析类型: %s", typeName)

	if isGolangPrimitiveType(typeName) {
		return primitiveSchema(transToValidSchemeType(typeName)), nil
	}

	typeSpecDef, err := p.scanner.GetType(typeName, astFile)
	if err != nil {
		return nil, err
	}
	if typeSpecDef == nil {
		return nil, fmt.Errorf("没有找到类型定义: %s", typeName)
	}

	schema, found := p.parsedSchemas[typeSpecDef]
	if found {
		return schema, nil
	}

	schema, err = p.parseTypeExpr(typeSpecDef.Name(), typeSpecDef.File, typeSpecDef.TypeSpec.Type, false)
	if err != nil {
		return nil, err
	}
	SchemaSetTypeFullName(schema, typeSpecDef.Id())

	p.parsedSchemas[typeSpecDef] = schema
	return schema, nil
}

func (p *Parser) parseTypeExpr(rootTypeName string, file *goscanner.AstFile, typeExpr ast.Expr, ref bool) (*spec.Schema, error) {
	switch expr := typeExpr.(type) {
	case *ast.Ident: // Baz
		return p.ParseType(expr.Name, file)
	case *ast.StarExpr: // *Baz
		return p.parseTypeExpr(rootTypeName, file, expr.X, ref)
	case *ast.InterfaceType: // interface{}
		return &spec.Schema{}, nil
	case *ast.SelectorExpr: // pkg.Bar
		if xIdent, ok := expr.X.(*ast.Ident); ok {
			return p.ParseType(xIdent.Name+"."+expr.Sel.Name, file)
		}
	case *ast.StructType: // struct {...}
		return p.parseStruct(rootTypeName, file, expr)
	case *ast.ArrayType: // []Baz
		itemSchema, err := p.parseTypeExpr(rootTypeName, file, expr.Elt, true)
		if err != nil {
			return nil, err
		}

		return spec.ArrayProperty(itemSchema), nil
	case *ast.MapType: // map[string]Bar
		if _, ok := expr.Value.(*ast.InterfaceType); ok {
			return spec.MapProperty(nil), nil
		}
		schema, err := p.parseTypeExpr(rootTypeName, file, expr.Value, true)
		if err != nil {
			return nil, err
		}

		return spec.MapProperty(schema), nil
	}

	return &spec.Schema{SchemaProps: spec.SchemaProps{Type: []string{"object"}}}, nil
}

func (p *Parser) parseStruct(rootTypeName string, file *goscanner.AstFile, st *ast.StructType) (*spec.Schema, error) {
	required, orders, properties := make([]string, 0), make([]string, 0), make(map[string]spec.Schema)
	for _, field := range st.Fields.List {
		dataType := parseFieldType(field.Type)
		// 匿名字段
		if len(field.Names) == 0 {
			nSchema, err := p.ParseType(dataType, file)
			if err != nil {
				return nil, err
			}
			if len(nSchema.Type) > 0 && nSchema.Type[0] == OBJECT && len(nSchema.Properties) > 0 {
				for k, v := range nSchema.Properties {
					properties[k] = v
					orders = append(orders, k)
				}
			}
			required = append(required, nSchema.SchemaProps.Required...)
			orders = append(orders, SchemaGetPropertiesOrders(nSchema)...)
			// 没有解析为有效类型，忽略该字段
			continue
		}

		fieldName := field.Names[0].Name
		if field.Tag != nil {
			// `json:"name" validate:"required"`
			tag := field.Tag.Value
			if jsonTag := getJsonTag(tag); jsonTag != "" {
				jsonName, tagOpts := parseJsonTag(jsonTag)
				if tagOpts.Contains("string") {
					// json 标签中定义了类型转换
					dataType = "string"
				}
				if jsonName == "-" {
					continue
				}
				if jsonName != "" {
					fieldName = jsonName
				}
			}
			if strings.Contains(tag, "required") {
				required = append(required, fieldName)
			}
		}

		comment := ""
		if field.Doc != nil {
			// 字段上行的注释
			// 忽略这种注释
		}
		if field.Comment != nil {
			// 字段后面的同行注释
			for _, comm := range field.Comment.List {
				comment += strings.TrimSpace(strings.TrimLeft(comm.Text, "//"))
			}
		}

		var fschema *spec.Schema
		switch {
		case dataType == "" || dataType == INTERFACE:
			fschema = primitiveSchema(OBJECT)
		case isGolangPrimitiveType(dataType):
			fschema = primitiveSchema(transToValidSchemeType(dataType))
		case dataType == TIME:
			fschema = primitiveSchema(STRING)
			fschema.Format = "date-time"
		case dataType == STRUCT:
			var err error
			fschema, err = p.parseStruct(rootTypeName, file, field.Type.(*ast.StructType))
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(dataType, "[]"):
			itemType := dataType[2:]
			if itemType == rootTypeName {
				// 避免无限循环解析递归类型
				fschema = spec.ArrayProperty(nil)
			} else {
				item, err := p.ParseType(itemType, file)
				if err != nil {
					return nil, err
				}
				fschema = spec.ArrayProperty(item)
			}
		case strings.HasPrefix(dataType, "map["):
			// ignore key type
			idx := strings.Index(dataType, "]")
			if idx < 0 {
				return nil, fmt.Errorf("invalid type: %s", dataType)
			}
			dataType = dataType[idx+1:]
			if dataType == rootTypeName ||
				dataType == INTERFACE || dataType == ANY {
				fschema = spec.MapProperty(nil)
			} else {
				item, err := p.ParseType(dataType, file)
				if err != nil {
					return nil, err
				}
				fschema = spec.MapProperty(item)
			}
		default:
			var err error
			if fschema, err = p.ParseType(dataType, file); err != nil {
				return nil, err
			}
		}

		fschema.WithDescription(comment)
		properties[fieldName] = *fschema
		orders = append(orders, fieldName)
	}

	sort.Strings(required)
	schema := &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       []string{OBJECT},
			Properties: properties,
			Required:   required,
		},
	}
	SchemaSetPropertiesOrders(schema, orders)
	return schema, nil
}

func parseFieldType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch expr.(type) {
	case *ast.Ident:
		id := expr.(*ast.Ident)
		if id.Obj != nil && id.Obj.Decl != nil {
			if ts, ok := id.Obj.Decl.(*ast.TypeSpec); ok {
				// 自定义类型（可能是基础类型，也可能是struct，struct返回空字符串）
				if _, ok = ts.Type.(*ast.Ident); ok {
					return parseFieldType(ts.Type)
				}
			}
		}
		return id.Name
	case *ast.ArrayType:
		arrt := expr.(*ast.ArrayType)
		return "[]" + parseFieldType(arrt.Elt)
	case *ast.MapType:
		mpt := expr.(*ast.MapType)
		kn := parseFieldType(mpt.Key)
		vn := parseFieldType(mpt.Value)
		return fmt.Sprintf("map[%s]%s", kn, vn)
	case *ast.SelectorExpr: // 包名.类型
		selt := expr.(*ast.SelectorExpr)
		pkgName := parseFieldType(selt.X)
		if pkgName == "" {
			return selt.Sel.Name
		}
		return fmt.Sprintf("%s.%s", pkgName, selt.Sel.Name)
	case *ast.StarExpr: // 指针
		star := expr.(*ast.StarExpr)
		return parseFieldType(star.X)
	case *ast.InterfaceType:
		return INTERFACE
	case *ast.StructType: // 内部类
		return STRUCT
	}
	return ""
}

func isGolangPrimitiveType(typeName string) bool {
	switch typeName {
	case "uint",
		"int",
		"uint8",
		"int8",
		"uint16",
		"int16",
		"byte",
		"uint32",
		"int32",
		"rune",
		"uint64",
		"int64",
		"float32",
		"float64",
		"bool",
		"string":
		return true
	}
	return false
}

const (
	TIME   = "time.Time"
	STRUCT = "struct{}"
)

const (
	// ARRAY represent a array value.
	ARRAY = "array"
	// OBJECT represent a object value.
	OBJECT = "object"
	// PRIMITIVE represent a primitive value.
	PRIMITIVE = "primitive"
	// BOOLEAN represent a boolean value.
	BOOLEAN = "boolean"
	// INTEGER represent a integer value.
	INTEGER = "integer"
	// NUMBER represent a number value.
	NUMBER = "number"
	// STRING represent a string value.
	STRING = "string"
	// FUNC represent a function value.
	FUNC = "func"
	// ERROR represent a error value.
	ERROR = "error"
	// INTERFACE represent a interface value.
	INTERFACE = "interface{}"
	// ANY represent a any value.
	ANY = "any"
	// NIL represent a empty value.
	NIL = "nil"

	// IgnoreNameOverridePrefix Prepend to model to avoid renaming based on comment.
	IgnoreNameOverridePrefix = '$'
)

func transToValidSchemeType(typeName string) string {
	switch typeName {
	case "uint", "int", "uint8", "int8", "uint16", "int16", "byte":
		return INTEGER
	case "uint32", "int32", "rune":
		return INTEGER
	case "uint64", "int64":
		return INTEGER
	case "float32", "float64":
		return NUMBER
	case "bool":
		return BOOLEAN
	case "string":
		return STRING
	}
	return typeName
}

func primitiveSchema(refType string) *spec.Schema {
	return &spec.Schema{SchemaProps: spec.SchemaProps{Type: []string{refType}}}
}

func parseFields(s string) []string {
	nestLevel := 0

	return strings.FieldsFunc(s, func(char rune) bool {
		if char == '{' {
			nestLevel++

			return false
		} else if char == '}' {
			nestLevel--

			return false
		}

		return char == ',' && nestLevel == 0
	})
}

const (
	schemaExtraTypeFullName     = "apigo-type-full-name"
	schemaExtraComposedFieldKey = "apigo-composed-field-key"
)

var schemaExtraPropertiesOrders = "apigo-properties-orders"

// SetSchemaExtraPropertiesOrdersKey 设置属性排序存储的KEY
func SetSchemaExtraPropertiesOrdersKey(key string) {
	schemaExtraPropertiesOrders = key
}

func SchemaSetTypeFullName(schema *spec.Schema, val string) {
	if schema.ExtraProps == nil {
		schema.ExtraProps = make(map[string]interface{})
	}
	schema.ExtraProps[schemaExtraTypeFullName] = val
}

func SchemaGetTypeFullName(schema *spec.Schema) string {
	if schema.ExtraProps != nil {
		if val, ok := schema.ExtraProps[schemaExtraTypeFullName]; ok {
			return val.(string)
		}
	}
	return ""
}

func SchemaSetComposedFieldKey(schema *spec.Schema, val string) {
	if schema.ExtraProps == nil {
		schema.ExtraProps = make(map[string]interface{})
	}
	schema.ExtraProps[schemaExtraComposedFieldKey] = val
}

func SchemaGetComposedFieldKey(schema *spec.Schema) string {
	if schema.ExtraProps != nil {
		if val, ok := schema.ExtraProps[schemaExtraComposedFieldKey]; ok {
			return val.(string)
		}
	}
	return ""
}

func SchemaSetPropertiesOrders(schema *spec.Schema, val []string) {
	if schema.ExtraProps == nil {
		schema.ExtraProps = make(map[string]interface{})
	}
	schema.ExtraProps[schemaExtraPropertiesOrders] = val
}

func SchemaGetPropertiesOrders(schema *spec.Schema) []string {
	if schema.ExtraProps != nil {
		if val, ok := schema.ExtraProps[schemaExtraPropertiesOrders]; ok {
			return val.([]string)
		}
	}
	return make([]string, 0)
}
