package valids

import (
	"fmt"
)

//FuncValidator 自定义函数验证器
//
//使用时必须提供自定义的func参数或modeFunc参数
type FuncValidator struct {
	f1 func(v interface{}) string
	f2 func(mode ModelInterface, v interface{}) string
	*Validator
}

//创建自定义函数验证器
func NewFuncValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	var f1 func(v interface{}) string
	var f2 func(mode ModelInterface, v interface{}) string

	if f, ok := params["func"].(func(v interface{}) string); ok {
		f1 = f
	}
	if f, ok := params["modeFunc"].(func(mode ModelInterface, v interface{}) string); ok {
		f2 = f
	}

	if f1 == nil && f2 == nil {
		panic("必须提供满足格式要求的自定义的func参数或modeFunc参数")
	}

	return &FuncValidator{
		f1:        f1,
		f2:        f2,
		Validator: v,
	}
}

//ValidateAttribute 验证属性，发送错误则在其模型中添加错误
func (fv FuncValidator) ValidateAttribute(m ModelInterface, f string) {
	if fv.f2 != nil {
		//使用自定义验证器来验证
		msg := fv.f2(m, f)

		if msg != "" {
			m.AddError(f, msg)
		}

		fmt.Printf("验证结果为%s = %s \n", f, msg)
		return
	}

	//降级，采用继承的，此时必须要有func参数
	fv.Validator.ValidateAttribute(m, f)
}

//将自定义函数转换为标准验证器
func (fv FuncValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	fmt.Println("********FuncValidator***********")

	if fv.f1 != nil {
		return fv.f1(v), nil
	}

	panic("当前没有提供满足格式要求的自定义的func参数，不能运行ValidateValue方法")
}
