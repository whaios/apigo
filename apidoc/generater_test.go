package apidoc

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseGoComment(t *testing.T) {
	Convey("测试解析go代码中的注释", t, func() {
		funcName := "TestMethod"
		comments := []string{
			`// TestMethod 示例接口`,
			`// `,
			`// @catalog	一级目录/二级目录`,
			`// @url		GET /pet/{petId}`,
			`// @desc		接口说明`,
			`// @header		Authorization	string	true	"bearer {{TOKEN}}"	"用户登录凭证"`,
			`// @path_var	petId			int		true	"1"					"宠物 id"`,
			`// @query		usr				int		true	"1"					"当前用户 id"`,
			`// @param_mode	json`,
			`// @param		name			string	true	"Hello Kitty"		"宠物名称"`,
			`// @resp		name			string	"宠物名称"`,
			`// @remark		接口备注`,
		}

		doc := NewApiDoc(nil)
		for _, comment := range comments {
			err := parseGoComment(doc, nil, funcName, comment)
			So(err, ShouldBeNil)
		}

		wantJson := `{"title":"示例接口","catalog":"一级目录/二级目录","description":"接口说明","remark":"接口备注","method":"get","path":"/pet/{petId}","parameters":{"path":[{"name":"petId","type":"int","required":true,"example":"1","description":"宠物 id"}],"query":[{"name":"usr","type":"int","required":true,"example":"1","description":"当前用户 id"}],"header":[{"name":"Authorization","type":"string","required":true,"example":"bearer {{TOKEN}}","description":"用户登录凭证"}]},"requestBody":{"type":"application/json","parameters":[{"name":"name","type":"string","required":true,"example":"Hello Kitty","description":"宠物名称"}]},"responses":[{"name":"成功"}]}`
		gotJsonData, err := json.Marshal(doc)
		So(err, ShouldBeNil)
		gotJson := string(gotJsonData)
		So(gotJson, ShouldEqual, wantJson)
	})
}
