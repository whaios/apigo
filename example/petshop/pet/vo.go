package pet

import "petshop/model"

type FindByStatusRsp struct {
	Pets []*model.Pet `json:"pets"` // 宠物列表
}

type DelPetReq struct {
	PetId int64 `json:"pet_id,string"` // 要删除的宠物 id
}