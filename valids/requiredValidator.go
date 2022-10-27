package valids

import (
	"fmt"
	"strings"

	"helpers.zhaowenming.cn/types"
)

//RequiredValidator
type RequiredValidator struct {
	//属性必须具有的所需值。
	//如果为空，验证程序将验证指定的属性是否为空。默认RequiredValue为空。
	//如果将其设置为非空值，验证程序将验证该属性是否具有与该属性值相同的值。
	RequiredValue interface{}
	*Validator
}

func NewRequiredValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	rv := RequiredValidator{
		Validator: v,
	}
	rv.Validator.SkipEmpty = types.BoolPtr(false)

	if v, ok := params["requiredValue"]; ok {
		rv.RequiredValue = v
	}

	if v, ok := params["message"]; ok {
		rv.Message = v.(string)
	} else {
		if rv.RequiredValue == nil {
			rv.Message = "{attribute} 不能为空"
		} else {
			rv.Message = "{attribute} 必须是 {requiredValue}"
		}
	}

	return &rv
}

func (rv *RequiredValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	if rv.RequiredValue == nil {
		//字符串，首尾去除空字符串
		str, ok := v.(string)
		if ok {
			v = strings.Trim(str, " ")
		}

		if !rv.IsEmpty(v) {
			return "", nil
		}
	} else if v == rv.RequiredValue {
		return "", nil
	}

	if rv.RequiredValue == nil {
		return rv.Message, nil
	}

	return rv.Message, map[string]string{"requiredValue": fmt.Sprintf("%v", rv.RequiredValue)}
}
