package model

// Pet 宠物资料
type Pet struct {
	Category  Category `json:"category"`                      // 分组
	Id        int64    `json:"id,string" validate:"required"` // 宠物ID编号
	Name      string   `json:"name" validate:"required"`      // 名称
	PhotoUrls []string `json:"photoUrls"`                     // 照片URL
	Status    Status   `json:"status" validate:"required"`    // 宠物销售状态
	Tags      []Tag    `json:"tags"`                          // 标签
}
