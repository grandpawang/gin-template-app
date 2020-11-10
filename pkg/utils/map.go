package utils

import (
	"fmt"
	"reflect"
)

// GetKeys map[string]interface{} get keys
func GetKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapEqual map equal
func MapEqual(a, b map[string]interface{}) bool {
	for k, v := range a {
		fmt.Println(reflect.TypeOf(b[k]), reflect.TypeOf(v))
		fmt.Println(k, b[k], v)
		if b[k] != v {
			return false
		}
	}
	return true
}
