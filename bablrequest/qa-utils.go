package bablrequest

import "reflect"

func isValidField(y interface{}, ytype reflect.Kind) bool {
	if y != nil && (reflect.TypeOf(y).Kind() == ytype) {
		return true
	}
	return false
}

func getFieldData(x interface{}, xtype reflect.Kind) interface{} {
	var result interface{} = nil
	if isValidField(x, xtype) {
		return x
	} else {
		switch xtype {
		case reflect.String:
			result = ""
		case reflect.Float64:
			result = float64(0.0)
		}
	}
	return result
}

/*
func getFieldDataStringArray(x map[string]interface{}, fieldname string) string {
	var result string
	list := x[fieldname].([]interface{})
	for _, v := range list {
		result += v.(string) + ","
	}
	fmt.Println("Result: ", result)
	return result
}*/

func getFieldDataString(x interface{}) string {
	return getFieldData(x, reflect.String).(string)
}

func getFieldDataInt(x interface{}) int {
	return int(getFieldData(x, reflect.Float64).(float64))
}

func getFieldDataInt32(x interface{}) int32 {
	return int32(getFieldData(x, reflect.Float64).(float64))
}

func getFieldDataFloat64(x interface{}) float64 {
	return getFieldData(x, reflect.Float64).(float64)
}
