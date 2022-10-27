package valids

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"helpers.zhaowenming.cn/strs"
	"helpers.zhaowenming.cn/types"
)

//验证器接口
//
//验证器只是自己去验证某个值
type ValidatorInterface interface {
	Init()                                            //初始化，注意是在对具体验证器进行设置之后的初始化
	ValidateAttributes(m ModelInterface, fs []string) //验证模型的一组属性，具体类型验证器可以不实现本方法，采用嵌入Validator来继承
	ValidateAttribute(m ModelInterface, f string)     //验证模型的属性，具体类型验证器可以不实现本方法，采用嵌入Validator来继承

	//TODO 具体验证器必须实现ValidateValue方法
	ValidateValue(v interface{}) (msg string, formatParams map[string]string) //验证值，如果验证成功，msg为空，否则msg是具有格式的错误消息，利用formatParams可以组织成最终的消息字符串
}

//内置验证器组
//
//由具体验证器实现以下创建函数再来添加
//
//func NewXxxValidator(params map[string]interface{}) ValidatorInterface {
//	进行本验证器自己的参数初始化等
//}
var BuiltInValidators = map[string]func(params map[string]interface{}, v *Validator) ValidatorInterface{
	"func":     NewFuncValidator,         //自定义函数验证器
	"boolean":  NewBooleanValidator,      //布尔验证器
	"default":  NewDefaultValueValidator, //设置默认值
	"required": NewRequiredValidator,     //不能为空
	"str":      NewStringValidator,       //字符串长短验证
	"num":      NewNumberValidator,       //数字大小验证
}

//规格及其验证器
//
//由其调用和创建具体验证器来进行验证，请在Mode的Rules()返回[]*Validator来使用
type Validator struct {
	Attributes []string                              //验证器中待验证的字段
	On         []string                              //场景，可以应用的场景
	Except     []string                              //当前验证程序不应应用于的场景，和On二选一
	SkipError  *bool                                 //如果正在验证的属性已存在某些验证错误，则是否应跳过此验证规则。
	SkipEmpty  *bool                                 //如果正在验证的属性是空时，则是否应跳过此验证规则。
	IsEmpty    func(v interface{}) bool              //空的判断，自定义判断当前值是否为空的函数，如果未定义按照验证器的默认实现
	When       func(m interface{}, attr string) bool //条件验证，当此函数返回true才进行验证，否则跳过本次验证
	Message    string                                //默认错误消息

	Name         string                 //具体验证器名称，空则采用自定义函数，Params中必须包含func且类型为func(v interface{}) string的映射
	Params       map[string]interface{} //具体验证器参数组，和Name参数一起来创建具体验证器，仅此作用
	theValidator ValidatorInterface     //具体应用的验证器
}

func (v *Validator) Init() {
	//初始化
	if v.SkipEmpty == nil {
		v.SkipEmpty = types.BoolPtr(true)
	}
	if v.SkipError == nil {
		v.SkipError = types.BoolPtr(true)
	}
	if v.IsEmpty == nil {
		v.IsEmpty = v.isEmpty
	}
	//When可以为nil
	// if v.When == nil {
	// 	v.When = v.when
	// }
}

//创建验证器
func (v *Validator) CreateValidator() ValidatorInterface {
	if valid, ok := BuiltInValidators[v.Name]; ok {
		if v.theValidator == nil {
			//创建具体验证器
			v.theValidator = valid(v.Params, v)
			//初始化验证器
			v.theValidator.Init()

			fmt.Printf("验证器%s： %+v 已经创建！\n", v.Name, v)
		}
		return v.theValidator
	}
	panic("不存在\"" + v.Name + "\"内置验证器！")
}

//ValidateAttributes 验证一组属性
func (v *Validator) ValidateAttributes(m ModelInterface, fs []string) {
	//只验证需要验证的属性
	attrs := v.getValidationAttributes(fs)
	fmt.Printf("%+v中待验证属性为%+v \n", fs, attrs)

	for _, attr := range attrs {
		fmt.Println("正在验证" + attr)

		//是否跳过错误或者空
		skip := *v.SkipError && m.HasErrors(attr) || *v.SkipEmpty && v.IsEmpty(m.GetAttributeValue(attr))
		if !skip {
			//是否需要验证
			if v.When == nil || v.When(m, attr) {
				//使用具体验证器来验证属性
				fmt.Println("验证中：" + attr)
				v.theValidator.ValidateAttribute(m, attr)
			} else {
				fmt.Println("当前无需验证：" + attr)
			}
		} else {
			fmt.Println("跳过验证：" + attr)
		}
	}
}

