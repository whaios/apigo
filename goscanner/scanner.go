// Package goscanner 扫描指定目录中的go代码，获取指定类型。
//
// 使用示例：
// scanner := New("../example/petshop/model")
// err := scanner.Scan()
//
//	err := scanner.WalkFile(func(file *AstFile) error {
//			...
//			astType, err = GetType("Pet", file)
//			...
//	})
package goscanner

import (
	"fmt"
	"github.com/whaios/apigo/log"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

// New go代码解析器
func New() *Scanner {
	return &Scanner{
		mode:     parser.ParseComments,
		files:    make([]*AstFile, 0),
		packages: newPackages(),
	}
}

// Scanner go代码解析器
type Scanner struct {
	rootDir string // 要解析的代码根目录
	rootPkg string // 代码根目录对应的go包名

	mode     parser.Mode
	files    []*AstFile // 收集目录中的go文件，按字母顺序排序
	packages *Packages  // 管理扫描到的所有包和类型
}

// Scan 指定要扫描的代码目录，并开始收集代码，解析类型。
func (p *Scanner) Scan(dir string) error {
	log.UpdateSpinner("获取Go包名")

	rootPkg, err := dirToPkg(dir)
	if err != nil {
		return fmt.Errorf("获取Go包名失败, dir: %s, error: %+v", dir, err.Error())
	}
	p.rootDir = dir
	p.rootPkg = rootPkg
	log.Info("获取Go包名 %s", p.rootPkg)

	// 按字母顺序遍历根目录的文件树（包括树中的每个目录和文件）。
	// 如果rootDir使用相对路径，那么 WalkDirFunc 回调函数中的 id 参数也为相对路径。
	err = filepath.WalkDir(p.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过目录
		if d.IsDir() {
			return nil
		}
		// 跳过非go业务代码文件
		if filepath.Ext(path) != ".go" ||
			strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 根据根包名，计算go文件所在的包
		relPath, err := filepath.Rel(p.rootDir, path)
		if err != nil {
			return err
		}
		pkgId := filepath.ToSlash(filepath.Dir(filepath.Join(p.rootPkg, relPath)))

		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, p.mode)
		if err != nil {
			return err
		}

		astFileInfo := p.packages.ParseFile(pkgId, path, astFile)
		p.files = append(p.files, astFileInfo)
		return nil
	})
	return err
}

// FileCount 获取扫描到的文件个数
func (p *Scanner) FileCount() int {
	return len(p.files)
}

// Files 获取扫描到的文件
func (p *Scanner) Files() []*AstFile {
	return p.files
}

// GetFile 获取指定文件
func (p *Scanner) GetFile(fileName string) *AstFile {
	for _, f := range p.files {
		if f.Name() == fileName {
			return f
		}
	}
	return nil
}

// GetType 获取类型，没有找到返回 nil。
//
//	name: 类型名称，如： User 或 model.User
//	astFile: 用到该类型的 go 代码文件
func (p *Scanner) GetType(name string, astFile *AstFile) (*AstTypeSpec, error) {
	var pkgName, typeName string
	{
		typeName = name
		if parts := strings.SplitN(name, ".", 2); len(parts) == 2 {
			pkgName = parts[0]
			typeName = parts[1]
		}
	}

	// 从以下包中查找类型
	targetPkgs := make([]string, 0)
	if pkgName != "" {
		// 有包名，引用的是其他包中的类型
		imptPkg := astFile.GetImportPkg(pkgName)
		// 没有找到对应的包名，无法解析该类型，则直接返回
		if imptPkg == "" {
			return nil, nil
		}
		targetPkgs = append(targetPkgs, imptPkg)
	} else {
		// 没有包名，引用的是本包中的类型
		pkgId := astFile.PkgId()
		targetPkgs = append(targetPkgs, pkgId)

		// 本包中没有找到该类型，从 . 包中查找
		for _, dotPkg := range astFile.DotImports() {
			targetPkgs = append(targetPkgs, dotPkg)
		}
	}

	// 循环以下包，查找指定类型
	for _, pkgId := range targetPkgs {
		pkg := p.packages.GetPkg(pkgId)
		// 还没有解析过此包
		if pkg == nil {
			// 加载指定包下的代码文件，并解析其中的类型
			if err := p.parseExternalPackage(pkgId); err != nil {
				return nil, err
			}
			// 解析后还是没有此包信息，可能此包下没有代码文件
			if pkg = p.packages.GetPkg(pkgId); pkg == nil {
				continue
			}
		}

		// 找到了该类型直接返回
		if astType := pkg.GetTypeByName(typeName); astType != nil {
			return astType, nil
		}
	}
	return nil, nil
}

// parseExternalPackage 加载指定包下的代码文件，并解析其中的类型
func (p *Scanner) parseExternalPackage(pkg string) error {
	cfg := &packages.Config{
		Dir:  p.rootDir,
		Mode: packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedCompiledGoFiles,
	}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for i, astFile := range pkg.Syntax {
			p.packages.ParseFile(pkg.ID, pkg.CompiledGoFiles[i], astFile)
		}
	}
	return nil
}

// DirToPkg 通过 go list 命令，获取指定目录的包名："../example/goparser" > "goparser"。
//
//	注意：指定目录下必须有 go 代码文件，否则会返回 "no Go files" 错误
func dirToPkg(dir string) (string, error) {
	var stdout, stderr strings.Builder
	var cmd = exec.Command("go", "list", "-f={{.ImportPath}}")
	{
		cmd.Dir = dir
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list command, %s, stdout:%s, stderr:%s", err, stdout.String(), stderr.String())
	}

	outStr := stdout.String()
	{
		if outStr[0] == '_' { // will shown like _/{GOPATH}/src/{YOUR_PACKAGE} when NOT enable GO MODULE.
			outStr = strings.TrimPrefix(outStr, "_"+build.Default.GOPATH+"/src/")
		}
		f := strings.Split(outStr, "\n")
		outStr = f[0]
	}
	return outStr, nil
}
