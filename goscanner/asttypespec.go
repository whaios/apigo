package goscanner

import "go/ast"

// AstTypeSpec go 类型申明
type AstTypeSpec struct {
	pkg      *Package
	File     *AstFile
	TypeSpec *ast.TypeSpec
}

// PkgId 完整包名
func (t *AstTypeSpec) PkgId() string {
	return t.pkg.id
}

// AbsPath 类型所在文件路径
func (t *AstTypeSpec) AbsPath() string {
	return t.File.absPath
}

// Id 完整包名.类型名
func (t *AstTypeSpec) Id() string {
	return TypeId(t.PkgId(), t.Name())
}

// Name 类型名称，如：Book
func (t *AstTypeSpec) Name() string {
	return t.TypeSpec.Name.Name
}

// TypeId 组合类型唯一名称：完整包名.类型名
func TypeId(pkgId, typeName string) string {
	return pkgId + "." + typeName
}
