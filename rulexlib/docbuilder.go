package rulexlib

/*
*
* 这个是个很简单的文档生成器 主要以后用来生成lua库的文档用
*
 */
import (
	"fmt"
	"log"
	"os"
)

/*
*
* 生成文档的最小单元
*
 */
type FunArg struct {
	Pos         int
	Type        string
	Description string
}
type ReturnValue struct {
	Pos         int
	Type        string
	Description string
}
type Fun struct {
	NameSpace   string        // 函数名
	FunName     string        // 函数名
	FunArgs     []FunArg      // 函数参数
	ReturnValue []ReturnValue // 函数返回值
	Description string        // 描述文本
	Example     string        // 示例
}
type RulexLibDoc struct {
	Name        string
	Version     string
	ReleaseTime string
	Funcs       []Fun
}

func ifEmpty(s string) string {
	if s == "" {
		return "暂无信息"
	} else {
		return s
	}
}
func (doc *RulexLibDoc) AddFunc(f Fun) {
	doc.Funcs = append(doc.Funcs, f)
}
func (doc *RulexLibDoc) BuildDoc() {
	body := "# " + doc.Name
	tHeader := "\n|版本|发布时间|\n| --- | --- |\n"
	body += tHeader + "|" + doc.Version + "|" + doc.ReleaseTime + "|\n"
	for _, v := range doc.Funcs {
		body += v.BuildSection()
	}
	fmt.Println("准备生成文档:", doc.Name+"-"+doc.Version+"-"+doc.ReleaseTime+".md")
	file, err := os.Create("./" + doc.Name + "-" + doc.Version + "-" + doc.ReleaseTime + ".md")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(body)
	file.Close()
	fmt.Println("文档生成结束")

}
func (fun *Fun) BuildSection() string {
	body := "## " + fun.NameSpace + ":" + fun.FunName + "\n"
	body += "命名空间:`" + fun.NameSpace + "`\n"
	body += "### 简介\n" + fun.Description + "\n"
	tHeader := "|位置|类型|描述|\n| --- | --- | --- |\n"
	argsLine := ""
	// 参数
	for _, arg := range fun.FunArgs {
		argsLine += "|" + fmt.Sprintf("%v", arg.Pos) + "|" + arg.Type + "|" + ifEmpty(arg.Description) + "|\n"
	}
	body += "### 参数\n" + tHeader + argsLine
	// 返回值
	returnLine := ""
	// 参数
	for _, arg := range fun.ReturnValue {
		returnLine += "|" + fmt.Sprintf("%v", arg.Pos) + "|" + arg.Type + "|" + ifEmpty(arg.Description) + "|\n"
	}
	body += "### 返回\n" + tHeader + returnLine
	body += "### 示例\n" + "```lua\n" + fun.Example + "\n```\n"
	return body
}
