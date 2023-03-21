package utils

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

//
// JSON String to a struct, (can't validate map!!!)
//
func TransformConfig(s1 []byte, s2 interface{}) error {
	if err := json.Unmarshal(s1, &s2); err != nil {
		return err
	}
	if err := validator.New().Struct(s2); err != nil {
		return err
	}
	return nil
}

//
// Bind config to struct
// config: a Map, s: a struct variable
//
func BindConfig(config map[string]interface{}, s interface{}) error {
	return BindSourceConfig(config, s)
}
func BindSourceConfig(config map[string]interface{}, s interface{}) error {
	configBytes, err0 := json.Marshal(&config)
	if err0 != nil {
		return err0
	}
	if err := json.Unmarshal(configBytes, &s); err != nil {
		return err
	}
	if err := validator.New().Struct(s); err != nil {
		return err
	}
	return nil
}

//
// Map转成结构体
//Decode takes an input structure and uses reflection to translate it to the output structure.
//    output must be a pointer to a map or struct.
//
func MaoToTerm(input interface{}, output interface{}) error {
	return mapstructure.Decode(input, output)
}
