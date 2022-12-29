package goparser

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

// Parser go代码解析器
type Parser struct {
	rootDir string      // 要解析的代码根目录
	rootPkg string      // 代码根目录对应的go包名
	mode    parser.Mode // parser.ParseComments

	files    []*AstFile // 收集目录中的go文件，按字母顺序排序
	packages *Packages  // 管理扫描到的所有包和类型
}

// Parse 指定要解析的代码目录，并开始收集代码
func (p *Parser) Parse(dir string) error {
	p.mode = parser.ParseComments
	p.rootDir = dir

	rootPkg, err := GetPkgName(p.rootDir)
	if err != nil {
		return fmt.Errorf("获取包名失败, dir: %s, error: %+v", p.rootDir, err.Error())
	}
	p.rootPkg = rootPkg

	if p.packages == nil {
		p.packages = newPackages()
	}
	return p.collectFile()
}

// collectFile 开始收集指定根目录下的go代码
func (p *Parser) collectFile() error {
	// 按字母顺序遍历根目录的文件树（包括树中的每个目录和文件）。
	// 如果rootDir使用相对路径，那么 WalkDirFunc 回调函数中的 path 参数也为相对路径。
	return filepath.WalkDir(p.rootDir, func(path string, d fs.DirEntry, err error) error {
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
		pkgId := filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(p.rootPkg, relPath))))

		// 获取文件的绝对路径名称
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, p.mode)
		if err != nil {
			return err
		}

		astFileInfo := p.packages.ParseFile(pkgId, absPath, astFile)
		p.files = append(p.files, astFileInfo)
		return nil
	})
}

// GetType 获取类型。
// 	name: 类型名称，如： User 或 model.User
// 	astFile: 用到该类型的 go 代码文件
func (p *Parser) GetType(name string, astFile *AstFile) *AstTypeSpec {
	var pkgName, typeName string
	{
		typeName = name
		if parts := strings.SplitN(name, ".", 2); len(parts) == 2 {
			pkgName = parts[0]
			typeName = parts[1]
		}
	}

	// 有包名，指定类型不在当前包
}

// GetPkgName 通过go list命令，获取指定目录的包名："../example/goparser" > "goparser"
func GetPkgName(dir string) (string, error) {
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
