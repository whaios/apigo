package apifox_test

import (
	"encoding/json"
	"fmt"
	"github.com/whaios/apigo/apifox"
	"github.com/whaios/apigo/parser"
	"os"
	"path/filepath"
)

func ExamplePostImportData() {
	apifox.ProjectId = "0"
	apifox.AccessToken = "a"

	openApiJsonData, err := os.ReadFile(filepath.Join("../example", "openapi_data.json"))
	if err != nil {
		panic(err)
	}
	result, err := apifox.PostImportData(string(openApiJsonData))
	if err != nil {
		panic(err)
	}

	resultData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultData))
	// Output:
	// {"success":false,"data":{"apiCollection":{"item":{"createCount":0,"updateCount":0,"errorCount":0,"ignoreCount":0},"folder":{"createCount":0,"updateCount":0,"errorCount":0,"ignoreCount":0}},"schemaCollection":{"item":{"createCount":0,"updateCount":0,"errorCount":0,"ignoreCount":0},"folder":{"createCount":0,"updateCount":0,"errorCount":0,"ignoreCount":0}}}}
}

func ExampleOpenApi2() {
	goParser := parser.NewParser()
	_, err := goParser.Scan("../example/petshop/pet")
	if err != nil {
		panic(err)
	}
	items, err := goParser.Parse()
	if err != nil {
		panic(err)
	}

	api2 := apifox.NewOpenApi2()
	apifox.OpenApi2AddPaths(api2, items)

	api2JsonData, err := json.Marshal(api2)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(api2JsonData))
	// Output:
	// {"swagger":"2.0","info":{"description":"解析 Go 代码文件中的注释生成 Api 文档。","title":"Apigo","version":"1.0.0"},"paths":{"/pet":{"put":{"consumes":["application/json"],"produces":["application/json"],"summary":"修改宠物信息","parameters":[{"type":"string","example":"bearer {{TOKEN}}","description":"用户登录凭证","name":"Authorization","in":"header","required":true},{"name":"petshop/model.Pet","in":"body","schema":{"type":"object","required":["id","name","status"],"properties":{"category":{"description":"分组","type":"object","properties":{"id":{"description":"分组ID编号","type":"string"},"name":{"description":"分组名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Category"},"id":{"description":"宠物ID编号","type":"string"},"name":{"description":"名称","type":"string"},"photoUrls":{"description":"照片URL","type":"array","items":{"type":"string"}},"status":{"description":"宠物销售状态","type":"string","apigo-type-full-name":"petshop/model.Status"},"tags":{"description":"标签","type":"array","items":{"type":"object","properties":{"id":{"description":"标签ID编号","type":"string"},"name":{"description":"标签名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Tag"}}},"apigo-properties-orders":["category","id","name","photoUrls","status","tags"],"apigo-type-full-name":"petshop/model.Pet"}}],"responses":{"200":{"description":"成功","schema":{"allOf":[{"type":"object","properties":{"errcode":{"description":"错误代码","type":"integer"},"errmsg":{"description":"错误说明","type":"string"}},"apigo-composed-field-key":"data","apigo-properties-orders":["errcode","errmsg"],"apigo-type-full-name":"petshop/comm.HttpCode"},{"type":"object","properties":{"data":{"type":"object","required":["id","name","status"],"properties":{"category":{"description":"分组","type":"object","properties":{"id":{"description":"分组ID编号","type":"string"},"name":{"description":"分组名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Category"},"id":{"description":"宠物ID编号","type":"string"},"name":{"description":"名称","type":"string"},"photoUrls":{"description":"照片URL","type":"array","items":{"type":"string"}},"status":{"description":"宠物销售状态","type":"string","apigo-type-full-name":"petshop/model.Status"},"tags":{"description":"标签","type":"array","items":{"type":"object","properties":{"id":{"description":"标签ID编号","type":"string"},"name":{"description":"标签名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Tag"}}},"apigo-properties-orders":["category","id","name","photoUrls","status","tags"],"apigo-type-full-name":"petshop/model.Pet"}}}]}}},"x-apifox-folder":"宠物商城/宠物管理","x-apifox-status":""},"post":{"consumes":["application/x-www-form-urlencoded"],"produces":["application/json"],"summary":"新建宠物信息","parameters":[{"type":"string","example":"bearer {{TOKEN}}","description":"用户登录凭证","name":"Authorization","in":"header","required":true},{"type":"string","example":"Hello Kitty","description":"宠物名","name":"name","in":"formData","required":true},{"type":"string","example":"sold","description":"宠物销售状态","name":"status","in":"formData","required":true}],"responses":{"200":{"description":"成功示例","schema":{"allOf":[{"type":"object","properties":{"errcode":{"description":"错误代码","type":"integer"},"errmsg":{"description":"错误说明","type":"string"}},"apigo-composed-field-key":"data","apigo-properties-orders":["errcode","errmsg"],"apigo-type-full-name":"petshop/comm.HttpCode"},{"type":"object","properties":{"data":{"type":"object","required":["id","name","status"],"properties":{"category":{"description":"分组","type":"object","properties":{"id":{"description":"分组ID编号","type":"string"},"name":{"description":"分组名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Category"},"id":{"description":"宠物ID编号","type":"string"},"name":{"description":"名称","type":"string"},"photoUrls":{"description":"照片URL","type":"array","items":{"type":"string"}},"status":{"description":"宠物销售状态","type":"string","apigo-type-full-name":"petshop/model.Status"},"tags":{"description":"标签","type":"array","items":{"type":"object","properties":{"id":{"description":"标签ID编号","type":"string"},"name":{"description":"标签名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Tag"}}},"apigo-properties-orders":["category","id","name","photoUrls","status","tags"],"apigo-type-full-name":"petshop/model.Pet"}}}]}}},"x-apifox-folder":"宠物商城/宠物管理","x-apifox-status":""}},"/pet/findByStatus":{"get":{"consumes":["none"],"produces":["application/json"],"summary":"根据状态查找宠物列表","parameters":[{"type":"string","example":"","description":"宠物销售状态","name":"status","in":"query","required":true},{"type":"string","example":"bearer {{TOKEN}}","description":"用户登录凭证","name":"Authorization","in":"header","required":true}],"responses":{"200":{"description":"成功","schema":{"allOf":[{"type":"object","properties":{"errcode":{"description":"错误代码","type":"integer"},"errmsg":{"description":"错误说明","type":"string"}},"apigo-composed-field-key":"data","apigo-properties-orders":["errcode","errmsg"],"apigo-type-full-name":"petshop/comm.HttpCode"},{"type":"object","properties":{"data":{"type":"object","properties":{"pets":{"description":"宠物列表","type":"array","items":{"type":"object","required":["id","name","status"],"properties":{"category":{"description":"分组","type":"object","properties":{"id":{"description":"分组ID编号","type":"string"},"name":{"description":"分组名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Category"},"id":{"description":"宠物ID编号","type":"string"},"name":{"description":"名称","type":"string"},"photoUrls":{"description":"照片URL","type":"array","items":{"type":"string"}},"status":{"description":"宠物销售状态","type":"string","apigo-type-full-name":"petshop/model.Status"},"tags":{"description":"标签","type":"array","items":{"type":"object","properties":{"id":{"description":"标签ID编号","type":"string"},"name":{"description":"标签名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Tag"}}},"apigo-properties-orders":["category","id","name","photoUrls","status","tags"],"apigo-type-full-name":"petshop/model.Pet"}}},"apigo-properties-orders":["pets"],"apigo-type-full-name":"petshop/pet.FindByStatusRsp"}}}]}}},"x-apifox-folder":"宠物商城/宠物管理","x-apifox-status":""}},"/pet/{petId}":{"get":{"description":"指定id查询宠物详情","consumes":["none"],"produces":["application/json"],"summary":"查询宠物详情","parameters":[{"type":"int","example":"1","description":"宠物 id","name":"petId","in":"path","required":true},{"type":"string","example":"bearer {{TOKEN}}","description":"用户登录凭证","name":"Authorization","in":"header","required":true}],"responses":{"200":{"description":"成功","schema":{"allOf":[{"type":"object","properties":{"errcode":{"description":"错误代码","type":"integer"},"errmsg":{"description":"错误说明","type":"string"}},"apigo-composed-field-key":"data","apigo-properties-orders":["errcode","errmsg"],"apigo-type-full-name":"petshop/comm.HttpCode"},{"type":"object","properties":{"data":{"type":"object","required":["id","name","status"],"properties":{"category":{"description":"分组","type":"object","properties":{"id":{"description":"分组ID编号","type":"string"},"name":{"description":"分组名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Category"},"id":{"description":"宠物ID编号","type":"string"},"name":{"description":"名称","type":"string"},"photoUrls":{"description":"照片URL","type":"array","items":{"type":"string"}},"status":{"description":"宠物销售状态","type":"string","apigo-type-full-name":"petshop/model.Status"},"tags":{"description":"标签","type":"array","items":{"type":"object","properties":{"id":{"description":"标签ID编号","type":"string"},"name":{"description":"标签名称","type":"string"}},"apigo-properties-orders":["id","name"],"apigo-type-full-name":"petshop/model.Tag"}}},"apigo-properties-orders":["category","id","name","photoUrls","status","tags"],"apigo-type-full-name":"petshop/model.Pet"}}}]}}},"x-apifox-folder":"宠物商城/宠物管理","x-apifox-status":"developing"},"delete":{"consumes":["application/json"],"produces":["application/json"],"summary":"删除宠物信息","parameters":[{"type":"string","example":"bearer {{TOKEN}}","description":"用户登录凭证","name":"Authorization","in":"header","required":true},{"name":"petshop/pet.DelPetReq","in":"body","schema":{"type":"object","properties":{"pet_id":{"description":"要删除的宠物 id","type":"string"}},"apigo-properties-orders":["pet_id"],"apigo-type-full-name":"petshop/pet.DelPetReq"}}],"responses":{"200":{"description":"组装响应类型","schema":{"type":"object","properties":{"errcode":{"description":"错误代码","type":"integer"},"errmsg":{"description":"错误说明","type":"string"}},"apigo-composed-field-key":"data","apigo-properties-orders":["errcode","errmsg"],"apigo-type-full-name":"petshop/comm.HttpCode"}}},"x-apifox-folder":"宠物商城/宠物管理","x-apifox-status":""}}}}
}
