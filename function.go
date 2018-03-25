package mysqlib

import (
	"strconv"
	"strings"
)

//过滤SQL敏感字符
func filterSQL(v string) string {
	// 判断 "
	n := strings.Count(v, "\"")
	if n%2 != 0 {
		v = strings.Replace(v, "\"", "\"\"", 1)
	}
	// 判断 '
	n = strings.Count(v, "'")
	if n%2 != 0 {
		v = strings.Replace(v, "'", "''", 1)
	}
	// 判断 `
	n = strings.Count(v, "`")
	if n%2 != 0 {
		v = strings.Replace(v, "`", "``", 1)
	}
	// 判断 \
	v = strings.Replace(v, "\\", "\\\\", -1)
	// 判断 --
	v = strings.Replace(v, "--", "---", -1)

	return v
}

//将interface{}类型的参数断言并转换成string以用于拼接sql语句
func interfaceToString(v interface{}) string {
	var value string
	switch v.(type) {
	case string:
		value = filterSQL(v.(string))
	case int:
		value = strconv.Itoa(v.(int))
	case int8:
		value = strconv.Itoa(int(v.(int8)))
	case int16:
		value = strconv.Itoa(int(v.(int16)))
	case int32:
		value = strconv.Itoa(int(v.(int32)))
	case int64:
		value = strconv.FormatInt(v.(int64), 10)
	case float32:
		value = strconv.FormatFloat(v.(float64), 'f', -1, 32)
	case float64:
		value = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		value = strconv.FormatBool(v.(bool))
	default:
		value = "[!mysqlib!]"
	}

	return value
}

//将赋值给interface{}类型的slice参数断言并转换成string以用于拼接sql语句
func sliceToString(v interface{}) []string {
	var values []string

	switch v.(type) {
	case []string:
		value := v.([]string)
		for _, v := range value {
			values = append(values, filterSQL(v))
		}
	case []int:
		value := v.([]int)
		for _, v := range value {
			values = append(values, strconv.Itoa(v))
		}
	case []int8:
		value := v.([]int8)
		for _, v := range value {
			values = append(values, strconv.Itoa(int(v)))
		}
	case []int16:
		value := v.([]int16)
		for _, v := range value {
			values = append(values, strconv.Itoa(int(v)))
		}
	case []int32:
		value := v.([]int32)
		for _, v := range value {
			values = append(values, strconv.Itoa(int(v)))
		}
	case []int64:
		value := v.([]int64)
		for _, v := range value {
			values = append(values, strconv.FormatInt(v, 10))
		}
	case []float32:
		value := v.([]float32)
		for _, v := range value {
			values = append(values, strconv.FormatFloat(float64(v), 'f', -1, 32))
		}
	case []float64:
		value := v.([]float64)
		for _, v := range value {
			values = append(values, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case []bool:
		value := v.([]bool)
		for _, v := range value {
			values = append(values, strconv.FormatBool(v))
		}
	}

	return values
}
