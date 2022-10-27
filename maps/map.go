package maps

import (
	"errors"
	"fmt"
	"reflect"

	"helpers.zhaowenming.cn/strs"
)

//Equal 比较两个映射是否相同
func Equal(x, y map[interface{}]interface{}) bool {
	//地址一样肯定相同
	if &x == &y {
		return true
	}
	if len(x) != len(y) {
		return false
	}
	//一一比对
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}

//SetMapToStruct 将映射中的值设置为结构体中对应的属性值
//
//params 键和结构体的字段对应，以首字母小写形式，只有匹配上的才会进行设置
//
//obj 必须是结构体指针，非指针不能进行设置
func SetMapToStruct(params map[string]interface{}, obj interface{}) (err error) {
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Ptr {
		return errors.New("目标对象必须是结构体指针")
	}
	rt = rt.Elem()
	if rt.Kind() != reflect.Struct {
		return errors.New("目标对象必须是结构体指针")
	}

	//设置recover，recover只能放在defer后面使用才有效恢复
	defer func() {
		//恢复宕机，并捕获宕机抛出的任意类型数据
		if msg := recover(); msg != nil {
			//出错了，才进行处理
			err = fmt.Errorf("%s", msg) //修改错误信息
		}
	}() //必须调用

	paramsLen := len(params)
	rv := reflect.ValueOf(obj).Elem()
	for i := 0; i < rt.NumField(); i++ {
		//如果已经把所有参数都比对过了，则结束比对
		if paramsLen == 0 {
			break
		}

		fieldType := rt.Field(i)
		fieldName := fieldType.Name
		fieldValue := rv.Field(i)
		fieldName2 := strs.FormatName(fieldName, 1)

		for paramName, paramValue := range params {
			paramName = strs.FormatName(paramName, 1)
			if paramName == fieldName2 {
				//找到一个就少一个
				paramsLen--

				if fieldValue.CanSet() {
					paramValueType := reflect.TypeOf(paramValue)
					if fieldType.Type.Kind() != paramValueType.Kind() && fieldType.Type.Kind() != reflect.Interface {
						return errors.New(paramName + "值(" + fieldType.Type.Kind().String() + ")和结构体属性" + fieldName + "值(" + paramValueType.Kind().String() + ")类型不一致")
					}

					//TODO 通用设置（效率可能没有以下具体类型设置高），如果类型不一致会宕机
					fieldValue.Set(reflect.ValueOf(paramValue))

					/*
						//特殊类型和无效类型直接报错
						switch fieldType.Type.Kind() {
						case reflect.Bool:
							fieldValue.SetBool(paramValue.(bool))
						case reflect.Int:
							fieldValue.SetInt(int64(paramValue.(int)))
						case reflect.Int8:
							fieldValue.SetInt(int64(paramValue.(int8)))
						case reflect.Int16:
							fieldValue.SetInt(int64(paramValue.(int16)))
						case reflect.Int32:
							fieldValue.SetInt(int64(paramValue.(int32)))
						case reflect.Int64:
							fieldValue.SetInt(int64(paramValue.(int64)))
						case reflect.Uint:
							fieldValue.SetUint(uint64(paramValue.(uint)))
						case reflect.Uint8:
							fieldValue.SetUint(uint64(paramValue.(uint8)))
						case reflect.Uint16:
							fieldValue.SetUint(uint64(paramValue.(uint16)))
						case reflect.Uint32:
							fieldValue.SetUint(uint64(paramValue.(uint32)))
						case reflect.Uint64:
							fieldValue.SetUint(uint64(paramValue.(uint64)))
						case reflect.Float32:
							fieldValue.SetFloat(float64(paramValue.(float32)))
						case reflect.Float64:
							fieldValue.SetFloat(float64(paramValue.(float64)))
						default:
							fieldValue.Set(reflect.ValueOf(paramValue))
						}
					*/
				} else {
					fmt.Println("fieldValue.Set 不能设置 ", fieldName)
				}

				break
			}
		}
	}

	return nil
}
