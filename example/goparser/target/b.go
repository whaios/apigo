package target

import (
	cmm "goparser/common"
)

type B struct {
	BFieldStr string
	BFieldInt int

	AliasPkgName cmm.Cmm // 测试包别名
}
