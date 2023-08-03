package target

import (
	bb "goparser/target/pkgb"
)

type B struct {
	BFieldStr string
	BFieldInt int

	AliasPkgName bb.PkgB // 测试包别名
}
