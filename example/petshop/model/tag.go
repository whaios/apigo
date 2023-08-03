package model

// Tag 宠物标签
type Tag struct {
	ID   *int64  `json:"id,string,omitempty"` // 标签ID编号
	Name *string `json:"name,omitempty"`      // 标签名称
}
