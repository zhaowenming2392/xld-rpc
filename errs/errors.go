package errs

import (
	"errors"
	"fmt"
)

//创建新错误
func NewErr(msg string) error {
	return errors.New(msg)
}

//不能为空
func EmptyErr(name string) error {
	return NewErr(name + " 不能为空！")
}

//不存在
func NotExistErr(name string) error {
	return NewErr(name + " 不存在！")
}

//已经存在
func AlreadyExistErr(name string) error {
	return NewErr(name + " 已经存在！")
}

//%s 必须在 %v 中
func RangeErr(name string, rag ...interface{}) error {
	return NewErr(fmt.Sprintf("%s 必须在 %v 中", name, rag))
}

//%s 必须在 %d - %d（不含）之间
func LimitErr(name string, min, max int) error {
	return NewErr(fmt.Sprintf("%s 必须在 %d - %d（不含）之间", name, min, max))
}
