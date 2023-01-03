package apidoc

import (
	"errors"
	"fmt"
	"github.com/whaios/apigo/goparser"
	"github.com/whaios/apigo/log"
	"go/ast"
	"path"
	"regexp"
	"strings"
)

// API 接口注释支持的标签
const (
	TagTitle       = "@title"       // 接口名称，默认取方法名称
	TagCatalog     = "@catalog"     // 接口目录，多级目录用 / 隔开
	TagUrl         = "@url"         // 接口URL，格式为：[method] [url]
	TagDesc        = "@desc"        // 可选，接口说明
	Tagdescription = "@description" // 可选，接口说明
	TagHeader      = "@header"      // 可选，请求 Header 参数。格式为 [字段名] [类型] [必填] ["示例值"] ["说明"]
	TagPathVar     = "@path_var"    // 可选，请求 Path 参数。格式为 [字段名] [类型] [必填] ["示例值"] ["说明"]
	TagQuery       = "@query"       // 可选，请求 Query 参数。支持结构体（如：Struct{}，一对大括号结尾） 或 单个参数（如：[字段名] [类型] [必填] ["值"] ["备注"]）两种方式。
	TagParamMode   = "@param_mode"  // 可选，请求 Body 参数方式。none、form-data(formdata)、x-www-form-urlencoded(urlencoded) 和 json
	TagParam       = "@param"       // 可选，请求 Body 参数。支持结构体（如：Struct{}，一对大括号结尾） 或 单个参数（如：[字段名] [类型] [必填] ["值"] ["备注"]）两种方式。
	TagResp        = "@resp"        // 可选，返回内容。支持结构体（如：Struct{}，一对大括号结尾） 或 单个参数（如：[字段名] [类型] ["备注"]）两种方式。
	TagResponse    = "@response"    // 可选，返回内容。支持结构体（如：Struct{}，一对大括号结尾） 或 单个参数（如：[字段名] [类型] ["备注"]）两种方式。
	TagRemark      = "@remark"      // 可选，备注信息
)

// Generate 扫描指定目录下的 go 代码文件，解析 go 注释，生成 API 接口文档
func Generate(searchDir string) (docs []*ApiDoc, err error) {
	docs = make([]*ApiDoc, 0)

	// 扫描 go 代码
	if err = goparser.Parse(searchDir); err != nil {
		return
	}
	// 循环解析每个 go 代码文件中的注释
	if err = goparser.WalkFile(func(file *goparser.AstFile) error {

		if ds, err := parseGoFile(file); err != nil {
			return err
		} else {
			docs = append(docs, ds...)
		}
		return nil

	}); err != nil {
		return
	}
	return
}

// parseGoFile 解析 go 代码文件中的注释
func parseGoFile(file *goparser.AstFile) ([]*ApiDoc, error) {
	docs := make([]*ApiDoc, 0)
	commdoc := NewApiDoc(nil)
	order := 1
	for _, astDescription := range file.File().Decls {
		switch astDescription.(type) {
		case *ast.GenDecl:
			// 解析类型上的通用注释
			astDecl := astDescription.(*ast.GenDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析通用注释: %s", file.FullName())
				for _, comment := range astDecl.Doc.List {
					if err := parseGoComment(commdoc, file, "", comment.Text); err != nil {
						return docs, fmt.Errorf("解析通用注释出错 %s :%+v", file.FullName(), err)
					}
				}
			}
		case *ast.FuncDecl:
			// 解析方法上的注释
			astDecl := astDescription.(*ast.FuncDecl)
			if astDecl.Doc != nil && astDecl.Doc.List != nil {
				log.Debug("解析方法注释: %s %s()", file.FullName(), astDecl.Name.Name)
				doc := NewApiDoc(commdoc)
				// 逐行解析方法上的注释块
				for _, comment := range astDecl.Doc.List {
					log.Debug("	> 注释: %s", comment.Text)
					if err := parseGoComment(doc, file, astDecl.Name.Name, comment.Text); err != nil {
						return docs, fmt.Errorf("解析方法注释出错 %s %s():%+v", file.FullName(), astDecl.Name.Name, err)
					}
				}
				// 检查是否合法的API文档
				if doc.Invalid() {
					log.Debug("忽略方法注释（没有 title 或 url）: %s()", astDecl.Name.Name)
					continue
				}

				log.Info("生成接口文档(%d) %s", order, doc.Name())
				docs = append(docs, doc)
				order++
			}
		}
	}
	return docs, nil
}

