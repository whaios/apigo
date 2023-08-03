package target

import (
	. "goparser/dotpkg"
	"goparser/target/pkga"
	_ "goparser/unusepkg"
)

type A struct {
	AFieldStr string
	AFieldInt int

	pkga.PkgA
	DotPkg
}
