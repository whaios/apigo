package goscanner

import (
	"go/ast"
	"strings"
)

func newPackage(id string) *Package {
	return &Package{
		id:    id,
		files: make(map[string]*AstFile),
		types: make(map[string]*AstTypeSpec),
	}
}

// Package 管理 go 包下面的文件和类型
type Package struct {
	id    string                  // 包完整名称
	files map[string]*AstFile     // 包下的go代码文件，key=absPath
	types map[string]*AstTypeSpec // 使用到的所有类型，key=类型唯一名称（包名+类型名 type.Id）
}

func (p *Package) AddFile(path string, file *ast.File) *AstFile {
	astFile := NewAstFile(p, path, file)
	p.files[astFile.AbsPath()] = astFile
	return astFile
}

func (p *Package) AddType(astFile *AstFile, typeSpec *ast.TypeSpec) *AstTypeSpec {
	astTypeSpec := &AstTypeSpec{
		pkg:      p,
		File:     astFile,
		TypeSpec: typeSpec,
	}
	p.types[astTypeSpec.Id()] = astTypeSpec
	return astTypeSpec
}

// GetType 获取类型
func (p *Package) GetType(typeId string) *AstTypeSpec {
	return p.types[typeId]
}

// GetTypeByName 获取类型
func (p *Package) GetTypeByName(typeName string) *AstTypeSpec {
	typeId := TypeId(p.id, typeName)
	return p.GetType(typeId)
}

// GetPkgName 根据包路径截取出包名
func GetPkgName(pkgId string) string {
	paths := strings.Split(pkgId, "/")
	return paths[len(paths)-1]
}
