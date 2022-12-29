package goparser

import "go/ast"

func newAstFile(pkg *Package, absPath string, file *ast.File) *AstFile {
	astFile := &AstFile{
		pkg:     pkg,
		absPath: absPath,
		file:    file,
	}
	pkg.addFile(astFile)
	return astFile
}

// AstFile Go 源码文件信息.
type AstFile struct {
	pkg     *Package  // Go 源码文件所在包
	absPath string    // Go 源码文件全名称
	file    *ast.File // Go 源码文件
}