//getValidationAttributes 返回需要进行验证的属性数组
//
//fs中可能存在不需要进行验证的属性，fs中所有属性会自动转成小写驼峰
func (v *Validator) getValidationAttributes(fs []string) []string {
	if fs == nil {
		return v.getAttributeNames()
	}
	valids := []string{}
	attrs := v.getAttributeNames()
	for _, f := range fs {
		f = strs.FormatName(f, 5)
		for _, attr := range attrs {
			if f == attr {
				valids = append(valids, f)
				break
			}
		}
	}
	fmt.Printf("需要验证属性为：%+v \n", valids)
	return valids
}

//getAttributeNames 从Attributes属性中分析出真实的属性名称数组
func (v *Validator) getAttributeNames() []string {
	valids := []string{}
	for _, attr := range v.Attributes {
		valids = append(valids, strs.FormatName(strings.TrimLeft(attr, "!"), 5))
	}

	fmt.Printf("从Attributes属性中分析出所有属性名称数组：%+v \n", valids)

	return valids
}

//isAvtive 判断场景是否是活动的
func (v *Validator) isAvtive(scenario string) bool {
	var isExcept bool
	for _, e := range v.Except {
		if e == scenario {
			isExcept = true
			break
		}
	}
	if isExcept {
		return false
	}

	var isOn bool
	if len(v.On) > 0 {
		for _, o := range v.On {
			if o == scenario {
				isOn = true
				break
			}
		}
	} else {
		isOn = true
	}

	return isOn
}

//AddError 增加错误消息
func (v *Validator) AddError(m ModelInterface, attr, msg string, params map[string]string) {
	if params == nil {
		params = make(map[string]string)
	}

	params["attribute"] = m.GetAttributeLabel(attr)
	if _, ok := params["value"]; !ok {
		params["value"] = fmt.Sprintf("%s", m.GetAttributeValue(attr))
	}

	fmt.Println("错误为：", params)
	m.AddError(attr, v.formatMessage(msg, params))
}

//formatMessage 格式化错误消息
//
//将消息中的{xx}替换成yy，其中xx为params的键，yy为相应的值
func (v *Validator) formatMessage(msg string, params map[string]string) string {
	//TODO i18n忽略

	if len(params) == 0 {
		return msg
	}

	for n, v := range params {
		msg = strings.ReplaceAll(msg, "{"+n+"}", v)
	}

	return msg
}

//ValidateAttribute 验证属性，发送错误则在其模型中添加错误
func (v *Validator) ValidateAttribute(m ModelInterface, f string) {
	f = strs.FormatName(f, 5)
	fv := m.GetAttributeValue(f)

	//使用具体验证器来验证
	msg, params := v.theValidator.ValidateValue(fv)
	if msg != "" {
		v.AddError(m, f, msg, params)
	}
	fmt.Printf("验证结果为%s=%+v \n", msg, params)
}

func (v *Validator) ValidateValue(value interface{}) (msg string, formatParams map[string]string) {
	panic("ValidateValue函数必须要实现")
}

//验证普通值，非模型值
func (v *Validator) Validate(value interface{}) error {
	//使用具体验证器来验证
	msg, params := v.theValidator.ValidateValue(value)
	if msg == "" {
		return nil
	}

	//组织错误
	if params == nil {
		params = make(map[string]string)
	}
	params["attribute"] = "the input value"
	params["value"] = fmt.Sprintf("%s", value)

	return errors.New(v.formatMessage(msg, params))
}

//判断值是不是空的。
//
//注意bool值是假时，返回true
//
//一切数值==0，返回true，字符串空时返回true，其他为nil时返回true
func (v *Validator) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	//反射值
	sv := reflect.ValueOf(value)
	//是不是空值

	if !sv.IsValid() || sv.IsZero() {
		return true
	}

	return false
}
