package x

import (
	"reflect"
	"strconv"

	"github.com/marianogappa/sqlparser"
	"github.com/ngaut/log"
)

func SqlParse(data *map[string]interface{}, sql string) (*map[string]interface{}, error) {
	log.Debug("InputData:", data)
	query, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	} else {
		if query.TableName == "INPUT_DATA" {
			fieldsResult := map[string]interface{}{}
			for _, f := range query.Fields {
				fieldsResult[f] = nil
			}
			conditions := query.Conditions
			for _, condition := range conditions {
				if condition.Operand1IsField {
					inputValue := (*data)[condition.Operand1]
					log.Debug("InputValueTypeOf", reflect.TypeOf(inputValue))
					if inputValue != nil {
						log.Debug("Field:", condition.Operand1, " InputValue:", inputValue, " GivenLimit:", condition.Operand2)
						switch {
						case condition.Operator == 1: //=
							{
								if inputValue == condition.Operand2 {
									fieldsResult[condition.Operand1] = inputValue
								}
							}
						case condition.Operator == 2: // "!="
							{
								if inputValue != condition.Operand2 {
									fieldsResult[condition.Operand1] = inputValue
								}
							}
						case condition.Operator == 3: // >
							{
								if reflect.TypeOf(inputValue).Kind() == reflect.Float64 {
									value, err := strconv.ParseFloat(condition.Operand2, 64)
									if err == nil {
										if inputValue.(float64) > value {
											fieldsResult[condition.Operand1] = inputValue
										}
									}
								}
							}
						case condition.Operator == 4: //<
							{
								if reflect.TypeOf(inputValue).Kind() == reflect.Float64 {
									value, err := strconv.ParseFloat(condition.Operand2, 64)
									if err == nil {
										if inputValue.(float64) < value {
											fieldsResult[condition.Operand1] = inputValue
										}
									}
								}

							}
						case condition.Operator == 5: // >=
							{
								if reflect.TypeOf(inputValue).Kind() == reflect.Float64 {
									value, err := strconv.ParseFloat(condition.Operand2, 64)
									if err == nil {
										if inputValue.(float64) >= value {
											fieldsResult[condition.Operand1] = inputValue
										}
									}
								}

							}
						case condition.Operator == 6: // <=
							{
								if reflect.TypeOf(inputValue).Kind() == reflect.Float64 {
									value, err := strconv.ParseFloat(condition.Operand2, 64)
									if err == nil {
										if inputValue.(float64) <= value {
											fieldsResult[condition.Operand1] = inputValue
										}
									}
								}
							}
						default:
							{
							}
						}

					}
				}

			}
			return &fieldsResult, nil
		}
	}
	return nil, err
}
