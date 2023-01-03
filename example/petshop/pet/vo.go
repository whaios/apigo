package pet

import "petshop/model"

type FindByStatusRsp struct {
	Pets []*model.Pet `json:"pets"` // 宠物列表
}
