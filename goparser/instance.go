package goparser

var instance *Parser

// Parse 指定要解析的代码目录，并开始收集代码，解析类型
func Parse(dir string) (err error) {
	instance, err = NewParser(dir)
	if err != nil {
		return
	}
	return instance.Parse()
}

// WalkFile 循环处理扫描到的每个文件
func WalkFile(fn WalkFileFunc) error {
	if instance == nil {
		return nil
	}
	return instance.WalkFile(fn)
}

// GetType 获取类型，没有找到返回 nil。
// 	name: 类型名称，如： User 或 model.User
// 	astFile: 用到该类型的 go 代码文件
func GetType(name string, astFile *AstFile) (*AstTypeSpec, error) {
	if instance == nil {
		return nil, nil
	}
	return instance.GetType(name, astFile)
}
