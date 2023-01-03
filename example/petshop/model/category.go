package model

// Category 宠物分组
type Category struct {
	ID   *int64  `json:"id,omitempty"`   // 分组ID编号
	Name *string `json:"name,omitempty"` // 分组名称
}
