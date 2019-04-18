package utils

import (
	"fmt"
	"strconv"
)

func FieldAsString(data map[string]interface{}, fieldName string) string {
	if len(data) < 1 {
		return ""
	}

	if v, ok := data[fieldName]; ok {
		if vv, ok := v.(string); ok {
			return vv
		}
		if vv, ok := v.(float64); ok {
			return strconv.FormatInt(int64(vv), 10)
		}
		if vv, ok := v.(int64); ok {
			return strconv.FormatInt(vv, 10)
		}
		return fmt.Sprintf("%v", v)
	}

	return ""
}

func FieldAsInt(data map[string]interface{}, fieldName string) int {
	if len(data) < 1 {
		return 0
	}

	if v, ok := data[fieldName]; ok {
		if vv, ok := v.(int); ok {
			return vv
		}
		if vv, ok := v.(int64); ok {
			return int(vv)
		}
		if vv, ok := v.(float64); ok {
			return int(vv)
		}
		if vv, ok := v.(string); ok {
			if vvv, e := strconv.Atoi(vv); e == nil {
				return vvv
			}
		}
	}

	return 0
}
