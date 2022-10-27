package valids

import (
	"regexp"

	"helpers.zhaowenming.cn/maps"
)

//MatchValidator 正则验证器
//
//验证是否满足正则表达式，如果not=true，则验证是否不满足
type MatchValidator struct {
	Pattern string //正则
	Not     bool   //进行匹配或不匹配

	*Validator
}

func NewMatchValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	mv := MatchValidator{}
	mv.Validator = v

	err := maps.SetMapToStruct(params, &mv)
	if err != nil {
		panic("NewMatchValidator err=" + err.Error())
	}

	return &mv
}

//Init 初始化，注意是在对具体验证器进行设置之后的初始化
func (mv *MatchValidator) Init() {
	//父级初始化
	mv.Validator.Init()

	if mv.Pattern == "" {
		panic("正则表达式必须提供")
	}

	//初始化
	if mv.Message == "" {
		mv.Message = "{attribute} 格式无效！"
	}
}

//ValidateValue 验证某个值
func (mv *MatchValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	reg := regexp.MustCompile(mv.Pattern)

	var ok bool
	switch value := v.(type) {
	case string:
		ok = reg.MatchString(value)
	case []byte:
		ok = reg.Match(value)
	default:
		return "{attribute} 只能是字符串或字节切片！", nil
	}

	ok = ok && !mv.Not || !ok && mv.Not

	if !ok {
		return mv.Message, nil
	}
	return "", nil
}
