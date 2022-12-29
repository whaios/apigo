package goparser

import (
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestGetPkgName(t *testing.T) {
	Convey("测试获取指定目录的go包名", t, func() {
		goparserRelPath := "../example/goparser"
		goparserAbsPath, _ := filepath.Abs(goparserRelPath)
		cases := []struct {
			Dir string
			Pkg string
		}{
			{goparserRelPath, "goparser"},
			{goparserAbsPath, "goparser"},
			{"../example/goparser/util", "goparser/util"},
		}

		for _, cs := range cases {
			gotPkg, err := GetPkgName(cs.Dir)
			So(err, ShouldBeNil)
			So(gotPkg, ShouldEqual, cs.Pkg)
		}
	})
}

func TestParser_Parse(t *testing.T) {
	Convey("测试解析go代码", t, func() {
		// absPath := `D:\Work\github.com\whaios\apigo\example\goparser`
		searchDir, _ := filepath.Abs("../example/goparser")

		p := new(Parser)
		So(p.Parse(searchDir), ShouldBeNil)
		So(len(p.files) > 0, ShouldBeTrue)
	})
}
