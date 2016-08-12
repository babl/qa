package kafkalogs

import "reflect"

func isValidField(y interface{}, ytype reflect.Kind) bool {
	if y != nil && reflect.TypeOf(y).Kind() == ytype {
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

func getFieldDataString(x interface{}) string {
	return getFieldData(x, reflect.String).(string)
}

func getFieldDataInt32(x interface{}) int32 {
	return int32(getFieldData(x, reflect.Float64).(float64))
}

func getFieldDataFloat64(x interface{}) float64 {
	return getFieldData(x, reflect.Float64).(float64)
}
