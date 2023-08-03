package parser

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/whaios/apigo/goscanner"
	"testing"
)

func TestParseGoComment(t *testing.T) {
	Convey("测试解析go代码中的注释", t, func() {
		wantApiItem := `{
    "title": "示例接口",
    "folder": "一级目录/二级目录",
    "status": "released",
    "description": "接口说明",
    "remark": "接口备注",
    "method": "get",
    "path": "/pet/{petId}",
    "parameters": {
        "path": [
            {
                "name": "petId",
                "type": "int",
                "required": true,
                "example": "1",
                "description": "宠物 id"
            }
        ],
        "query": [
            {
                "name": "usr",
                "type": "int",
                "required": true,
                "example": "1",
                "description": "当前用户 id"
            }
        ],
        "header": [
            {
                "name": "Authorization",
                "type": "string",
                "required": true,
                "example": "bearer {{TOKEN}}",
                "description": "用户登录凭证"
            }
        ],
        "bodyType": "application/json",
        "formData": [
            {
                "name": "name",
                "type": "string",
                "required": true,
                "example": "Hello Kitty",
                "description": "宠物名称"
            }
        ]
    },
    "content_type": ""
}`

		funcName := "TestMethod"
		comments := []string{
			`// TestMethod 示例接口`,
			`// `,
			`// @folder	一级目录/二级目录`,
			`// @status 	released`,
			`// @desc		接口说明`,
			`// @remark		接口备注`,
			`// @url		GET /pet/{petId}`,
			`// @bodytype	json`,
			`// @param		header	Authorization	string	true	"bearer {{TOKEN}}"	"用户登录凭证"`,
			`// @param		path	petId			int		true	"1"					"宠物 id"`,
			`// @param		query	usr				int		true	"1"					"当前用户 id"`,
			`// @param		form	name			string	true	"Hello Kitty"		"宠物名称"`,
		}

		apiItem := &ApiItem{}
		p := NewParser()
		for _, comment := range comments {
			err := p.parseGoComment(apiItem, nil, funcName, comment)
			So(err, ShouldBeNil)
		}

		gotJsonData, err := json.MarshalIndent(apiItem, "", "    ")
		So(err, ShouldBeNil)
		So(string(gotJsonData), ShouldEqual, wantApiItem)
	})
}

func Test_ParseSchema(t *testing.T) {
	Convey("测试解析结构体", t, func() {
		scanner := goscanner.New()
		err := scanner.Scan("../example/goparser/simple")
		So(err, ShouldBeNil)
		astFile := scanner.GetFile("simple.go")

		tp := NewParser()
		tp.SetScanner(scanner)
		schema, err := tp.ParseType("SomeStruct", astFile)
		So(err, ShouldBeNil)

		wantSchema := `{
    "type": "object",
    "required": [
        "string"
    ],
    "properties": {
        "bool": {
            "description": "布尔",
            "type": "boolean"
        },
        "bytes": {
            "type": "array",
            "items": {
                "type": "integer"
            }
        },
        "float64": {
            "description": "浮点数",
            "type": "number"
        },
        "int": {
            "description": "整数",
            "type": "integer"
        },
        "int64": {
            "description": "大整数",
            "type": "integer"
        },
        "json": {
            "type": "array",
            "items": {
                "type": "integer"
            },
            "apigo-type-full-name": "encoding/json.RawMessage"
        },
        "map": {
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        },
        "ptr": {
            "type": "string"
        },
        "slice": {
            "type": "array",
            "items": {
                "type": "string"
            }
        },
        "string": {
            "description": "字符串",
            "type": "string"
        },
        "struct": {
            "type": "object",
            "required": [
                "x"
            ],
            "properties": {
                "x": {
                    "type": "string"
                }
            },
            "apigo-properties-orders": [
                "x"
            ]
        },
        "structWithoutFields": {
            "type": "object",
            "properties": {
                "Y": {
                    "type": "string"
                }
            },
            "apigo-properties-orders": [
                "Y"
            ]
        },
        "time": {
            "type": "string",
            "format": "date-time"
        }
    },
    "apigo-properties-orders": [
        "bool",
        "int",
        "int64",
        "float64",
        "string",
        "bytes",
        "json",
        "time",
        "slice",
        "map",
        "struct",
        "structWithoutFields",
        "ptr"
    ],
    "apigo-type-full-name": "goparser/simple.SomeStruct"
}`

		data, err := json.MarshalIndent(schema, "", "    ")
		So(err, ShouldBeNil)
		//fmt.Printf("%s\n", data)
		So(string(data), ShouldEqual, wantSchema)
	})
}

func Test_ParseTypeRecursive(t *testing.T) {
	Convey("测试解析递归类型", t, func() {
		scanner := goscanner.New()
		err := scanner.Scan("../example/goparser/recursive")
		So(err, ShouldBeNil)
		astFile := scanner.GetFile("recursive.go")

		tp := NewParser()
		tp.SetScanner(scanner)
		schema, err := tp.ParseType("TypeRecursive", astFile)
		So(err, ShouldBeNil)

		wantSchema := `{
    "type": "object",
    "properties": {
        "a": {
            "description": "a",
            "type": "string"
        },
        "childs": {
            "description": "递归类型数组",
            "type": "array"
        },
        "m_childs": {
            "description": "递归类型字典",
            "type": "object",
            "additionalProperties": true
        }
    },
    "apigo-properties-orders": [
        "a",
        "childs",
        "m_childs"
    ],
    "apigo-type-full-name": "goparser/recursive.TypeRecursive"
}`
		data, err := json.MarshalIndent(schema, "", "    ")
		So(err, ShouldBeNil)
		//fmt.Printf("%s\n", data)
		So(string(data), ShouldEqual, wantSchema)
	})
}

func Test_ParseSameName(t *testing.T) {
	Convey("测试解析同名不同包的类型", t, func() {
		scanner := goscanner.New()
		err := scanner.Scan("../example/goparser/recursive")
		So(err, ShouldBeNil)
		astFile := scanner.GetFile("recursive.go")

		tp := NewParser()
		tp.SetScanner(scanner)
		schema, err := tp.ParseType("SameName", astFile)
		So(err, ShouldBeNil)

		wantSchema := `{
    "type": "object",
    "properties": {
        "desc": {
            "description": "介绍",
            "type": "string"
        },
        "name": {
            "description": "名称",
            "type": "string"
        }
    },
    "apigo-properties-orders": [
        "name",
        "name",
        "desc"
    ],
    "apigo-type-full-name": "goparser/recursive.SameName"
}`
		data, err := json.MarshalIndent(schema, "", "    ")
		So(err, ShouldBeNil)
		//fmt.Printf("%s\n", data)
		So(string(data), ShouldEqual, wantSchema)
	})
}
