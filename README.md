# Apigo

Apigo 主要用于解析 Go (Golang) 代码注释，快速生成 API 文档，并同步到 [Apifox](https://apifox.com/)，实现代码零入侵。

## 目录

- [命令说明](#命令说明)
- [安装](#安装)
- [设置环境变量](#设置环境变量)
- [使用示例](#使用示例)
- [注释格式](#注释格式)
  - [API信息](#API信息)
  - [请求参数](#请求参数)
  - [返回响应](#返回响应)

## 命令说明

```
$ apigo.exe
NAME:
   Apigo - 快速生成 Go (Golang) API 文档

USAGE:
   Apigo [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

DESCRIPTION:
   Apigo 主要用于解析 Go (Golang) 代码注释，快速生成 API 文档，并同步到 Apifox，实现代码零入侵。

COMMANDS:
   apifox, af  快速生成 API 文档，并同步到 Apifox。
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        开启调试模式。 (default: false)
   --help         显示帮助 (default: false)
   --version, -v  print the version
```

## 安装

已经编译好的平台有： [点击下载](https://github.com/whaios/apigo/releases)
- windows/amd64
- linux/amd64
- darwin/amd64

## 设置环境变量

### Apifox

生成的文档同步到 Apifox，需要指定 Apifox 的 `个人访问令牌` 和 `项目 ID`。
为了避免每次同步文档时都需要输入这两个变量，建议将其配置到系统环境变量中。

- 个人访问令牌: `ApifoxAccessToken`，查看[如何获取个人访问令牌](https://apifox.com/help/openapi/)
- 项目 ID: `ApifoxProjectId`，打开 Apifox 进入项目里的“项目设置”查看

## 使用示例

### Go 代码

详细代码可参考 `example/petshop/pet/handler.go`

```go
package pet

// Handler 宠物管理
//
// @folder	宠物商城/宠物管理
// @param	header Authorization string true "bearer {{TOKEN}}" "用户登录凭证"
// @resp	200	"组装响应类型"	comm.HttpCode{data}
type Handler struct {
}

// GetPet 查询宠物详情
//
// @desc 	指定id查询宠物详情
// @remark 	本接口需要登录
// @status 	developing
// @url 	GET /pet/{petId}
// @param 	path petId int true "1" "宠物 id"
// @success model.Pet{}
func (h *Handler) GetPet() {
}

// CreatePet 新建宠物信息
//
// @url 		POST /pet
// @bodytype	x-www-form-urlencoded
// @param 		form	name 	string	true	"Hello Kitty" 	"宠物名"
// @param 		form	status 	string	true	"sold" 			"宠物销售状态"
// @resp 		200	"成功示例"	model.Pet{}
func (h *Handler) CreatePet() {
}

// EditPet 修改宠物信息
//
// @url 	PUT /pet
// @param 	body model.Pet{}
// @success model.Pet{}
func (h *Handler) EditPet() {
}

// DelPet 删除宠物信息
//
// @url 	DELETE /pet/{petId}
// @param 	body DelPetReq{}
func (h *Handler) DelPet() {
}

// FindByStatus 根据状态查找宠物列表
//
// @url 	GET /pet/findByStatus
// @param 	query status string true "" "宠物销售状态"
// @success	FindByStatusRsp{}
func (h *Handler) FindByStatus() {
}
```

### 生成文档并同步到 Apifox

```shell
$ apigo.exe apifox --dir ./example/petshop/pet/

扫描目录 ./example/petshop/pet/
获取Go包名 petshop/pet
采集到2个Go代码文件
生成接口文档(1) 宠物商城/宠物管理/查询宠物详情
生成接口文档(2) 宠物商城/宠物管理/新建宠物信息
生成接口文档(3) 宠物商城/宠物管理/修改宠物信息
生成接口文档(4) 宠物商城/宠物管理/删除宠物信息
生成接口文档(5) 宠物商城/宠物管理/根据状态查找宠物列表
新增接口 0，修改接口 5，出错接口 0，忽略接口 0
新增模型 0，修改模型 0，出错模型 0，忽略模型 0
同步 Apifox 成功

```

## 注释格式

### API信息

| 注释      | 说明                               | 示例                       |
|---------|----------------------------------|--------------------------|
| @title  | **必须**，接口名称                      | // @title 查询宠物详情         |
| @folder | **必须**，接口所属目录，多级目录使用斜杠`/`分隔      | // @folder 一级/二级/三级      |
| @url    | **必须**，接口URL，格式：`[method] [url]` | // @url GET /pet/{petId} |
| @status | [接口状态](#Apifox 接口状态)             | // @status released      |
| @desc   | 接口说明                             | // @desc 指定id查询宠物详情      |

### 请求参数

| 注释        | 说明                                                                                                  | 示例                                                                          |
|-----------|-----------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| @param    | 请求参数，用空格分隔。<br>格式一：`参数类型` `参数名` `数据类型` `必填` `"值"` `"备注"` <br>格式二：`参数类型` `结构体`<br>查看支持的[参数类型](#参数类型) | // @param path petId int true "123" "宠物 id" <br> // @param body model.Pet{} |
| @bodytype | Body 类型，查看支持的[Mime类型](#Mime类型)                                                                      | // @bodytype json                                                           |

### 返回响应

| 注释           | 说明                                        | 示例                                                               |
|--------------|-------------------------------------------|------------------------------------------------------------------|
| @contenttype | 响应类型，查看支持的[Mime类型](#Mime类型)               | // @contenttype json                                             |
| @resp        | 响应内容，用空格分隔。<br>格式：`http 状态码` `名称` `结构体{}` | // @resp 200 "成功" model.Pet{}                                    |
| @success     | 成功响应内容                                    | // @success	model.Pet{}<br>等效于：<br>// @resp 200 "成功" model.Pet{} |

### Apifox 接口状态

| 状态  | 	代码          |
|-----|--------------|
| 设计中 | 	designing   |
| 待确定 | 	pending     |
| 开发中 | 	developing  |
| 联调中 | 	integrating |
| 测试中 | 	testing     |
| 已测完 | 	tested      |
| 已发布 | 	released    |
| 已废弃 | 	deprecated  |
| 有异常 | 	exception   |
| 已废弃 | 	obsolete    |
| 将废弃 | 	deprecated  |

### 参数类型

| 类型     | 
|--------|
| path   |
 | query  |
 | header |
 | cookie |
 | form   |
 | body   |

### Mime类型

| 别名                    | 	类型                               |
|-----------------------|-----------------------------------|
| form-data             | multipart/form-data               |
 | x-www-form-urlencoded | application/x-www-form-urlencoded |
 | json                  | application/json                  |
 | xml                   | application/xml                   |
 | html                  | text/html                         |
 | raw                   | text/plain                        |
 | binary                | application/octet-stream          |

### 数据类型
 - integer (byte, uint, int, int32, int64) 
 - number (float32, float64)  
 - boolean (bool) 
 - string (string)  
 - struct

## 参考
- [OpenAPI 规范 (中文版)](https://openapi.apifox.cn/)
- [Go OpenAPI 3.0](https://github.com/getkin/kin-openapi)
- [Go Swagger 2.0](https://github.com/swaggo/swag)