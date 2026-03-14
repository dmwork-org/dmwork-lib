package util

import (
	"reflect"
)

// AttrToUnderscore 获取struct的所有属性并转为下划线模式
// 参数必须为指向结构体的非nil指针，否则返回空切片
func AttrToUnderscore(st interface{}) []string {
	if st == nil {
		return []string{}
	}
	t := reflect.ValueOf(st)
	if t.Kind() != reflect.Ptr || t.IsNil() {
		return []string{}
	}
	vType := t.Elem().Type()
	if vType.Kind() != reflect.Struct {
		return []string{}
	}
	names := make([]string, 0)
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
