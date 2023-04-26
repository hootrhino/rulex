package rulexlib

import (
	"encoding/json"
	"errors"

	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* JsonEncodeLib,_Encode converts Lua values to JSON.
*
 */
type JsonEncodeLib struct {
}

func NewJsonEncodeLib() typex.XLib {
	return &JsonEncodeLib{}
}

func (l *JsonEncodeLib) Name() string {
	return "JsonEncode"
}
func (l *JsonEncodeLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return apiEncode
}

/*
*
* JsonDecodeLib, _Decode converts the JSON encoded data to Lua values.
*
 */
type JsonDecodeLib struct {
}

func NewJsonDecodeLib() typex.XLib {
	return &JsonDecodeLib{}
}

func (l *JsonDecodeLib) Name() string {
	return "JsonDecode"
}
func (l *JsonDecodeLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return apiDecode
}
func JSONE(rx typex.RuleX) func(l *lua.LState) int {
	return apiEncode
}
func JSOND(rx typex.RuleX) func(l *lua.LState) int {
	return apiDecode
}
func apiDecode(L *lua.LState) int {
	str := L.CheckString(2)

	value, err := _Decode(L, []byte(str))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	L.Push(lua.LNil)
	return 2
}

func apiEncode(L *lua.LState) int {
	value := L.CheckAny(2)

	data, err := _Encode(value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	L.Push(lua.LNil)
	return 2
}

var (
	errNested      = errors.New("cannot encode recursively nested tables to JSON")
	errSparseArray = errors.New("cannot encode sparse array")
	errInvalidKeys = errors.New("cannot encode mixed or invalid key types")
)

type invalidTypeError lua.LValueType

func (i invalidTypeError) Error() string {
	return `cannot encode ` + lua.LValueType(i).String() + ` to JSON`
}

// _Encode returns the JSON encoding of value.
func _Encode(value lua.LValue) ([]byte, error) {
	return json.Marshal(jsonValue{
		LValue:  value,
		visited: make(map[*lua.LTable]bool),
	})
}

type jsonValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j jsonValue) MarshalJSON() (data []byte, err error) {
	switch converted := j.LValue.(type) {
	case lua.LBool:
		data, err = json.Marshal(bool(converted))
	case lua.LNumber:
		data, err = json.Marshal(float64(converted))
	case *lua.LNilType:
		data = []byte(`null`)
	case lua.LString:
		data, err = json.Marshal(string(converted))
	case *lua.LTable:
		if j.visited[converted] {
			return nil, errNested
		}
		j.visited[converted] = true

		key, value := converted.Next(lua.LNil)

		switch key.Type() {
		case lua.LTNil: // empty table
			data = []byte(`[]`)
		case lua.LTNumber:
			arr := make([]jsonValue, 0, converted.Len())
			expectedKey := lua.LNumber(1)
			for key != lua.LNil {
				if key.Type() != lua.LTNumber {
					err = errInvalidKeys
					return
				}
				if expectedKey != key {
					err = errSparseArray
					return
				}
				arr = append(arr, jsonValue{value, j.visited})
				expectedKey++
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(arr)
		case lua.LTString:
			obj := make(map[string]jsonValue)
			for key != lua.LNil {
				if key.Type() != lua.LTString {
					err = errInvalidKeys
					return
				}
				obj[key.String()] = jsonValue{value, j.visited}
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(obj)
		default:
			err = errInvalidKeys
		}
	default:
		err = invalidTypeError(j.LValue.Type())
	}
	return
}

// _Decode converts the JSON encoded data to Lua values.
func _Decode(L *lua.LState, data []byte) (lua.LValue, error) {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return DecodeValue(L, value), nil
}

// DecodeValue converts the value to a Lua value.
//
// This function only converts values that the encoding/json package decodes to.
// All other values will return lua.LNil.
func DecodeValue(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case json.Number:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(DecodeValue(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), DecodeValue(L, item))
		}
		return tbl
	case nil:
		return lua.LNil
	}

	return lua.LNil
}
