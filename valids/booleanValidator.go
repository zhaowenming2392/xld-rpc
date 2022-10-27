package valids

import (
	"fmt"
)

//BooleanValidator 布尔验证器
//
//Kind 0默认，验证真假；1通过数字如1、0来验证真假；2通过字符串如'1'、'0'来判断真假；3通过任何值来判断
//
//非0的Kind，可以自定义具体真假对应的值，如不一定非要用1、0，也可以用10、20
type BooleanValidator struct {
	Kind     int         //初始话时判断，采用什么类型进行判断
	TrueInt  int         //作为数字类型进行判断时，true的值
	FalseInt int         //作为数字类型进行判断时，false的值
	TrueStr  string      //作为字符串类型进行判断时，true的值
	FalseStr string      //作为字符串类型进行判断时，false的值
	TrueOt   interface{} //作为任意值判断时的真值，类型必须相同，部分类型不可比较，如map，[]切片等动态类型的值
	FalseOt  interface{} //作为任意值判断时的假值

	*Validator
}

func NewBooleanValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	bv := BooleanValidator{}
	bv.Validator = v

	if v, ok := params["message"]; ok {
		bv.Message = v.(string)
	}
	if v, ok := params["kind"]; ok {
		bv.Kind = v.(int)
	}
	if v, ok := params["trueInt"]; ok {
		bv.TrueInt = v.(int)
	}
	if v, ok := params["falseInt"]; ok {
		bv.FalseInt = v.(int)
	}
	if v, ok := params["trueStr"]; ok {
		bv.TrueStr = v.(string)
	}
	if v, ok := params["falseStr"]; ok {
		bv.FalseStr = v.(string)
	}
	if v, ok := params["trueOt"]; ok {
		bv.TrueOt = v
	}
	if v, ok := params["falseOt"]; ok {
		bv.FalseOt = v
	}

	return &bv
}

//Init 初始化，注意是在对具体验证器进行设置之后的初始化
func (bv *BooleanValidator) Init() {
	//父级初始化
	bv.Validator.Init()

	//初始化
	if bv.Message == "" {
		bv.Message = "{attribute} 必须是 \"{true}\" 或者 \"{false}\""
	}

	//没有指定类型，直接使用默认
	switch bv.Kind {
	case 0: //默认，按照bool
	case 1: //指定整数
		if bv.TrueInt == 0 {
			bv.TrueInt = 1
		}
		if bv.FalseInt == bv.TrueInt {
			panic("\"{true}\" 和 \"{false}\"的真假值不能相同")
		}
	case 2:
		if bv.TrueStr == "" {
			bv.TrueStr = "1"
		}
		if bv.FalseStr == "" {
			bv.FalseStr = "0"
		}
		if bv.FalseStr == bv.TrueStr {
			panic("\"{true}\" 和 \"{false}\"的真假值不能相同")
		}
	case 3:
		if bv.TrueOt == bv.FalseOt {
			panic("\"{true}\" 和 \"{false}\"的真假值不能相同")
		}
	}
}

func (bv *BooleanValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	formatParams = make(map[string]string)

	switch bv.Kind {
	case 0: //默认，按照bool
		formatParams["true"] = "true"
		formatParams["false"] = "false"
	case 1: //指定整数
		formatParams["true"] = fmt.Sprintf("%d", bv.TrueInt)
		formatParams["false"] = fmt.Sprintf("%d", bv.FalseInt)
	case 2:
		formatParams["true"] = bv.TrueStr
		formatParams["false"] = bv.FalseStr
	case 3:
		formatParams["true"] = fmt.Sprintf("%s", bv.TrueOt)
		formatParams["false"] = fmt.Sprintf("%s", bv.FalseOt)
	}

	if v == nil {
		return bv.Message, formatParams
	}

	switch bv.Kind {
	case 0:
		if _, ok := v.(bool); ok {
			return "", nil
		}
	case 1:
		if sv, ok := v.(int); ok {
			if sv == bv.TrueInt || sv == bv.FalseInt {
				return "", nil
			}
		}
	case 2:
		if sv, ok := v.(string); ok {
			if sv == bv.TrueStr || sv == bv.FalseStr {
				return "", nil
			}
		}
	case 3:
		if v == bv.TrueOt || v == bv.FalseOt {
			return "", nil
		}
	}

	return bv.Message, formatParams
}