// parseGoComment 解析单行注释
func parseGoComment(doc *ApiDoc, file *goparser.AstFile, funcName, commentLine string) error {
	// 移除注释开头的 // 和空格
	comment := strings.TrimSpace(strings.TrimLeft(commentLine, "/"))
	if comment == "" {
		// 没有注释内容
		return nil
	}

	funcName = strings.ToLower(funcName)
	tagName := strings.ToLower(strings.Fields(comment)[0])
	lineRemainder := strings.TrimSpace(comment[len(tagName):])

	var err error
	switch tagName {
	case funcName, TagTitle:
		doc.Title = lineRemainder
	case TagCatalog:
		doc.Catalog = path.Join(doc.Catalog, lineRemainder)
	case TagUrl:
		err = parseUrlComment(doc, lineRemainder)
	case TagDesc, Tagdescription:
		if doc.Description != "" && lineRemainder != "" {
			doc.Description += "\n"
		}
		doc.Description += lineRemainder
	case TagHeader:
		err = parseHeaderComment(doc, lineRemainder)
	case TagPathVar:
		err = parsePathVarComment(doc, lineRemainder)
	case TagQuery:
		err = parseQueryComment(doc, lineRemainder, file)
	case TagParamMode:
		err = parseParamModeComment(doc, lineRemainder)
	case TagParam:
		err = parseParamComment(doc, lineRemainder, file)
	case TagResp, TagResponse:
		err = parseRespComment(doc, lineRemainder, file)
	case TagRemark:
		if doc.Remark != "" && lineRemainder != "" {
			doc.Remark += "\n"
		}
		doc.Remark += lineRemainder
	}

	return err
}

// parseUrlComment 解析URL，格式为：[method] [url]
//
// 如：GET /pet/{petId}
func parseUrlComment(doc *ApiDoc, comment string) error {
	fields := strings.Fields(comment)
	if len(fields) != 2 {
		return fmt.Errorf("无法解析 url 注释 \"%s\"", comment)
	}
	doc.Method = strings.ToLower(fields[0]) // 和runapi保持一致使用小写
	doc.Path = fields[1]

	// POST请求默认Body类型为JSON
	if doc.RequestBody.Type == "" && doc.Method == MethodPost {
		doc.RequestBody.Type = RequestBodyTypeJSON
	}
	return nil
}

// parseHeaderComment 解析Header。
//
// 如：	Authorization	string	true	"bearer {{TOKEN}}"	"用户登录凭证"
//		[字段名]			[类型]	[必填]	[值]				[说明]
func parseHeaderComment(doc *ApiDoc, comment string) error {
	param, err := stringToParameter(comment)
	if err != nil {
		return fmt.Errorf("无法解析 header 注释 \"%s\"\n%s", comment, err.Error())
	}
	doc.Parameters.Header = append(doc.Parameters.Header, param)
	return nil
}

// parsePathVarComment 解析 Path 参数
//
// 如：	petId	int		true	"1"		"宠物 id"
//		[字段名]	[类型]	[必填]	[值]	[说明]
func parsePathVarComment(doc *ApiDoc, comment string) error {
	param, err := stringToParameter(comment)
	if err != nil {
		return fmt.Errorf("无法解析 path_var 注释 \"%s\"\n%s", comment, err.Error())
	}
	doc.Parameters.Path = append(doc.Parameters.Path, param)
	return nil
}

