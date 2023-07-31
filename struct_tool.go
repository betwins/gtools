package gtools

import (
	"reflect"
	"strings"
)

func Struct2MapAny(in any, tagName string) map[string]any {
	out := make(map[string]interface{})

	if in == nil {
		return out
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return out
	}

	for i := 0; i < v.NumField(); i++ {
		fi := v.Type().Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			var keyName string
			if tagName == "gorm" {
				begin := strings.Index(tagValue, "column:") + 7
				end := strings.Index(tagValue[begin:], ";")
				if end < 0 {
					end = len(tagValue)
				} else {
					end = begin + end
				}
				keyName = tagValue[begin:end]
			} else {
				keyName = tagValue
			}
			if !isBlank(v.Field(i)) {
				out[keyName] = v.Field(i).Interface()
			}
		}
	}
	return out
}

func Struct2MapString(in any, tagName string) map[string]string {
	out := make(map[string]string)

	if in == nil {
		return out
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return out
	}

	for i := 0; i < v.NumField(); i++ {
		fi := v.Type().Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			var keyName string
			if tagName == "gorm" {
				begin := strings.Index(tagValue, "column:") + 7
				end := strings.Index(tagValue[begin:], ";")
				if end < 0 {
					end = len(tagValue)
				} else {
					end = begin + end
				}
				keyName = tagValue[begin:end]
			} else {
				keyName = tagValue
			}
			if !isBlank(v.Field(i)) {
				out[keyName] = Any(v.Field(i).Interface()).ToString()
			}
		}
	}
	return out
}

// CopyStruct 结构体复制, 忽略空值，暂不支持结构体内部map复制（有需要可扩展）
func CopyStruct[DST any](src any) DST {

	if reflect.ValueOf(src).Kind() != reflect.Pointer {
		panic("src must be pointer to struct object")
	}

	var dst DST
	dstValue := reflect.ValueOf(&dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	if srcValue.Kind() != reflect.Struct {
		panic("src must point to struct object")
	}

	if dstValue.Kind() != reflect.Struct {
		panic("DST must be struct type")
	}

	if src == nil {
		return dst
	}

	// Recursively copy the original.
	copyStructRecursive(srcValue.Addr().Interface(), dstValue.Addr().Interface(), true)

	return dst
}

// CopyStructTo 复制到目标结构体对象，忽略空值，暂不支持结构体内部map复制（有需要可扩展）
func CopyStructTo(dst any, src any) {

	if reflect.ValueOf(src).Kind() != reflect.Pointer {
		panic("src must be pointer to struct object")
	}

	if reflect.ValueOf(dst).Kind() != reflect.Pointer {
		panic("dst must be pointer to struct object")
	}

	// Make the interface a reflect.Value
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	if srcValue.Kind() != reflect.Struct {
		panic("src must point to struct object")
	}

	if dstValue.Kind() != reflect.Struct {
		panic("dst must point to struct object")
	}

	if src == nil {
		return
	}

	// Recursively copy the original.
	copyStructRecursive(srcValue.Addr().Interface(), dstValue.Addr().Interface(), false)

	return
}

func copyStructRecursive(src, dst interface{}, replaceAll bool) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()
	if (replaceAll || srcValue.Kind() != reflect.Struct) && srcValue.Type() == dstValue.Type() && !isBlank(srcValue) && dstValue.CanSet() {
		dstValue.Set(srcValue)
		return
	}

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcName := srcValue.Type().Field(i).Name
		dstFieldByName := dstValue.FieldByName(srcName)
		if !dstFieldByName.IsValid() {
			continue
		}

		switch srcField.Kind() {
		case reflect.Map:
			continue
		case reflect.Slice:
			// Make a new slice and copy each element.
			if srcField.IsNil() {
				continue
			}
			if dstFieldByName.Kind() != reflect.Slice {
				continue
			}
			dstFieldByName.Set(reflect.MakeSlice(dstFieldByName.Type(), srcField.Len(), srcField.Cap()))
			for j := 0; j < srcField.Len(); j++ {

				srcInterface := srcField.Index(j).Addr().Interface()
				dstInterface := dstFieldByName.Index(j).Addr().Interface()
				copyStructRecursive(srcInterface, dstInterface, replaceAll)
			}
		case reflect.Struct:
			if dstFieldByName.Kind() != reflect.Struct {
				continue
			}
			copyStructRecursive(srcField.Addr().Interface(), dstFieldByName.Addr().Interface(), replaceAll)

		default:
			if dstFieldByName.CanSet() && !isBlank(srcField) && dstFieldByName.Type() == srcField.Type() {
				dstFieldByName.Set(srcField)
			}
		}
	}
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
