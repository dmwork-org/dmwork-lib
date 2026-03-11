package util

import (
	"reflect"
)

// AttrToUnderscore 获取struct的所有属性并转为下划线模式
func AttrToUnderscore(st interface{}) []string {
	if st == nil {
		return nil
	}
	t := reflect.ValueOf(st)
	if t.Kind() != reflect.Ptr {
		return nil
	}
	if t.IsNil() {
		return nil
	}
	vType := t.Elem().Type()
	if vType.Kind() != reflect.Struct {
		return nil
	}
	names := make([]string, 0, vType.NumField())
	for i := 0; i < vType.NumField(); i++ {
		if vType.Field(i).Type.Kind() == reflect.Struct {
			continue
		}
		name := vType.Field(i).Name
		if name != "" {
			names = append(names, UnderscoreName(name))
		}
	}
	return names
}
