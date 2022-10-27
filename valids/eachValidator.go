package valids

import (
	"fmt"
	"reflect"

	"helpers.zhaowenming.cn/maps"
)

type EachValidator struct {
	Rule   string
	Params map[string]interface{}

	AllowMessageFromRule bool
	StopOnFirstError     bool

	*Validator
	valid ValidatorInterface
}

func NewEachValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	ev := EachValidator{}
	ev.Validator = v
	ev.AllowMessageFromRule = true
	ev.StopOnFirstError = true

	if r, ok := params["rule"]; !ok || r == "" {
		panic("NewEachValidator err=rule参数必须提供")
	}

	err := maps.SetMapToStruct(params, &ev)
	if err != nil {
		panic("NewEachValidator err=" + err.Error())
	}

	return &ev
}

//Init 初始化，注意是在对具体验证器进行设置之后的初始化
func (ev *EachValidator) Init() {
	//父级初始化
	ev.Validator.Init()

	//初始化
	if ev.Message == "" {
		ev.Message = "{attribute} 的格式无效！"
	}
}

func (ev *EachValidator) getValidator() ValidatorInterface {
	if ev.valid == nil {
		if valid, ok := BuiltInValidators[ev.Rule]; ok {
			//创建具体验证器
			ev.valid = valid(ev.Params, ev.Validator)
			//初始化验证器
			ev.valid.Init()

			fmt.Printf("验证器%s： %+v 已经创建！\n", ev.Rule, ev.Validator)
		}

		panic("不存在\"" + ev.Rule + "\"内置验证器！")
	}

	return ev.valid
}

func (ev *EachValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice {
		valid := ev.getValidator()
		for i := 0; i < rv.Len(); i++ {
			value := rv.Index(i).Interface()
			if *ev.SkipEmpty && ev.IsEmpty(value) {
				continue
			}

			msg, ps := valid.ValidateValue(value)
			if msg != "" {
				if ev.AllowMessageFromRule {
					if ps == nil {
						ps = make(map[string]string)
					}
					ps["value"] = fmt.Sprintf("%s", value)
					return msg, ps
				}

				return ev.Message, map[string]string{
					"value": fmt.Sprintf("%s", value),
				}
			}
		}
	} else if rt.Kind() == reflect.Map {
		valid := ev.getValidator()
		for _, mv := range rv.MapKeys() {
			value := rv.MapIndex(mv).Interface()
			if *ev.SkipEmpty && ev.IsEmpty(value) {
				continue
			}

			msg, ps := valid.ValidateValue(value)
			if msg != "" {
				if ev.AllowMessageFromRule {
					if ps == nil {
						ps = make(map[string]string)
					}
					ps["value"] = fmt.Sprintf("%s", value)
					return msg, ps
				}

				return ev.Message, map[string]string{
					"value": fmt.Sprintf("%s", value),
				}
			}
		}
	} else {
		return ev.Message, nil
	}

	return "", nil
}

func (ev *EachValidator) ValidateAttribute(m ModelInterface, f string) {
	v := m.GetAttributeValue(f)
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice {
		valid := ev.getValidator()
		for i := 0; i < rv.Len(); i++ {
			value := rv.Index(i).Interface()
			if *ev.SkipEmpty && ev.IsEmpty(value) {
				continue
			}

			msg, ps := valid.ValidateValue(value)
			if msg != "" {
				if ev.AllowMessageFromRule {
					if ps == nil {
						ps = make(map[string]string)
					}
					ps["value"] = fmt.Sprintf("%s", value)
					m.AddError(f, ev.formatMessage(msg, ps))
				} else {
					m.AddError(f, ev.formatMessage(ev.Message, map[string]string{
						"value": fmt.Sprintf("%s", value),
					}))
				}

				if ev.StopOnFirstError {
					break
				}
			}
		}
	} else if rt.Kind() == reflect.Map {
		valid := ev.getValidator()
		for _, mv := range rv.MapKeys() {
			value := rv.MapIndex(mv).Interface()
			if *ev.SkipEmpty && ev.IsEmpty(value) {
				continue
			}

			msg, ps := valid.ValidateValue(value)
			if msg != "" {
				if ev.AllowMessageFromRule {
					if ps == nil {
						ps = make(map[string]string)
					}
					ps["value"] = fmt.Sprintf("%s", value)
					m.AddError(f, ev.formatMessage(msg, ps))
				} else {
					m.AddError(f, ev.formatMessage(ev.Message, map[string]string{
						"value": fmt.Sprintf("%s", value),
					}))
				}

				if ev.StopOnFirstError {
					break
				}
			}
		}
	} else {
		m.AddError(f, ev.formatMessage(ev.Message, nil))
	}
}
