package goparser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"os/exec"
	"strings"
)

// ParseDir 收集指定目录下的所有Go代码文件。
func ParseDir(dir string, mode parser.Mode) (*Packages, error) {
	pkgName, err := GetPkgName(dir)
	if err != nil {
		return nil, fmt.Errorf("获取包名失败, dir: %s, error: %s", dir, err.Error())
	}

	pkgs := &Packages{
		dir:          dir,
		pkgName:      pkgName,
		mode:         mode,
		astFileMap:   make(map[*ast.File]*AstFileInfo),
		pkgMap:       make(map[string]*Package),
		uniqueDefMap: make(map[string]*TypeSpecDef),
	}
	return pkgs, pkgs.collectFile()
}

// GetPkgName 通过go list命令，获取指定目录的包名："./example/ginweb/handler" > "ginweb/handler"
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

func IsPrimitiveType(typeName string) bool {
	switch typeName {
	case "uint",
		"int",
		"uint8",
		"int8",
		"uint16",
		"int16",
		"byte",
		"uint32",
		"int32",
		"rune",
		"uint64",
		"int64",
		"float32",
		"float64",
		"bool",
		"string":
		return true
	}
	return false
}

func ParseFieldType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch expr.(type) {
	case *ast.Ident:
		id := expr.(*ast.Ident)
		if id.Obj != nil && id.Obj.Decl != nil {
			if ts, ok := id.Obj.Decl.(*ast.TypeSpec); ok {
				// 自定义类型（可能是基础类型，也可能是struct，struct返回空字符串）
				if _, ok = ts.Type.(*ast.Ident); ok {
					return ParseFieldType(ts.Type)
				}
			}
		}
		return id.Name
	case *ast.ArrayType:
		arrt := expr.(*ast.ArrayType)
		return "[]" + ParseFieldType(arrt.Elt)
	case *ast.MapType:
		mpt := expr.(*ast.MapType)
		kn := ParseFieldType(mpt.Key)
		vn := ParseFieldType(mpt.Value)
		return fmt.Sprintf("map[%s]%s", kn, vn)
	case *ast.SelectorExpr: // 包名.类型
		selt := expr.(*ast.SelectorExpr)
		pkgName := ParseFieldType(selt.X)
		if pkgName == "" {
			return selt.Sel.Name
		}
		return fmt.Sprintf("%s.%s", pkgName, selt.Sel.Name)
	case *ast.StarExpr: // 指针
		star := expr.(*ast.StarExpr)
		return ParseFieldType(star.X)
	case *ast.InterfaceType:
		return "interface{}"
	}
	return ""
}
