package test

import (
	"fmt"
	"reflect"
	"testing"
)

// go test -timeout 30s -run ^TestJsonFilter github.com/hootrhino/rulex/test -v -count=1
func TestJsonFilter(t *testing.T) {
	// SELECT * FROM DATA WHERE a1.value GT(>) 3.14
	r := Reduce(map[string]interface{}{
		"a1": float64(3.14),
		"a2": float64(3.15),
		"a3": float64(3.16),
		"a4": float64(3.17),
	}, []FilterCondition{
		{
			Key:      "*",
			Label:    "tag-a1",
			Operator: "GT",
			InValue:  float64(3.14),
			Type:     "number",
		},
	})
	t.Log(r)

}

/*
*
* 用 filterCondition 去 in 里面拿数据
*
 */
type FilterCondition struct {
	Key      string      // 期望的Key, 如果是 * 则表示只要满足条件的全部拿出来
	Label    string      // 期望的Label
	Operator string      // 操作条件 EQ:==, LT:<, GT:>, LTE:<=, GTE:>=
	Type     string      // 数据类型 1 number 2 string
	InValue  interface{} // 对比的值1
	OutValue interface{} // 对比的值2
}

func Reduce(in map[string]interface{}, fcs []FilterCondition) map[string]interface{} {
	retV := map[string]interface{}{}
	for k, v := range in {
		for _, fc := range fcs {

			v, MatchedOk, err := Match(v, fc)
			if err != nil {
				return nil
			}
			if MatchedOk {
				// TODO : SELECT *
				if fc.Key == "*" {
					fc.OutValue = v
					retV[k] = fc
				} else {
					if fc.Key != k {
						continue
					}
					fc.OutValue = v
					retV[fc.Key] = fc
				}

			}

		}
	}
	return retV
}

func checkType(in interface{}) (string, bool) {
	tin := reflect.TypeOf(in).Name()
	ok := false
	t := ""
	switch tin {
	case "byte":
		{
			ok = true
			t = "number"

		}
	case "int":
		{
			ok = true
			t = "number"
		}
	case "int32":
		{
			ok = true
			t = "number"
		}
	case "int64":
		{
			ok = true
			t = "number"
		}
	case "float32":
		{
			ok = true
			t = "number"
		}
	case "float64":
		{
			ok = true
			t = "number"
		}
	case "string":
		{
			ok = true
			t = "string"
		}
	default:
	}
	return t, ok
}

func Match(in interface{}, fc FilterCondition) (interface{}, bool, error) {
	tin := reflect.TypeOf(in).Name()
	InType, checkOk := checkType(in)
	if InType != fc.Type {
		return nil, false, fmt.Errorf(
			"type of Data and type of 'Filter' must have same type: %v, %v",
			tin, fc.Type)
	}
	if !checkOk {
		return nil, false, fmt.Errorf("target data type invalid, in: %v, fc: %v", tin, fc.Type)
	}

	switch fc.Type {
	case "string":
		{
			s := in.(string)
			if fc.Operator == "EQ" {
				r, ok := matchS(s, fc.InValue.(string), func(in, fc string) (string, bool) {
					return StringEQ(in, fc)
				})
				return r, ok, nil
			}
			if fc.Operator == "NEQ" {
				r, ok := matchS(s, fc.InValue.(string), func(in, fc string) (string, bool) {
					return StringNEQ(in, fc)
				})
				return r, ok, nil
			}
		}
	case "number":
		{
			if fc.Operator == "EQ" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericEQ(in, fc)
					})
				return r, ok, nil
			}
			if fc.Operator == "NEQ" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericNEQ(in, fc)
					})
				return r, ok, nil
			}
			if fc.Operator == "LT" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericLT(in, fc)
					})
				return r, ok, nil
			}
			if fc.Operator == "LTE" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericLTE(in, fc)
					})
				return r, ok, nil
			}
			if fc.Operator == "GT" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericGT(in, fc)
					})
				return r, ok, nil
			}
			if fc.Operator == "GTE" {
				r, ok := matchN(in.(float64), fc.InValue.(float64),
					func(in, fc float64) (float64, bool) {
						return NumericGTE(in, fc)
					})
				return r, ok, nil
			}
		}
	default:
		{
		}
	}
	return nil, false, fmt.Errorf("unsupported data type,in: %v, fc: %v", tin, fc.Type)
}

// -------------------------------------------------------------------------------------------------
// 具体的比较实现
// -------------------------------------------------------------------------------------------------

type Numeric interface {
	~byte | ~int | ~int32 | ~int64 | ~float32 | ~float64
}

type String interface {
	~string
}

func matchN[T Numeric](in, fc T, mf func(in, fc T) (T, bool)) (T, bool) {
	return mf(in, fc)
}
func matchS[T String](in, fc T, mf func(in, fc T) (T, bool)) (T, bool) {
	return mf(in, fc)
}

func StringEQ[T String](in, fc T) (T, bool) {
	if in == fc {
		return in, true
	}
	return in, false
}
func StringNEQ[T String](in, fc T) (T, bool) {
	if in != fc {
		return in, true
	}
	return "", false
}

func NumericEQ[T Numeric](in, fc T) (T, bool) {
	if in == fc {
		return in, true
	}
	return 0, false
}

func NumericNEQ[T Numeric](in, fc T) (T, bool) {
	if in != fc {
		return in, true
	}
	return 0, false
}

func NumericLT[T Numeric](in, fc T) (T, bool) {
	if in < fc {
		return in, true
	}
	return 0, false
}

func NumericLTE[T Numeric](in, fc T) (T, bool) {
	if in <= fc {
		return in, true
	}
	return 0, false
}

func NumericGT[T Numeric](in, fc T) (T, bool) {
	if in > fc {
		return in, true
	}
	return 0, false
}

func NumericGTE[T Numeric](in, fc T) (T, bool) {
	if in >= fc {
		return in, true
	}
	return 0, false
}
