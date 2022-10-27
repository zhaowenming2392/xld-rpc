package slice

import (
	"errors"
)

type AnySlice []interface{}

//ToAnySlice 将多个不同类型的参数追加到一个任意值类型的切片中
func ToAnySlice(values ...interface{}) []interface{} {
	newSlice := []interface{}{}
	return append(newSlice, values...)
}

//Unique 返回去重后的新切片，不损害原切片
func Unique(oldSlice interface{}) (newSlice []interface{}) {
	os := oldSlice.(AnySlice)
	s := make(map[interface{}]bool)
	for _, v := range os {
		if _, ok := s[v]; !ok {
			s[v] = true
			newSlice = append(newSlice, v)
		}
	}

	return
}

//UniqueStrs 返回去重后的新string切片，不损害原切片
func UniqueStrs(oldSlice []string) (newSlice []string) {
	s := make(map[string]bool)
	for _, v := range oldSlice {
		if _, ok := s[v]; !ok {
			s[v] = true
			newSlice = append(newSlice, v)
		}
	}

	return
}

//UniqueInts 返回去重后的新int切片，不损害原切片
func UniqueInts(oldSlice []int) (newSlice []int) {
	s := make(map[int]bool)
	for _, v := range oldSlice {
		if _, ok := s[v]; !ok {
			s[v] = true
			newSlice = append(newSlice, v)
		}
	}

	return
}

//Combine 返回一个 切片，用来自 keys 切片的值作为键名，来自 values 切片的值作为相应的值。
//注意只能都是 []int 类型
func Combine(keys, values []int) (newSlice []int) {
	for k := range keys {
		newSlice[k] = values[k]
	}
	return
}

//Merge 将一个或多个切片的单元合并起来，一个切片中的值附加在前一个切片的后面
func Merge(slices ...interface{}) (newSlice []interface{}) {
	for _, v := range slices {
		vs := v.(AnySlice)
		newSlice = append(newSlice, vs...)
	}
	return
}

//Equal 比较两个切片是否相同
func Equal(sliceX, sliceY interface{}) bool {
	x := sliceX.(AnySlice)
	y := sliceY.(AnySlice)
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

//IsEmpty 判断切片是不是空的
func IsEmpty(slice interface{}) bool {
	s := slice.(AnySlice)
	if s == nil {
		return true
	}
	if len(s) == 0 {
		return true
	}
	return false
}

//Shift 将切片开头的x个单元移出数组，保留x下标开始的元素
func Shift(slice interface{}, x int) []interface{} {
	s := slice.(AnySlice)
	return s[x:]
}

//Pop 弹出切片最后x个单元（出栈）
func Pop(slice interface{}, x int) []interface{} {
	s := slice.(AnySlice)
	return s[:len(s)-x]
}

//Remove 删除切片的某个元素，并且保证原来的排序
func Remove(slice interface{}, i int) []interface{} {
	s := slice.(AnySlice)
	copy(s[i:], s[i+1:])
	return s[:len(s)-1]
}

//Reverse 反转切片,返回单元顺序相反的切片
func Reverse(slice interface{}) {
	s := slice.(AnySlice)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

//InSlice 判断某个值是否在切片中，并返回其下标
func InSlice(slice, v interface{}) (int, error) {
	s := slice.(AnySlice)
	for sk, sv := range s {
		if sv == v {
			return sk, nil
		}
	}

	return 0, errors.New("未找到")
}
