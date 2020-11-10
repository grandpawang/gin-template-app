package utils

import (
	"reflect"
)

// GetByKey interface{} get item by it's Field Name
func GetByKey(form interface{}, fieldName string) reflect.Value {
	val := reflect.ValueOf(form)
	if val.Kind() != reflect.Ptr {
		val = val.FieldByName(fieldName)
	} else {
		val = val.Elem().FieldByName(fieldName)
	}
	return val
}

func structMap(form interface{}, out map[string]interface{}) map[string]interface{} {
	formVal := indirect(reflect.ValueOf(form))
	formType := formVal.Type()

	if out == nil {
		out = make(map[string]interface{})
	}

	for i := 0; i < formVal.NumField(); i++ {
		fieldVal := formVal.Field(i)
		fieldType := formType.Field(i)
		tagName := fieldType.Tag.Get("json")
		if fieldType.Type.Kind() == reflect.Struct {
			out[tagName] = structMap(fieldVal.Interface(), nil)
		} else {
			if tmp := indirect(fieldVal); tmp.Kind() != reflect.Invalid {
				out[tagName] = tmp.Interface()
			}
		}
	}
	return out
}

// StructMap struct to single dimension map[string]interface{}
// the form item must be ptr and had tag `json`
func StructMap(form interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	return structMap(form, res)
}

// Struct1DMapHadNil struct to any dimension map[string]interface{} field had nil map
// the form data must be ptr and had tag `form`
func Struct1DMapHadNil(form interface{}, out map[string]interface{}) map[string]interface{} {
	formVal := indirect(reflect.ValueOf(form))
	formType := formVal.Type()

	if out == nil {
		out = make(map[string]interface{})
	}

	for i := 0; i < formVal.NumField(); i++ {
		fieldVal := indirect(formVal.Field(i))
		fieldType := formType.Field(i)
		tagName := fieldType.Tag.Get(" json")

		if fieldType.Type.Kind() == reflect.Struct { // prt is pass
			Struct1DMapHadNil(fieldVal.Interface(), out)
		} else {
			if !fieldVal.CanAddr() {
				out[tagName] = nil
			} else {
				out[tagName] = fieldVal.Interface()
			}
		}

	}
	return out
}

// CopyInterface copy src to target, like: CopyInterface(src / &src, &tar).
// the data must be:
// 1. if Resources(FieldName) == Target(FieleName) && Resources(Type) == Target(Type) => copy
// 2. target => type Form struct { filed interface }
func CopyInterface(src interface{}, tar interface{}) {
	srcVal := indirect(reflect.ValueOf(src))
	tarVal := reflect.ValueOf(tar).Elem()
	tarType := tarVal.Type()
	for i := 0; i < tarVal.NumField(); i++ {
		tarField := tarType.Field(i)
		srcVal := srcVal.FieldByName(tarField.Name)
		if srcKind := srcVal.Kind(); srcKind != reflect.Invalid &&
			tarField.Type == srcVal.Type().Elem() && !srcVal.IsNil() {
			tarVal.FieldByName(tarField.Name).Set(srcVal.Elem())
		}
	}
}
