package goparser

import (
	"fmt"
	"github.com/whaios/apigo/log"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Packages 存储扫描到的 Go 文件、包路径和他们之间的关系。
// 主要用于解析注释中可能会引用到的类型。
type Packages struct {
	// 目标Go项目所在目录（支持绝对路径和相对路径），用于加载外部包时获取go包名。
	// 为空时默认为当前运行目录。
	// 如："../example/ginweb/handler"
	dir     string
	pkgName string // dir 目录对应的包名
	mode    parser.Mode

	astFiles   []*AstFileInfo             // 收集目录中的go文件，按字母顺序排序
	astFileMap map[*ast.File]*AstFileInfo // 收集目录中和引用到的go文件
	pkgMap     map[string]*Package        // key=完整包名. 如：ginweb/book
	// 在目录下收集到的具有唯一名称（完整包名+类型名）的类型。如果存在同名的不会出现到该字典中。
	uniqueDefMap map[string]*TypeSpecDef // key=类型全名. 如：ginweb.handler.book.Book
}

func (p *Packages) collectFile() error {
	err := filepath.Walk(p.dir, func(path string, f fs.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), "_test.go") || filepath.Ext(path) != ".go" {
			return nil
		}

		// path = "../example/ginweb/handler/book/handler.go"
		// pkgPath = "ginweb/handler/book"
		var pkgPath string
		{
			relPath, err := filepath.Rel(p.dir, path) // "book/handler.go"
			if err != nil {
				return err
			}
			pkgPath = filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(p.pkgName, relPath))))
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, p.mode)
		if err != nil {
			return fmt.Errorf("parser.ParseFile has error:%+v", err)
		}

		// absPath = "/Users/whai/gowork/github.com/whaios/apigo/example/ginweb/handler/book/handler.go"
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		fileInfo := p.addFile(pkgPath, absPath, astFile)
		p.astFiles = append(p.astFiles, fileInfo)
		return nil
	})
	if err != nil {
		return err
	}

	// 对代码文件进行排序
	sort.Slice(p.astFiles, func(i, j int) bool {
		return strings.Compare(p.astFiles[i].FileName, p.astFiles[j].FileName) < 0
	})
	return nil
}

// 添加 Go 源码文件，并解析代码中的结构体
func (p *Packages) addFile(pkgPath, fileName string, astFile *ast.File) *AstFileInfo {
	info := &AstFileInfo{
		File:     astFile,
		FileName: fileName,
		PkgPath:  pkgPath,
	}
	p.astFileMap[astFile] = info
	log.Debug("收集Go文件: %s", info.FileName)

	// 解析go文件中的结构体
	p.collectTypes(info.File, info.PkgPath)
	return info
}

// 解析go文件中的结构体
func (p *Packages) collectTypes(astFile *ast.File, pkgPath string) {
	for _, astDeclaration := range astFile.Decls {
		if generalDeclaration, ok := astDeclaration.(*ast.GenDecl); ok && generalDeclaration.Tok == token.TYPE {
			for _, astSpec := range generalDeclaration.Specs {
				if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
					typeSpecDef := &TypeSpecDef{
						PkgPath:  pkgPath,
						File:     astFile,
						TypeSpec: typeSpec,
					}

					fullName := typeSpecDef.FullName()
					log.Debug("收集类型: %s", fullName)
					if anotherTypeDef, ok := p.uniqueDefMap[fullName]; ok {
						if typeSpecDef.PkgPath == anotherTypeDef.PkgPath {
							continue
						} else {
							delete(p.uniqueDefMap, fullName)
						}
					} else {
						p.uniqueDefMap[fullName] = typeSpecDef
					}

					if p.pkgMap[typeSpecDef.PkgPath] == nil {
						p.pkgMap[typeSpecDef.PkgPath] = &Package{
							Name:            astFile.Name.Name,
							TypeDefinitions: map[string]*TypeSpecDef{typeSpecDef.Name(): typeSpecDef},
						}
					} else if _, ok = p.pkgMap[typeSpecDef.PkgPath].TypeDefinitions[typeSpecDef.Name()]; !ok {
						p.pkgMap[typeSpecDef.PkgPath].TypeDefinitions[typeSpecDef.Name()] = typeSpecDef
					}
				}
			}
		}
	}
}

