package valids

import (
	"fmt"
	"reflect"

	"helpers.zhaowenming.cn/maps"
)

//NumberValidator 数值验证器
//
//验证数值的大小区间，根据验证方式不同需要设置不同类型的区间
//
//min和max在默认情况为int或float64类型，通过给出的数值自动判断
type NumberValidator struct {
	kind     int         //验证方式：0默认整数，1浮点数
	Min      interface{} //最小，含，自动转换为int或者float64，其他类型均会报错
	Max      interface{} //最大，含
	SmallMsg string      //小了的时候提醒
	BigMsg   string      //大了的时候提醒
	Message  string      //默认错误消息

	*Validator
}

func NewNumberValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	nv := NumberValidator{}
	nv.Validator = v

	err := maps.SetMapToStruct(params, &nv)
	if err != nil {
		panic("NewNumberValidator err=" + err.Error())
	}

	return &nv
}

func (n *NumberValidator) Init() {
	//父级初始化
	n.Validator.Init()

	min1, ok1 := n.Min.(int)
	min2, ok2 := n.Min.(float64)
	max1, ok3 := n.Max.(int)
	max2, ok4 := n.Max.(float64)

	if !ok1 && !ok2 || !ok3 && !ok4 {
		panic("NewNumberValidator err= min、max只支持int或者float64类型")
	}

	if ok1 && ok4 {
		//转换min1
		n.Min = float64(min1)
		if max2 < float64(min1) {
			panic("NewNumberValidator err= min不能比max大")
		}

		n.kind = 1
	}
	if ok2 && ok3 {
		//转换max1
		n.Max = float64(max1)
		if float64(max1) < min2 {
			panic("NewNumberValidator err= min不能比max大")
		}

		n.kind = 1
	}

	if n.Message == "" {
		n.Message = "{attribute} 必须是数字"
	}
	if n.BigMsg == "" {
		n.BigMsg = "{attribute} 数值不能大于等于 {max}"
	}
	if n.SmallMsg == "" {
		n.SmallMsg = "{attribute} 数值不能小于 {min}"
	}
}

func (n *NumberValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	formatParams = make(map[string]string)
	if n.kind == 0 {
		formatParams["min"] = fmt.Sprintf("%d", n.Min)
		formatParams["max"] = fmt.Sprintf("%d", n.Max)
	} else {
		formatParams["min"] = fmt.Sprintf("%f", n.Min)
		formatParams["max"] = fmt.Sprintf("%f", n.Max)
	}

	rt := reflect.TypeOf(v)
	sk := rt.Kind()
	sv := reflect.ValueOf(v)
	switch sk {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n.kind == 0 {
			value := int(sv.Int())
			if value < n.Min.(int) {
				return n.SmallMsg, formatParams
			}
			if value > n.Max.(int) {
				return n.BigMsg, formatParams
			}
		} else {
			value := float64(sv.Int())
			if value < n.Min.(float64) {
				return n.SmallMsg, formatParams
			}
			if value > n.Max.(float64) {
				return n.BigMsg, formatParams
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n.kind == 0 {
			value := int(sv.Uint())
			if value < n.Min.(int) {
				return n.SmallMsg, formatParams
			}
			if value > n.Max.(int) {
				return n.BigMsg, formatParams
			}
		} else {
			value := float64(sv.Uint())
			if value < n.Min.(float64) {
				return n.SmallMsg, formatParams
			}
			if value > n.Max.(float64) {
				return n.BigMsg, formatParams
			}
		}
	case reflect.Float32, reflect.Float64:
		value := sv.Float()
		if value < n.Min.(float64) {
			return n.SmallMsg, formatParams
		}
		if value > n.Max.(float64) {
			return n.BigMsg, formatParams
		}
	default:
		return n.Message, formatParams
	}

	return "", nil
}
