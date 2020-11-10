package utils

import (
	"reflect"
)

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// Substring get sub string
func Substring(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// DifferentIDs two slice ids differents, how src to tar, true is delete, false is insert.
func DifferentIDs(src, tar []uint) map[uint]bool {
	res := make(map[uint]bool)
	for _, v := range src {
		res[v] = true
	}
	for _, v := range tar {
		if res[v] == false { // src don't set
			res[v] = false
		} else {
			delete(res, v)
		}
	}
	return res
}

// DuplicateUint make the duplicate slice be unique and the slice must be sort
func DuplicateUint(a *[]uint) []uint {
	tmp := make(map[uint]bool)
	for _, i := range *a {
		tmp[i] = true
	}
	a = new([]uint)
	for k := range tmp {
		*a = append(*a, k)
	}
	return *a
}
