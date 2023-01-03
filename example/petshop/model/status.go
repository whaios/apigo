package model

// Status 宠物销售状态
type Status string

const (
	Available Status = "available" // 可售
	Pending   Status = "pending"   // 待售
	Sold      Status = "sold"      // 已售
)
