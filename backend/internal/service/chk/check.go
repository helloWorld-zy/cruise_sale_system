//go:build ignore
// +build ignore

// check 是一个开发辅助工具，用于扫描 Go 源码中缺少注释的导出符号。
// 遍历指定包的 AST，检查导出的函数、类型、变量和常量是否具有文档注释。
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// main 解析上级目录中的 Go 源文件，检查所有导出符号是否缺少注释。
func main() {
	fset := token.NewFileSet()
	// 解析上级目录（service 包）中的所有 Go 源文件
	pkgs, err := parser.ParseDir(fset, "..", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 遍历所有包及其文件
	for pkgName, pkg := range pkgs {
		for filename, f := range pkg.Files {
			// 遍历 AST 节点，查找缺少注释的导出符号
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:
					// 检查导出函数是否缺少文档注释
					if x.Name.IsExported() && x.Doc == nil {
						fmt.Printf("%s: func %s lacks comment (pkg %s)\n", filename, x.Name.Name, pkgName)
					}
				case *ast.GenDecl:
					// 检查导出的类型、变量和常量是否缺少文档注释
					if x.Tok == token.TYPE || x.Tok == token.VAR || x.Tok == token.CONST {
						for _, spec := range x.Specs {
							switch s := spec.(type) {
							case *ast.TypeSpec:
								if s.Name.IsExported() {
									if x.Doc == nil && s.Doc == nil {
										fmt.Printf("%s: type %s lacks comment (pkg %s)\n", filename, s.Name.Name, pkgName)
									}
								}
							case *ast.ValueSpec:
								for _, name := range s.Names {
									if name.IsExported() {
										if x.Doc == nil && s.Doc == nil {
											fmt.Printf("%s: var/const %s lacks comment (pkg %s)\n", filename, name.Name, pkgName)
										}
									}
								}
							}
						}
					}
				}
				return true
			})
		}
	}
}
