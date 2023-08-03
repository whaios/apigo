package goscanner

import (
	"github.com/whaios/apigo/log"
	"go/ast"
	"go/token"
)

func newPackages() *Packages {
	return &Packages{
		files: make(map[string]*AstFile),
		pkgs:  make(map[string]*Package),
	}
}

// Packages 管理扫描到的所有包和文件
type Packages struct {
	files map[string]*AstFile // 使用到的所有go代码文件，key=absPath
	pkgs  map[string]*Package // 使用到的所有go包，key=pkgId
}

// ParseFile 解析go代码文件中的类型
func (p *Packages) ParseFile(pkgId, path string, file *ast.File) *AstFile {
	pkg, ok := p.pkgs[pkgId]
	if !ok {
		pkg = newPackage(pkgId)
		p.pkgs[pkg.id] = pkg
	}

	astFile := pkg.AddFile(path, file)
	p.files[astFile.absPath] = astFile

	for _, decl := range astFile.file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				// 循环获取代码中定义的类型申明
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					tpe := pkg.AddType(astFile, typeSpec)
					log.Debug("	> 采集类型: %s", tpe.Name())
				}
			}
		}
	}
	return astFile
}

// GetPkg 获取解析过的指定包，没有找到返回nil
func (p *Packages) GetPkg(pkgId string) *Package {
	return p.pkgs[pkgId]
}
