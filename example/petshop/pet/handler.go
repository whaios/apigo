package pet

import (
	_ "petshop/comm"
	_ "petshop/model"
)

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
