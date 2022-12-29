package goparser

import (
	"go/ast"
	"go/token"
)

func newPackages() *Packages {
	return &Packages{
		files: make(map[string]*AstFile),
		types: make(map[string]*AstTypeSpec),
		pkgs:  make(map[string]*Package),
	}
}

// Packages 管理扫描到的所有包和类型
type Packages struct {
	files map[string]*AstFile     // 使用到的所有go代码文件，key=absPath
	pkgs  map[string]*Package     // 使用到的所有go包
	types map[string]*AstTypeSpec // 使用到的所有类型，key=类型唯一名称（包名+类型名 type.Id）
}

// ParseFile 解析go代码文件中的类型
func (p *Packages) ParseFile(pkgId, absPath string, file *ast.File) *AstFile {
	pkg := p.getPkg(pkgId) // 获取文件所在包
	astFile := newAstFile(pkg, absPath, file)

	p.files[astFile.absPath] = astFile

	for _, decl := range astFile.file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				// 循环获取代码中定义的类型申明
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					astTypeSpec := newAstTypeSpec(astFile, typeSpec)

					p.types[astTypeSpec.Id()] = astTypeSpec
				}
			}
		}
	}
	return astFile
}

func (p *Packages) getPkg(pkgId string) *Package {
	pkg, ok := p.pkgs[pkgId]
	if !ok {
		pkg = newPackage(pkgId)
		p.pkgs[pkg.id] = pkg
	}
	return pkg
}
