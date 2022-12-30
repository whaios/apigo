package goparser

import (
	"go/ast"
	"path/filepath"
	"strings"
)

// AstFile Go 源码文件信息.
type AstFile struct {
	pkg     *Package  // Go 源码文件所在包
	absPath string    // Go 源码文件全名称
	file    *ast.File // Go 源码文件
}

// PkgId 获取 go 源码文件所属包
func (f *AstFile) PkgId() string {
	return f.pkg.id
}

// Name 获取 文件名.后缀
func (f *AstFile) Name() string {
	return filepath.Base(f.absPath)
}

// GetImportPkg 从 import 的包中查找指定包，没有找到返回空字符串
func (f *AstFile) GetImportPkg(pkgName string) string {
	for _, impt := range f.file.Imports {
		path := strings.Trim(impt.Path.Value, `"`)
		// 匹配到有别名的包
		if impt.Name != nil &&
			impt.Name.Name != "_" { // 代码中不会引用到别名为"_"的包，这里为了支持注释里面以这种方式引用类型
			if impt.Name.Name == pkgName {
				return path
			}
			continue
		}

		// 常规导入包
		if name := GetPkgName(path); name == pkgName {
			return path
		}
	}
	return ""
}

// DotImports 获取文件中的 . 包路径
func (f *AstFile) DotImports() []string {
	dotimpts := make([]string, 0)
	for _, impt := range f.file.Imports {
		if impt.Name != nil && impt.Name.Name == "." {
			path := strings.Trim(impt.Path.Value, `"`)
			dotimpts = append(dotimpts, path)
		}
	}
	return dotimpts
}
