package device

import (
	"log"
	"os"
	"time"

	"github.com/i4de/rulex/typex"
)

type DeviceDoc struct {
	Name           string
	Version        string
	ReleaseTime    string
	DeviceSections []DeviceSection
}

func (doc *DeviceDoc) AddDeviceSection(s DeviceSection) {
	doc.DeviceSections = append(doc.DeviceSections, s)
}

type DeviceSection struct {
	Name        string
	Description string
	Config      Section
	OnRead      Section
	OnWrite     Section
	Example     string
}

type Section struct {
	Struct      []Var
	ExampleJson string
}
type Var struct {
	Name        string
	Type        string
	Description string
}

func ifEmpty(s string) string {
	if s == "" {
		return "暂无信息"
	} else {
		return s
	}
}

func (doc *DeviceDoc) BuildDoc() {
	body := "# " + doc.Name
	tHeader := "\n|版本|发布时间|\n| --- | --- |\n"
	body += tHeader + "|" + doc.Version + "|" + doc.ReleaseTime + "|\n"
	file, err := os.Create("./" + doc.Name + "-" + doc.Version + "-" + doc.ReleaseTime + ".md")
	if err != nil {
		log.Fatal(err)
	}
	sectiontH := `
|名称|类型|描述|
| --- | --- | --- |`
	for _, section := range doc.DeviceSections {
		body += "## " + ifEmpty(section.Name) + "\n"
		body += "### 简介\n"
		body += ifEmpty(section.Description) + "\n"
		body += "### 配置参数\n"
		body += sectiontH + "\n"
		sectionLine := ""
		for _, vars := range section.Config.Struct {
			sectionLine += "|" + ifEmpty(vars.Name) + "|" + ifEmpty(vars.Type) + "|" + ifEmpty(vars.Description) + "|\n"
		}
		sectionLine += "### 示例配置\n```json\n" + ifEmpty(section.Config.ExampleJson) + "\n```" + "\n"
		body += sectionLine
		body += "### 示例程序\n```lua\n" + ifEmpty(section.Example) + "\n```" + "\n"
	}
	file.WriteString(body)
	file.Close()
	log.Println("文档生成结束")
}

/*
*
* 构建文档
*
 */
func BuildDoc() {
	currentTime := time.Now()
	var deviceDoc DeviceDoc = DeviceDoc{
		Name:        "RULEX-外部设备文档",
		Version:     typex.DefaultVersion.Version,
		ReleaseTime: currentTime.Format("2006-01-02"),
	}
	deviceDoc.AddDeviceSection(DeviceSection{
		Name:        "TSS200",
		Description: "TSS200是一款国产环境参数传感器",
		Config: Section{
			ExampleJson: `
{
	"uuid": "TSS200V02",
	"name": "TSS200V02",
	"type": "TSS200V02",
	"description": "TSS200V02",
	"config": {
		"config": {
			"baudRate": 9600,
			"dataBits": 8,
			"ip": "127.0.0.1",
			"parity": "N",
			"port": 502,
			"stopBits": 1,
			"uart": "COM2"
		},
		"frequency": 5,
		"mode": "RTU",
		"registers": [
			{
				"address": 17,
				"function": 3,
				"quantity": 2,
				"slaverId": 1,
				"tag": "node1"
			}
		],
		"timeout": 10
	}
}`,
		},
		OnRead: Section{
			ExampleJson: "{}",
		},
		OnWrite: Section{

			ExampleJson: "{}",
		},
		Example: `
Actions = {
	function(data)
		print('data ==> ', data)
		return true, data
	end
}
`,
	})
	deviceDoc.BuildDoc()
}
