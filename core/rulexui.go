package core

import (
	"fmt"
	"reflect"
	"rulex/typex"
	"strings"
)

//
// 空串
//
const EMPTY_STRING string = ""

/*
*
* 这里是一些前端渲染支持
*
 */
type viewType string

const (
	_NUMBER viewType = "el-number"
	_TEXT   viewType = "el-text"
	_INLINE viewType = "el-inline"
	_SELECT viewType = "el-select"
	_FILE   viewType = "el-upload"
)

type view struct {
	Order       int         `json:"order"`       // 界面顺序
	Type        viewType    `json:"type"`        // 组件类型
	Name        string      `json:"name"`        // 表单字段名
	Info        string      `json:"info"`        // 表单提示
	Label       string      `json:"label"`       // 界面显示标签
	Value       interface{} `json:"value"`       // 字段的值
	Required    bool        `json:"required"`    // 是否必填
	Hidden      bool        `json:"hidden"`      // 是否隐藏
	Placeholder string      `json:"placeholder"` // 占位文本
}
type numberInputView struct {
	view
}

func NewNumberInputView() numberInputView {
	v := numberInputView{}
	v.Type = _NUMBER
	return v

}

type textInputView struct {
	view
}

func NewTextInputView() textInputView {
	v := textInputView{}
	v.Hidden = false
	v.Required = true
	v.Type = _TEXT
	return v

}

type inLineView struct {
	view
	Children []view `json:"children"`
}
type fileView struct {
	view
}

func NewFileView() fileView {
	v := fileView{}
	v.Type = _FILE
	return v
}
func NewInlineView() inLineView {
	v := inLineView{}
	v.Type = _INLINE
	return v
}

type selectView struct {
	view
	SelectOptions []selectOption `json:"selectOptions"`
}
type selectOption struct {
	Label string `json:"label"` // 下拉选择框 UI
	Value string `json:"value"` // 下拉选择框 值
}

func NewSelectView() selectView {
	v := selectView{}
	v.Type = _SELECT
	return v

}

/*
*
* 渲染 select 组件, Select 在Config上用tag 表示
* 其格式如下: options:"label1,value1|label2,value2 ...."
*
 */
func renderSelect(s string) ([]selectOption, error) {
	options := []selectOption{}
	var err error
	splits := strings.Split(s, "|")
	if len(splits)%2 != 0 {
		err = fmt.Errorf("expression error:%v", s)
		goto END
	}
	for _, kv := range splits {
		kv := strings.Split(kv, ",")
		if len(splits) != 2 {
			err = fmt.Errorf("syntax error:%v", kv)
			goto END
		}
		options = append(options, selectOption{Label: kv[0], Value: kv[1]})
	}
END:
	return options, err
}

/*
*
* 渲染UI界面
*
 */

func RenderConfig(Type string, helpTip string, config interface{}) (typex.XConfig, error) {
	var err error
	typee := reflect.TypeOf(config)
	views := make([]interface{}, 1)
	for i := 0; i < typee.NumField(); i++ {
		filedName := typee.Field(i).Name
		filedType := typee.Field(i).Type.String()
		tag := typee.Field(i).Tag
		//
		Info := tag.Get("info")
		if Info == EMPTY_STRING {
			Info = "暂无内容" // 当看到这个文本的时候就需要去增加这个Tag字段了
		}
		Label := tag.Get("label")
		if Label == EMPTY_STRING {
			Label = "暂无内容" // 当看到这个文本的时候就需要去增加这个Tag字段了
		}
		Required := tag.Get("required")
		Hidden := tag.Get("hidden")
		Placeholder := tag.Get("placeholder")
		//
		xcfg := typex.XConfig{}
		// 数字输入
		if (filedType == "int") ||
			(filedType == "int64") ||
			(filedType == "int32") ||
			(filedType == "float32") {
			xcfg.Type = string(_NUMBER)
			nv := NewNumberInputView()
			//
			nv.Order = i
			nv.Name = filedName
			nv.Label = Label
			nv.Info = Info
			nv.Placeholder = Placeholder
			if Required == "false" {
				nv.Required = false
			}
			if Hidden == "true" {
				nv.Hidden = true
			}
			views = append(views, nv)
		}
		// 文本输入
		if filedType == "string" || filedType == "*string" {

			// 文件框
			fileTag := tag.Get("file") // file:"uploadfile"
			if fileTag != EMPTY_STRING {
				xcfg.Type = string(_FILE)
				fv := NewFileView()
				//
				fv.Order = i
				fv.Name = fileTag
				fv.Label = Label
				fv.Info = Info
				fv.Placeholder = Placeholder
				if Required == "false" {
					fv.Required = false
				}
				if Hidden == "true" {
					fv.Hidden = true
				}
				views = append(views, fv)
			} else {
				xcfg.Type = string(_TEXT)
				tv := NewTextInputView()
				//
				tv.Order = i
				tv.Name = filedName
				tv.Label = Label
				tv.Info = Info
				tv.Placeholder = Placeholder
				if Required == "false" {
					tv.Required = false
				}
				if Hidden == "true" {
					tv.Hidden = true
				}

				views = append(views, tv)
			}
		}
		// 动态数组
		if (filedType == "[]string") ||
			(filedType == "[]int") ||
			(filedType == "[]int32") ||
			(filedType == "[]int64") {
			xcfg.Type = string(_INLINE)
			iv := NewInlineView()
			//
			iv.Order = i
			iv.Name = filedName
			iv.Label = Label
			iv.Info = Info
			iv.Placeholder = Placeholder
			if Required == "false" {
				iv.Required = false
			}
			if Hidden == "true" {
				iv.Hidden = true
			}

			views = append(views, iv)
		}

		optionsTag := tag.Get("options")
		// 下拉框输入
		if optionsTag != EMPTY_STRING {
			xcfg.Type = string(_SELECT)
			sv := NewSelectView()
			//
			sv.Order = i
			sv.Name = filedName
			sv.Label = Label
			sv.Info = Info
			sv.Placeholder = Placeholder
			if Required == "false" {
				sv.Required = false
			}
			if Hidden == "true" {
				sv.Hidden = true
			}
			selectOptions, err1 := renderSelect(optionsTag)
			if err1 != nil {
				err = err1
				goto END
			}
			sv.SelectOptions = selectOptions
			views = append(views, sv)
		}

	}
END:
	return typex.XConfig{Type: Type, Views: views, HelpTip: helpTip}, err
}
