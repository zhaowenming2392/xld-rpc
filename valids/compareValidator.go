package valids

import (
	"fmt"

	"helpers.zhaowenming.cn/maps"
)

//CompareValidator 比较验证器
//
//Kind 0默认，比较字符串；1比较数字(int，float64)
//
//默认和CompareValue比较，CompareAttribute和属性比较
type CompareValidator struct {
	Kind int //初始话时判断，采用什么类型进行比较

	CompareAttribute string      //和某个属性进行比较
	CompareValue     interface{} //和某个值进行比较
	Operator         string      //比较操作符

	*Validator
}

func NewCompareValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	cv := CompareValidator{}
	cv.Validator = v
	cv.Operator = "=="

	err := maps.SetMapToStruct(params, &cv)
	if err != nil {
		panic("NewNumberValidator err=" + err.Error())
	}

	return &cv
}

//Init 初始化，注意是在对具体验证器进行设置之后的初始化
func (cv *CompareValidator) Init() {
	//父级初始化
	cv.Validator.Init()

	//初始化
	if cv.Message == "" {
		//没有指定类型，直接使用默认
		switch cv.Operator {
		case "==": //默认
			cv.Message = "{attribute} 必须等于 \"{compareValueOrAttribute}\"."
		case "!=":
			cv.Message = "{attribute} 不能等于 \"{compareValueOrAttribute}\"."
		case ">":
			cv.Message = "{attribute} 必须大于 \"{compareValueOrAttribute}\"."
		case ">=":
			cv.Message = "{attribute} 必须大于等于 \"{compareValueOrAttribute}\"."
		case "<":
			cv.Message = "{attribute} 必须小于 \"{compareValueOrAttribute}\"."
		case "<=":
			cv.Message = "{attribute} 必须小于等于 \"{compareValueOrAttribute}\"."
		default:
			panic("比较操作符 " + cv.Operator + " 不支持！")
		}
	}
}

func (cv *CompareValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	if cv.CompareValue == nil {
		panic("比较的值 compareValue 必须提供！")
	}
	if !cv.compareValues(v, cv.CompareValue) {
		return cv.Message, map[string]string{
			"compareValueOrAttribute": fmt.Sprintf("%s", cv.CompareValue),
		}
	}

	return "", nil
}

func (cv *CompareValidator) ValidateAttribute(m ModelInterface, f string) {
	value := m.GetAttributeValue(f)

	var compareAttribute string
	var compareValueOrAttribute string
	var compareValue interface{}

	if cv.CompareValue != nil {
		compareValueOrAttribute = fmt.Sprintf("%s", cv.CompareValue)
		compareValue = cv.CompareValue
	} else {
		if cv.CompareAttribute == "" {
			compareAttribute = f + "2"
		} else {
			compareAttribute = cv.CompareAttribute
		}
		compareValue = m.GetAttributeValue(compareAttribute)
		compareValueOrAttribute = m.GetAttributeLabel(compareAttribute)
	}

	if !cv.compareValues(value, compareValue) {
		m.AddError(f, cv.formatMessage(cv.Message, map[string]string{
			"compareValueOrAttribute": compareValueOrAttribute,
		}))
	}
}

func (cv *CompareValidator) compareValues(value, compareValue interface{}) bool {
	if cv.Kind == 0 {
		v1 := value.(string)
		v2 := compareValue.(string)

		switch cv.Operator {
		case "==": //默认
			return v1 == v2
		case "!=":
			return v1 != v2
		case ">":
			return v1 > v2
		case ">=":
			return v1 >= v2
		case "<":
			return v1 < v2
		case "<=":
			return v1 <= v2
		default:
			return false
		}
	} else {
		v1, ok1 := value.(int)
		if ok1 {
			v11, ok11 := compareValue.(int)
			if !ok11 {
				panic("两个比较值类型不相同")
			}

			switch cv.Operator {
			case "==": //默认
				return v1 == v11
			case "!=":
				return v1 != v11
			case ">":
				return v1 > v11
			case ">=":
				return v1 >= v11
			case "<":
				return v1 < v11
			case "<=":
				return v1 <= v11
			default:
				return false
			}
		}

		v2, ok2 := value.(float64)
		if ok2 {
			v22, ok22 := compareValue.(float64)
			if !ok22 {
				panic("两个比较值类型不相同")
			}

			switch cv.Operator {
			case "==": //默认
				return v2 == v22
			case "!=":
				return v2 != v22
			case ">":
				return v2 > v22
			case ">=":
				return v2 >= v22
			case "<":
				return v2 < v22
			case "<=":
				return v2 <= v22
			default:
				return false
			}
		}

		panic("两个比较值的数值类型不支持，仅支持int和float64类型数值")
	}
}
