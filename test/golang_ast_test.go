package test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

/*
*
* 准备写个简单的翻译器，把Golang代码直接翻译成Lua
* 不过看起来应该是个体力活，需要精打细磨，先立Flag再说
*
 */
func Test_parse_AST(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "data/ast_test_example.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	// ast.Print(fset, node)
	for _, v := range node.Decls {
		switch typa := v.(type) {
		case *ast.FuncDecl:
			{
				if typa.Name.String() == "GoToLuaDemo" {
					t.Log("函数注释:", typa.Doc.Text())
					t.Log("函数名:", typa.Name)
					for _, param := range typa.Type.Params.List {
						t.Log("函数参数:", param.Type, param.Names)
					}
					for _, param := range typa.Type.Results.List {
						t.Log("函数返回值:", param.Type, param.Names)
					}
				}
			}
		}
	}

}

/*
*
* 尝试使用AST来生成文档
*
 */
func Test_gen_rulexlib_with_ast(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "data/ast_test_example.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	//ast.Print(fset, node)
	for _, v := range node.Decls {
		switch typa := v.(type) {
		case *ast.FuncDecl:
			{
				if strings.HasPrefix(typa.Name.String(), "__RULEXLIB") {
					fmt.Println("## " + typa.Name.String()[11:])
					for _, line := range typa.Doc.List {
						if strings.HasPrefix(line.Text, "//@") {
							specs := strings.Split(line.Text, ":")
							if len(specs) != 2 {
								continue
							}
							name := specs[0][3:]
							value := specs[1]
							fmt.Println(name, value)
						}
					}
					fmt.Println("### 参数表")
					tHeader := "|名称|类型|描述|\n| --- | --- | --- |"
					fmt.Println(tHeader)
					for _, param := range typa.Type.Params.List {
						l := fmt.Sprintf("|%s|%s|%v|", param.Type, param.Names[0], param.Comment)
						fmt.Println(l)
					}

					fmt.Println("### 返回值表")
					for _, result := range typa.Type.Results.List {
						fmt.Println(result.Type)
					}

				}
			}
		}
	}

}
