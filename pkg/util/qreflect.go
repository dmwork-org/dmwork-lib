package util

import (
	"errors"
	"reflect"
)

// AttrToUnderscore 获取struct的所有属性并转为下划线模式
func AttrToUnderscore(st interface{}) ([]string, error) {
	if st == nil {
		return nil, errors.New("input cannot be nil")
	}
	t := reflect.ValueOf(st)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("input must be a pointer to struct")
	}
	if t.IsNil() {
		return nil, errors.New("input cannot be nil pointer")
	}
	vType := t.Elem().Type()
	if vType.Kind() != reflect.Struct {
		return nil, errors.New("input must be a pointer to struct")
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
	return names, nil
}
