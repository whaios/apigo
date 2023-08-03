package goscanner

import (
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestDirToPkg(t *testing.T) {
	Convey("测试获取指定目录的go包名", t, func() {
		goparserRelPath := "../example/goparser/target"
		goparserAbsPath, _ := filepath.Abs(goparserRelPath)

		cases := []struct {
			Dir string
			Pkg string
		}{
			{goparserRelPath, "goparser/target"},
			{goparserAbsPath, "goparser/target"},
			{"../example/goparser/target/pkga", "goparser/target/pkga"},
		}

		for _, cs := range cases {
			gotPkg, err := dirToPkg(cs.Dir)
			So(err, ShouldBeNil)
			So(gotPkg, ShouldEqual, cs.Pkg)
		}
	})
}

func TestParser_Parse(t *testing.T) {
	Convey("测试解析go代码", t, func() {
		p := New()
		err := p.Scan("../example/goparser/target")
		So(err, ShouldBeNil)

		So(p.FileCount() > 0, ShouldBeTrue)
		// a.go 文件排在第一
		So(p.files[0].Name(), ShouldEqual, "a.go")
		// z.go 文件排在最后
		So(p.files[p.FileCount()-1].Name(), ShouldEqual, "z.go")

		for _, file := range p.files {
			// 测试获取 a.go 文件中使用到的类型
			if file.Name() == "a.go" {
				var astType *AstTypeSpec
				// ---- 获取不存在的类型 --------------------
				{
					astType, err = p.GetType("INVALID", file)
					So(err, ShouldBeNil)
					So(astType, ShouldBeNil)
				}
				// ---- 获取本文件中定义的类型 --------------------
				{
					astType, err = p.GetType("A", file)
					So(err, ShouldBeNil)
					So(astType, ShouldNotBeNil)
					So(astType.Id(), ShouldEqual, "goparser/target.A")
				}
				// ---- 获取本包中定义的类型 --------------------
				{
					astType, err = p.GetType("B", file)
					So(err, ShouldBeNil)
					So(astType, ShouldNotBeNil)
					So(astType.Id(), ShouldEqual, "goparser/target.B")
				}
				// ---- 获取外部包中定义的类型 --------------------
				{
					astType, err = p.GetType("pkga.PkgA", file)
					So(err, ShouldBeNil)
					So(astType, ShouldNotBeNil)
					So(astType.Id(), ShouldEqual, "goparser/target/pkga.PkgA")
				}
				// ---- 获取 . 包中定义的类型 --------------------
				{
					astType, err = p.GetType("DotPkg", file)
					So(err, ShouldBeNil)
					So(astType, ShouldNotBeNil)
					So(astType.Id(), ShouldEqual, "goparser/dotpkg.DotPkg")
				}
				// ---- 获取 _ 包中定义的类型 --------------------
				{
					astType, err = p.GetType("unusepkg.UnUsePkg", file)
					So(err, ShouldBeNil)
					So(astType, ShouldNotBeNil)
					So(astType.Id(), ShouldEqual, "goparser/unusepkg.UnUsePkg")
				}
			}
		}
	})
}