// parseQueryComment 解析 Query 参数
//
// 方式一：	page		int		true	"1"		"第几页"
//			[字段名]		[类型]	[必填]	[值]	[备注]
//
// 方式二： model.Pet{}
func parseQueryComment(doc *ApiDoc, comment string, file *goparser.AstFile) error {
	// 方式二 结构体的解析
	if strings.HasSuffix(comment, "{}") {
		// TODO 结构体的解析逻辑
		return nil
	}

	// 方式一的解析
	param, err := stringToParameter(comment)
	if err != nil {
		return fmt.Errorf("无法解析 query 注释 \"%s\"\n%s", comment, err.Error())
	}
	doc.Parameters.Query = append(doc.Parameters.Query, param)
	return nil
}

// parseParamModeComment 解析 Body 参数方式。
//	- none
//	- json、application/json
//	- formdata、form-data、multipart/form-data
//	- urlencoded、x-www-form-urlencoded、application/x-www-form-urlencoded
func parseParamModeComment(doc *ApiDoc, comment string) error {
	comment = strings.ToLower(comment)
	switch comment {
	case RequestBodyTypeNone:
		doc.RequestBody.Type = RequestBodyTypeNone
	case "json", RequestBodyTypeJSON:
		doc.RequestBody.Type = RequestBodyTypeJSON
	case "form-data", "formdata", RequestBodyTypeFormData:
		doc.RequestBody.Type = RequestBodyTypeFormData
	case "x-www-form-urlencoded", "urlencoded", RequestBodyTypeUrlEncoded:
		doc.RequestBody.Type = RequestBodyTypeUrlEncoded
	default:
		return fmt.Errorf("不支持 %s 请求参数模式", comment)
	}
	return nil
}

// parseParamComment 解析 Param 参数
//
// 方式一：	page		int		true	"1"		"第几页"
//			[字段名]		[类型]	[必填]	[值]	[备注]
//
// 方式二： model.Pet{}
func parseParamComment(doc *ApiDoc, comment string, file *goparser.AstFile) error {
	// 方式二 结构体的解析
	if strings.HasSuffix(comment, "{}") {
		// TODO 结构体的解析逻辑
		return nil
	}

	// 方式一的解析
	param, err := stringToParameter(comment)
	if err != nil {
		return fmt.Errorf("无法解析 param 注释 \"%s\"\n%s", comment, err.Error())
	}
	doc.RequestBody.Parameters = append(doc.RequestBody.Parameters, param)
	return nil
}

// parseRespComment 解析返回样例
//
// 方式一：	page		int		"第几页"
//			[字段名]		[类型]	[说明]
//
// 方式二： model.Pet{}
func parseRespComment(doc *ApiDoc, comment string, file *goparser.AstFile) error {
	// 方式二 结构体的解析
	if strings.HasSuffix(comment, "{}") {
		// TODO 结构体的解析逻辑
		return nil
	}

	// TODO 方式一的解析
	_, err := stringToRespParameter(comment)
	if err != nil {
		return fmt.Errorf("无法解析 response 注释 \"%s\"\n%s", comment, err.Error())
	}
	doc.Responses = append(doc.Responses, NewResponse("成功"))
	return nil
}

var reqParamPattern = regexp.MustCompile(`(\S+)[\s]+([\w]+)[\s]+([\w]+)[\s]+"([^"]*)"[\s]+"([^"]*)"`)

// stringToParameter 解析参数，格式如： [字段名] [类型] [必填] ["值"] ["说明"]
func stringToParameter(str string) (*Parameter, error) {
	matches := reqParamPattern.FindStringSubmatch(str)
	if len(matches) != 6 {
		return nil, errors.New(`不符合格式 [字段名] [类型] [必填] ["值"] ["备注"]`)
	}
	param := NewParameter(matches[1], matches[2], matches[3], matches[4], matches[5])
	return param, nil
}

var respParamPattern = regexp.MustCompile(`(\S+)[\s]+([\w]+)[\s]+"([^"]*)"`)

// stringToRespParameter 解析返回参数，格式如： [字段名] [类型] ["说明"]
func stringToRespParameter(str string) (*Parameter, error) {
	matches := respParamPattern.FindStringSubmatch(str)
	if len(matches) != 4 {
		return nil, errors.New(`不符合格式 [字段名] [类型] ["备注"]`)
	}
	param := NewRespParameter(matches[1], matches[2], matches[3])
	return param, nil
}
