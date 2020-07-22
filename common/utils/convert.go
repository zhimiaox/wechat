/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

import (
	"reflect"
)

// 使用反射，转换结构体 仅支持单层级
func SuperConvert(fromStruct interface{}, toStruct interface{}) {
	fromStructMap := structToMap(fromStruct)
	toStructV := reflect.ValueOf(toStruct).Elem()
	toStructT := reflect.TypeOf(toStruct).Elem()
	for i := 0; i < toStructV.NumField(); i++ {
		fieldName := toStructT.Field(i).Name
		if sourceVal, ok := fromStructMap[fieldName]; ok {
			if !sourceVal.IsValid() {
				continue
			}
			toStructVal := toStructV.Field(i)
			if toStructVal.CanSet() {
				if sourceVal.Type().PkgPath()+sourceVal.Type().Name() == toStructVal.Type().PkgPath()+toStructVal.Type().Name() {
					toStructVal.Set(sourceVal)
				}
			}
		}
	}
}

func structToMap(structName interface{}) map[string]reflect.Value {
	t := reflect.TypeOf(structName).Elem()
	v := reflect.ValueOf(structName).Elem()
	fieldNum := t.NumField()
	resMap := make(map[string]reflect.Value, fieldNum)
	for i := 0; i < fieldNum; i++ {
		resMap[t.Field(i).Name] = v.Field(i)
	}
	return resMap
}
