package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJson(t *testing.T) {

	rule := map[string]interface{}{
		"name":        "just_a_test",
		"description": "just_a_test",
		"actions": strings.Replace(`
					local json = require("json")
					Actions = {
						function(data)
							local s = '{"temp":100,"hum":30, "co2":123.4, "lex":22.56}'
							print(s == data)
							DataToMongo("$${OUT}", s)
							return true, data
						end
					}`, "$${OUT}", "m_Out_id_1.UUID", -1),
		"from": "mIn_id_1.UUID",
		"failed": `
			   function Failed(error)
			   print("call error:",error)
			   end`,
		"success": `
			   function Success()
			   print("call success")
			   end`,
	}
	b, _ := json.Marshal(&rule)
	fmt.Printf(string(b))
}
