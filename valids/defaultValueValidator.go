package valids

import (
	"helpers.zhaowenming.cn/types"
)

//DefaultValueValidator 默认值验证器
//
//当属性空时，给属性设置默认值
type DefaultValueValidator struct {
	Value interface{}
	*Validator
}

func NewDefaultValueValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	dv := DefaultValueValidator{
		Validator: v,
	}

	dv.Validator.SkipEmpty = types.BoolPtr(false)

	if v, ok := params["value"]; ok {
		dv.Value = v
	}

	return &dv
}

func (rv *DefaultValueValidator) ValidateAttribute(m ModelInterface, f string) {
	if rv.IsEmpty(m.GetAttributeValue(f)) {
		m.SetAttributeValue(f, rv.Value)
	}
}
