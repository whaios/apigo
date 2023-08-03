package recursive

import "goparser/samename"

// TypeRecursive 递归类型
type TypeRecursive struct {
	A       string                    `json:"a"`        // a
	Childs  []*TypeRecursive          `json:"childs"`   // 递归类型数组
	MChilds map[string]*TypeRecursive `json:"m_childs"` // 递归类型字典
}

type SameName struct {
	samename.SameName        // 测试同名包+同名结构体
	Desc              string `json:"desc"` // 介绍
}
