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
	node, err := parser.ParseFile(fset, "data/ast_test_example1.spec", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	commentsMap := map[int]string{}
	for _, comment := range node.Comments {
		// @arg 表示参数注释
		if strings.HasPrefix(comment.Text(), "@arg:") {
			position := fset.PositionFor(comment.Pos(), false)
			text := comment.Text()[5:] // 去掉前面的@arg标识
			if text[len(text)-1] == '\n' {
				commentsMap[position.Line] = text[:len(text)-1] //去掉换行符
			} else {
				commentsMap[position.Line] = text
			}
		}
	}
	//ast.Print(fset, node)
	for _, v := range node.Decls {
		switch typa := v.(type) {
		case *ast.FuncDecl:
			{
				// __RULEX_rulexlib_DataToTdEngine
				if strings.HasPrefix(typa.Name.String(), "__RULEX_") {
					fmt.Println("## " + typa.Name.String()[8:])
					fmt.Println("### 命名空间")
					fmt.Println("- " + node.Name.Name + "\n")
					for _, line := range typa.Doc.List {
						if strings.HasPrefix(line.Text, "//@") {
							specs := strings.Split(line.Text, ":")
							if len(specs) != 2 {
								continue
							}
							key := specs[0][3:]
							value := specs[1]
							if strings.HasPrefix(key, "desc") {
								fmt.Println("### 概述")
								fmt.Println(value + "\n")
							}
						}
					}
					fmt.Println("### 参数表")
					tHeader1 := "|名称|类型|描述|\n| --- | --- | --- |"
					fmt.Println(tHeader1)
					// exampleLua := "function($P) return $R end"
					args := []string{}
					returns := []string{}
					for _, param := range typa.Type.Params.List {
						position := fset.File(param.Pos()).PositionFor(param.Pos(), false)
						l := fmt.Sprintf("|%s|%s|%v|", param.Names[0], param.Type, commentsMap[position.Line])
						args = append(args, param.Names[0].String())
						fmt.Println(l)
					}
					if typa.Type.Results.List != nil {
						fmt.Println("### 返回值表")
						tHeader2 := "|名称|类型|描述|\n| --- | --- | --- |"
						fmt.Println(tHeader2)
						for _, param := range typa.Type.Results.List {
							position := fset.File(param.Pos()).PositionFor(param.Pos(), false)
							l := fmt.Sprintf("|%s|%s|%v|", param.Type, param.Type, commentsMap[position.Line])
							returns = append(returns, fmt.Sprintf("%s", param.Type))
							fmt.Println(l)
						}
					}
					// Example
					fmt.Println("### 示例代码")

					lua := `
local $R = $N:$F($A)
if (error~=nil)
    -- do something
end
`
					lua = strings.Replace(lua, "$N", node.Name.Name, 1)
					lua = strings.Replace(lua, "$F", typa.Name.String()[8:], 1)
					lua = strings.Replace(lua, "$A", strings.Join(args, ","), 1)
					lua = strings.Replace(lua, "$R", strings.Join(returns, ","), 1)
					fmt.Println("```lua")
					fmt.Println(lua)
					fmt.Println("```")

				}
			}
		}
	}

}
