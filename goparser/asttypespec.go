package goparser

import "go/ast"

func newAstTypeSpec(astFile *AstFile, typeSpec *ast.TypeSpec) *AstTypeSpec {
	astTypeSpec := &AstTypeSpec{
		file:     astFile,
		pkg:      astFile.pkg,
		typeSpec: typeSpec,
	}
	astTypeSpec.pkg.addType(astTypeSpec)
	return astTypeSpec
}

// AstTypeSpec go 类型申明
type AstTypeSpec struct {
	file     *AstFile
	pkg      *Package
	typeSpec *ast.TypeSpec
}

// AbsPath 类型所在文件路径
func (t *AstTypeSpec) AbsPath() string {
	return t.file.absPath
}

// Id 完整包名.类型名
func (t *AstTypeSpec) Id() string {
	return GetTypeId(t.PkgId(), t.Name())
}

// PkgId 完整包名
func (t *AstTypeSpec) PkgId() string {
	return t.pkg.id
}

// Name 类型名称，如：Book
func (t *AstTypeSpec) Name() string {
	return t.typeSpec.Name.Name
}

// GetTypeId 组合类型唯一名称：完整包名.类型名
func GetTypeId(pkgId, typeName string) string {
	return pkgId + "." + typeName
}
