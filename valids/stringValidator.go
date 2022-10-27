package valids

import (
	"strconv"
	"unicode/utf8"

	"helpers.zhaowenming.cn/maps"
)

type StringValidator struct {
	ValidByte bool   //是否验证字节长度，默认验证字符个数
	Min       int    //最短，含
	Max       int    //最长
	Length    []int  //短、长，切片形式提供，如果不为nil会替换到min max
	ShortMsg  string //短了时候的提醒
	LongMsg   string //过长时候的提醒

	*Validator
}

func NewStringValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	sv := StringValidator{}
	sv.Validator = v

	err := maps.SetMapToStruct(params, &sv)
	if err != nil {
		panic("NewStringValidator err=" + err.Error())
	}

	return &sv
}

func (s *StringValidator) Init() {
	//父级初始化
	s.Validator.Init()

	if s.Length != nil {
		if len(s.Length) != 2 || s.Length[0] > s.Length[1] {
			panic("如果提供Length，则必须为一小一大两个参数")
		} else {
			s.Min = s.Length[0]
			s.Max = s.Length[1]
		}
	}

	if s.Message == "" {
		s.Message = "{attribute} 必须是字符串"
	}

	if s.ShortMsg == "" {
		s.ShortMsg = "{attribute} 不能少于 {min} 个字"
	}

	if s.LongMsg == "" {
		s.LongMsg = "{attribute} 不能多于 {max} 个字"
	}
}

func (s *StringValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	formatParams = make(map[string]string)
	formatParams["min"] = strconv.Itoa(s.Min)
	formatParams["max"] = strconv.Itoa(s.Max)

	vv, ok := v.(string)
	if !ok {
		return s.Message, formatParams
	}

	if s.ValidByte {
		if len(vv) < s.Min {
			return s.ShortMsg + "节", formatParams
		}

		if len(vv) > s.Max {
			return s.LongMsg + "节", formatParams
		}
	} else {
		vl := utf8.RuneCountInString(vv)
		if vl < s.Min {
			return s.ShortMsg, formatParams
		}

		if vl > s.Max {
			return s.LongMsg, formatParams
		}
	}

	return "", nil
}
