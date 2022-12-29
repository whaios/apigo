package goparser

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

func (p Package) addFile(astFile *AstFile) {
	p.files[astFile.absPath] = astFile
}

func (p Package) addType(astTypeSpec *AstTypeSpec) {
	p.types[astTypeSpec.Id()] = astTypeSpec
}