// FindTypeSpec 查找类型
//
// @param shortName 包名.类型名，如：ListRsp 或 book.Book
func (p *Packages) FindTypeSpec(shortName string, file *ast.File) *TypeSpecDef {
	if file == nil { // for test
		return p.uniqueDefMap[shortName]
	}

	var pkgName, typeName string
	{
		typeName = shortName
		if parts := strings.SplitN(shortName, ".", 2); len(parts) == 2 {
			pkgName = parts[0]
			typeName = parts[1]
		}
	}

	// 有包名，查找外部包
	if pkgName != "" {
		// 从文件中导入的包中查找指定包路径
		imptPkgPath, _ := p.findPackagePathFromImports(pkgName, file)
		// 没有找到对应的包名
		if imptPkgPath == "" {
			return nil
		}

		// 收集外部包
		p.loadExternalPackage(imptPkgPath)
		return p.findTypeSpec(imptPkgPath, typeName)
	}

	var pkgPath, fullName string
	if fileInfo, ok := p.astFileMap[file]; ok {
		pkgPath = fileInfo.PkgPath
		fullName = pkgPath + "." + typeName
	}

	// 从目录包中查找
	typeDef, ok := p.uniqueDefMap[fullName]
	if ok {
		return typeDef
	}

	typeDef = p.findTypeSpec(pkgPath, typeName)
	if typeDef != nil {
		return typeDef
	}

	// 载入 . 包
	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == "." {
			imptPkgPath := strings.Trim(imp.Path.Value, `"`)
			// 收集外部包
			p.loadExternalPackage(imptPkgPath)

			if typeDef = p.findTypeSpec(imptPkgPath, typeName); typeDef != nil {
				return typeDef
			}
		}
	}
	return nil
}

// findTypeSpec 从收集的指定包中查找类型
//
// @param pkgPath 如：refstruct/employee
// @param typeName 如：Employee
func (p *Packages) findTypeSpec(pkgPath string, typeName string) *TypeSpecDef {
	if p.pkgMap == nil {
		return nil
	}
	pd, found := p.pkgMap[pkgPath]
	if found {
		typeSpec, ok := pd.TypeDefinitions[typeName]
		if ok {
			return typeSpec
		}
	}
	return nil
}

// imptPkgPath 加载指定包
func (p *Packages) loadExternalPackage(imptPkgPath string) {
	if p.pkgMap != nil {
		if _, ok := p.pkgMap[imptPkgPath]; ok {
			// 已经收集过该包
			return
		}
	}
	log.Debug("加载外部包: %s", imptPkgPath)

	cfg := &packages.Config{
		Dir:  p.dir,
		Mode: packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedCompiledGoFiles,
	}
	pkgs, _ := packages.Load(cfg, imptPkgPath)
	for _, pkg := range pkgs {
		for i, astFile := range pkg.Syntax {
			p.addFile(pkg.ID, pkg.CompiledGoFiles[i], astFile)
		}
	}
}

// 从文件中导入的包中查找指定包路径。
// @pkgName 包名
func (p *Packages) findPackagePathFromImports(pkgName string, file *ast.File) (pkgPath string, isAliasPkgName bool) {
	for _, imp := range file.Imports {
		// 有别名的包，别名相同，直接取得该包路径
		if imp.Name != nil && imp.Name.Name == pkgName {
			pkgPath = strings.Trim(imp.Path.Value, `"`)
			isAliasPkgName = true
			break
		}

		// 普通导入，包没有别名
		path := strings.Trim(imp.Path.Value, `"`)
		paths := strings.Split(path, "/")
		if paths[len(paths)-1] == pkgName {
			// 找到包路径
			pkgPath = path
			break
		}
	}
	return
}

// AstFileInfo ast.File 文件信息.
type AstFileInfo struct {
	File     *ast.File
	FileName string // Go 源码文件全名称
	PkgPath  string // Go 源码文件完整包名
}
