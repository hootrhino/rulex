package core

import (
	"errors"
	"fmt"
	"reflect"
	"rulex/typex"
	"strings"
)

/*
*
* 这里是一些前端渲染支持
*
 */
type selctView struct {
	Name string
	Key  string
}
type stringView struct {
	Name string
	Key  string
}
type numberView struct {
	Name string
	Key  string
}

/*
*
* 渲染 select 组件
*
 */
func RenderSelect(s string) ([]selctView, error) {
	selctors := []selctView{}
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
		selctors = append(selctors, selctView{Key: kv[0], Name: kv[1]})
	}
END:
	return selctors, err
}

/*
*
* 渲染UI界面
*
 */
func RenderConfig(i interface{}) ([]typex.XConfig, error) {
	var err error
	typee := reflect.TypeOf(i)
	data := make([]typex.XConfig, 0)
	for i := 0; i < typee.NumField(); i++ {
		xcfg := typex.XConfig{}
		field := typee.Field(i)
		tag := field.Tag
		title := tag.Get("title")
		if title == "" {
			err = errors.New("'title' tag can't empty")
			goto END
		}
		info := tag.Get("info")
		json := tag.Get("json")
		if json == "" {
			err = errors.New("'json' tag can't empty")
			goto END
		}
		// 枚举常量
		enum := tag.Get("enum")
		if enum != "" {
			data, err1 := RenderSelect(enum)
			if err1 != nil {
				err = err1
				goto END
			}
			xcfg.Enum = data
		}
		// 枚举的类型
		fieldType := typee.Field(i).Type.String()
		if fieldType[:2] == "[]" {
			xcfg.Multiple = true
		} else {
			xcfg.Multiple = false
		}
		// 文件
		file := tag.Get("file")
		if file != "" {
			xcfg.FieldType = "upload"
		} else {
			xcfg.FieldType = fieldType
		}
		xcfg.Field = json
		xcfg.Info = info
		xcfg.Title = title
		data = append(data, xcfg)
	}
END:
	return data, err
}
