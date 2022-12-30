package pet

import _ "petshop/comm"

// Handler 宠物管理
//
// @catalog 宠物商城/宠物管理
// @header Authorization string true "bearer {{TOKEN}}" "用户登录凭证"
// @resp comm.HttpCode{}
type Handler struct {
}

// GetPet 查询宠物详情
//
// @desc 指定id查询宠物详情
// @url GET /pet/{petId}
// @path_var petId int true "" "宠物 id"
// @resp Pet{}
func (h *Handler) GetPet() {

}

// CreatePet 新建宠物信息
//
// @url POST /pet
// @param_mode urlencoded
// @param name string true "Hello Kitty" "宠物名"
// @param status string true "sold" "宠物销售状态"
// @resp Pet{}
func (h *Handler) CreatePet() {

}

// EditPet 修改宠物信息
//
// @url PUT /pet
// @query Pet{}
// @resp Pet{}
func (h *Handler) EditPet() {

}

// DelPet 删除宠物信息
//
// @url DELETE /pet/{petId}
// @path_var petId int true "" "要删除的宠物 id"
func (h *Handler) DelPet() {

}

// FindByStatus 根据状态查找宠物列表
//
// @url GET /pet/findByStatus
// @query status string true "" "宠物销售状态"
// @resp List{}
func (h *Handler) FindByStatus() {

}
