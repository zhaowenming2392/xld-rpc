//定义一些常量、常用类型、函数等
package types

import (
	"math"
	"reflect"
	"time"
)

// Time 复制 time.Time 对象，并返回复制体的指针
func TimePtr(t time.Time) *time.Time {
	return &t
}

// String 复制 string 对象，并返回复制体的指针
func StringPtr(s string) *string {
	return &s
}

// Bool 复制 bool 对象，并返回复制体的指针
func BoolPtr(b bool) *bool {
	return &b
}

// Float64 复制 float64 对象，并返回复制体的指针
func Float64Ptr(f float64) *float64 {
	return &f
}

// Float32 复制 float32 对象，并返回复制体的指针
func Float32Ptr(f float32) *float32 {
	return &f
}

// Int64 复制 int64 对象，并返回复制体的指针
func Int64Ptr(i int64) *int64 {
	return &i
}

// Int32 复制 int64 对象，并返回复制体的指针
func Int32Ptr(i int32) *int32 {
	return &i
}

//IsInteger 判断数字是否是整数
func IsInteger(f float64) bool {
	//取得浮点数的整数部分与之比较
	return f == math.Trunc(f)
}

//IsArray 是否是数组
func IsArray(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Array
}

//IsSlice 是否是切片
func IsSlice(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Slice
}

//IsMap 是否是映射
func IsMap(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Map
}

//IsFunc 是否是函数
func IsFunc(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Func
}

//IsStruct 是否是结构体
func IsStruct(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Struct
}

/*
type TypeAndValue struct {
	TypeName   string//通过类型名称来获取值
	BoolValue  bool
	IntValue   int64
	FloatValue float64
	StrValue   string
	//---
}

//获取一个v的类型和值
func GetType(v interface{}) (tv TypeAndValue) {
	rt := reflect.TypeOf(v)
	sk := rt.Kind()
	sv := reflect.ValueOf(v)
	switch sk {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tv.IntValue = sv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tv.IntValue = int64(sv.Uint())

	case reflect.Float32, reflect.Float64:
		tv.FloatValue = sv.Float()

	default:
		return fmt.Errorf(n.Message, v)
	}

	return nil
}
*/
