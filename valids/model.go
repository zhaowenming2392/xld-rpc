/*
一：建立业务模型，嵌入本基础模型

二：安全验证器要求去编写验证器规则

三：可选择性丰富业务模型其他相关方法

四：通过模型去验证业务模型属性
*/
package valids

import (
	"fmt"
	//"net/http"
	"reflect"
	"strings"

	"helpers.zhaowenming.cn/maps"
	"helpers.zhaowenming.cn/strs"
)

const DefaultScenario = "default"
const CreateScenario = "create"
const UpdateScenario = "update"
const DeleteScenario = "delete"

//验证模型接口
//
//此模型专门来进行规则设置、验证相关及验证
type ModelInterface interface {
	Rules() []*Validator //当前验证器及其规则组：['password', 'compare', 'compareAttribute' => 'password2', 'on' => 'register']

	Attributes() []string                      //当前模型的所有属性，默认返回全部，改写可以使用属性
	AttributeLabels() map[string]string        //当前模型的各属性名称：name=>"名称"
	GetAttributeLabel(n string) string         //当前模型的某个属性名称
	GetAttributeValue(f string) interface{}    //获取属性值
	SetAttributeValue(f string, v interface{}) //设置属性值

	BeforeValid() bool //验证前运行，如果返回false立刻终止验证
	AfterValid()       //验证后运行

	AddError(f, e string)    //给属性f附加错误信息e
	HasErrors(f string) bool //判断有误错误，f=""则判断全部否则只判断f属性
}

//默认场景
const SCENARIO_DEFAULT = "default"

type Mode struct {
	//当前错误
	errors map[string][]string //属性及其一组错误消息

	//当前验证器
	validators []*Validator

	//当前场景
	scenario string
	//所有场景及其属性
	scenarios map[string][]string

	//当前模型
	mode ModelInterface
	//当前模型所有属性
	attributes []string
}

//Rules 返回模型的验证规则
func (m *Mode) Rules() []*Validator {
	return m.getMode().Rules()
}

//SetMode 设置具体业务模型实例
func (m *Mode) SetMode(mode ModelInterface) {
	if mode == nil {
		panic("业务模型实例不能为空")
	}
	m.mode = mode
}

//getMode 获取具体业务模型实例
func (m *Mode) getMode() ModelInterface {
	if m.mode == nil {
		panic("请先运行 SetMode() 设置业务模型实例")
	}
	return m.mode
}

//Scenarios 根据验证器中的定义返回场景及其活动属性
//
//未对场景进行大小写进行处理，但对其对应的属性进行了小写驼峰
func (m *Mode) Scenarios() map[string][]string {
	if m.scenarios != nil {
		return m.scenarios
	}

	scenarios := make(map[string][]string)
	scenarios[SCENARIO_DEFAULT] = []string{}

	//所有出现的场景
	for _, rule := range m.getMode().Rules() {
		for _, on := range rule.On {
			scenarios[on] = []string{}
		}
		for _, except := range rule.Except {
			scenarios[except] = []string{}
		}
	}

	//场景名称数组
	names := []string{}
	for n := range scenarios {
		names = append(names, n)
	}

	for _, rule := range m.getMode().Rules() {
		if rule.On == nil && rule.Except == nil {
			for _, n := range names {
				scenarios[n] = append(scenarios[n], rule.Attributes...)
			}
		} else if rule.On == nil {
			for _, n := range names {
				//场景是否除外
				var isExcept bool
				for _, scenario := range rule.Except {
					if scenario == n {
						isExcept = true
						break
					}
				}
				//的确要使用
				if !isExcept {
					scenarios[n] = append(scenarios[n], rule.Attributes...)
				}
			}
		} else {
			for _, n := range rule.On {
				scenarios[n] = append(scenarios[n], rule.Attributes...)
			}
		}
	}

	//去除无效场景
	for scenario, attrs := range scenarios {
		if len(attrs) == 0 {
			delete(scenarios, scenario)
		} else {
			//去重
			nowAttrs := []string{}
			for _, attr := range attrs {
				attr = strs.FormatName(attr, 5)

				var isExist bool
				for _, now := range nowAttrs {
					if now == attr {
						//存在
						isExist = true
						break
					}
				}

				if !isExist {
					nowAttrs = append(nowAttrs, attr)
				}
			}
			scenarios[scenario] = nowAttrs
		}
	}

	m.scenarios = scenarios

	return m.scenarios
}

//FormName 返回默认表单的名称
func (m *Mode) FormName() string {
	rt := reflect.TypeOf(m.getMode()).Elem()
	//首字母小写驼峰
	return strs.FormatName(rt.Name(), 5)
}

//Attributes 返回模型的所有属性,大小写不变
func (m *Mode) Attributes() []string {
	if m.attributes != nil {
		return m.attributes
	}

	rt := reflect.TypeOf(m.getMode()).Elem()
	attrs := []string{}
	for i := 0; i < rt.NumField(); i++ {
		fieldType := rt.Field(i)
		attrs = append(attrs, fieldType.Name)
	}

	m.attributes = attrs
	return m.attributes
}

//getAttributes 返回属性和值，同时会除去不应返回的属性，首字母小写
func (m *Mode) GetAttributes(names []string, except []string) map[string]interface{} {
	values := map[string]interface{}{}
	if names == nil {
		names = m.getMode().Attributes()
	}
	for _, name := range names {
		name = strs.FormatName(name, 5)
		values[name] = m.getMode().GetAttributeValue(name)
	}
	for _, name := range except {
		name = strs.FormatName(name, 5)
		delete(values, name)
	}

	return values
}

//SetAttributes 批量设置模型的属性值
//
//values属性及其值，safeOnly只对安全属性赋值
func (m *Mode) SetAttributes(values map[string]interface{}, safeOnly bool) {
	if len(values) > 0 {
		var names []string
		if safeOnly {
			names = m.SafeAttributes()
		} else {
			names = m.getMode().Attributes()
		}

		for name := range values {
			var isExist bool
			for _, n := range names {
				if n == name {
					isExist = true
					break
				}
			}
			if !isExist {
				delete(values, name)
				//提示不安全属性的赋值
				fmt.Printf("Failed to set unsafe attribute '%s' \n", name)
			}
		}

		if len(values) > 0 {
			//批量设置
			//TODO 是否可行？
			maps.SetMapToStruct(values, m.getMode())
		}
	}
}

//SafeAttributes 返回所有安全属性
func (m *Mode) SafeAttributes() []string {
	scenarios := m.Scenarios()
	scenario := m.GetScenario()

	if _, ok := scenarios[scenario]; !ok {
		return nil
	}
	attrs := []string{}
	for _, attr := range scenarios[scenario] {
		if attr[:1] != "!" {
			var isExist bool
			for _, attr2 := range scenarios[scenario] {
				if ("!" + attr) == attr2 {
					isExist = true
					break
				}
			}
			if !isExist {
				attrs = append(attrs, attr)
			}
		}
	}

	return attrs
}

//isAttributeRequired 是不是必须的属性，使用了RequiredValidator验证器
func (m *Mode) isAttributeRequired(attr string) bool {
	for _, validator := range m.getActiveValidators(attr) {
		_, ok := validator.theValidator.(*RequiredValidator)
		if ok && validator.When == nil {
			return true
		}
	}

	return false
}

//isAttributeSafe 是不是安全属性，安全属性可赋值
func (m *Mode) isAttributeSafe(attr string) bool {
	for _, name := range m.SafeAttributes() {
		if name == attr {
			return true
		}
	}

	return false
}

//isAttributeActive 是不是活动属性，活动属性需要验证
func (m *Mode) isAttributeActive(attr string) bool {
	for _, name := range m.activeAttributes() {
		if name == attr {
			return true
		}
	}

	return false
}

//GetAttributeValue 获取属性值
func (m *Mode) GetAttributeValue(f string) interface{} {
	if f == "" {
		return nil
	}

	rt := reflect.TypeOf(m.getMode()).Elem()
	rv := reflect.ValueOf(m.getMode()).Elem()
	for i := 0; i < rt.NumField(); i++ {
		if strings.EqualFold(strs.FormatName(rt.Field(i).Name, 5), strs.FormatName(f, 5)) {
			fmt.Printf("比对成功：%+v\n", rv.Field(i).Interface())
			return rv.Field(i).Interface()
		}
	}

	return nil
}

//SetAttributeValue 设置属性值
func (m *Mode) SetAttributeValue(f string, v interface{}) {
	if f != "" && v != nil {
		rt := reflect.TypeOf(m.getMode()).Elem()
		rv := reflect.ValueOf(m.getMode()).Elem()
		for i := 0; i < rt.NumField(); i++ {
			if strings.EqualFold(strs.FormatName(rt.Field(i).Name, 5), strs.FormatName(f, 5)) {
				if rt.Field(i).Type.Kind() == reflect.Ptr {
					/*
						fmt.Printf("目标属性是reflect.Ptr，当前值是 %+v\n", v)
						if !reflect.ValueOf(&v).Elem().CanAddr() {
							fmt.Printf("%p\n", &v)
						}
						if !rv.Field(i).CanAddr() {
							fmt.Printf("1 %+v\n", rv.Field(i))
						}
						if !rv.Field(i).Elem().CanAddr() {
							fmt.Println("rv.Field(i).Elem() not CanAddr")
						}
					*/
					//rv.Field(i).Elem().Set(reflect.ValueOf(&v).Elem())
					//TODO 尚未完成：针对指数类型的SET，发现最终类型不匹配，如：*int int
					//rv.Field(i).Set(reflect.ValueOf(&v).Elem())
					//rv.Field(i).Set(reflect.ValueOf(&v))
					//rv.Field(i).Set(reflect.ValueOf(v))
					//rv.Field(i).Elem().Set(reflect.ValueOf(&v))
					rv.Field(i).Elem().Set(reflect.ValueOf(v))
					//rv.Field(i).Set(reflect.ValueOf(&v).Elem())
				} else {
					rv.Field(i).Set(reflect.ValueOf(v))
				}
			}

		}
	}
}

//Validate 验证模型相关属性
//
//attrNames 需要验证的属性组，nil为验证全部，目前会自动对所有属性名称进行驼峰格式化
//
//clearErrors 是否清除原有错误，在某次验证流程下多次调用去验证某些个属性时，可以保留之前的错误信息
func (m *Mode) Validate(attrNames []string, clearErrors bool) bool {
	if clearErrors {
		fmt.Println("清除所有错误")
		m.clearErrors("")
	}

	if !m.getMode().BeforeValid() {
		fmt.Println("模型的BeforeValid返回了假，停止验证")
		return false
	}

	scenarios := m.Scenarios()
	fmt.Printf("当前所有场景为：%+v\n", scenarios)
	scenario := m.GetScenario()
	fmt.Printf("当前场景为：%+v\n", scenario)

	if _, ok := scenarios[scenario]; !ok {
		panic("Unknown scenario: " + scenario)
	}

	if len(attrNames) == 0 {
		fmt.Println("提供的待验证属性数值为空，现在验证所有属性")
		attrNames = m.activeAttributes()
	}

	fmt.Printf("所有待验证属性为：%+v\n", attrNames)

	for _, valid := range m.getActiveValidators("") {
		fmt.Printf("当前验证规则为：%+v\n", valid)
		fmt.Printf("当前验证器为：%+v\n", valid.theValidator)
		//使用具体验证器来验证
		valid.theValidator.ValidateAttributes(m.getMode(), attrNames)
	}

	m.getMode().AfterValid()

	fmt.Println("验证结束")

	return !m.HasErrors("")
}

//activeAttributes 返回当前场景需要验证的活动属性
func (m *Mode) activeAttributes() []string {
	scenarios := m.Scenarios()
	scenario := m.GetScenario()

	if _, ok := scenarios[scenario]; !ok {
		return nil
	}

	attrs := []string{}
	for _, f := range scenarios[scenario] {
		//以 ! 开头的需要验证但是不可赋值
		if f[:1] == "!" {
			f = f[1:]
		}
		attrs = append(attrs, strs.FormatName(f, 5))
	}

	fmt.Printf("所有活动属性为%+v \n", attrs)
	return attrs
}

//getActiveValidators 返回适用于当前[[场景]]的验证程序
//f 返回适用此属性的验证程序，空返回全部
func (m *Mode) getActiveValidators(f string) []*Validator {
	fmt.Println("(m *Mode) getActiveValidators")
	activeAttrs := m.activeAttributes()
	if f != "" {
		var isExist bool
		for _, attr := range activeAttrs {
			if attr == f {
				isExist = true
			}
		}
		if !isExist {
			fmt.Println(f + "不在活动属性中无需验证！")
			return nil
		}
	}

	scenario := m.GetScenario()
	validators := []*Validator{}
	for _, validator := range m.getValidators() {
		var attrValid bool
		if f == "" {
			validatorAttributes := validator.getValidationAttributes(activeAttrs)
			attrValid = len(validatorAttributes) > 0
		} else {
			validatorAttributes := validator.getValidationAttributes(activeAttrs)
			for _, attr := range validatorAttributes {
				if attr == f {
					attrValid = true
					break
				}
			}
		}
		if attrValid && validator.isAvtive(scenario) {
			fmt.Println("当前活动验证器：", validator.Name, "----", validator)
			validators = append(validators, validator)
		}
	}

	return validators
}

//返回所有验证器
func (m *Mode) getValidators() []*Validator {
	fmt.Println("(m *Mode) getValidators()")
	if m.validators == nil {
		m.validators = m.createValidators()
	}

	return m.validators
}

func (m *Mode) createValidators() (vs []*Validator) {
	fmt.Println("(m *Mode) createValidators()")
	if m.getMode().Rules() == nil {
		return nil
	}

	for _, rule := range m.getMode().Rules() {
		//创建具体验证器
		rule.CreateValidator()
		vs = append(vs, rule)
	}

	return vs
}

//BeforeValid 验证之前，返回false停止验证
func (m *Mode) BeforeValid() bool {
	return true
}

//AfterValid 验证后
func (m *Mode) AfterValid() {
}

//GetScenario 获取当前场景
func (m *Mode) GetScenario() string {
	if m.scenario == "" {
		m.scenario = DefaultScenario
	}
	return m.scenario
}

//SetScenario 设置当前场景
func (m *Mode) SetScenario(v string) {
	m.scenario = v
}

//AttributeLabels 默认返回英文标签
//
//注意键是小写驼峰，值是大写驼峰
func (m *Mode) AttributeLabels() map[string]string {
	attrs := m.getMode().Attributes()
	labels := make(map[string]string)
	for _, attr := range attrs {
		labels[strs.FormatName(attr, 5)] = strs.FormatName(attr, 6)
	}

	return labels
}

//GetAttributeLabel 获取某个属性的标签，如果不存在原样返回
func (m *Mode) GetAttributeLabel(n string) string {
	labels := m.getMode().AttributeLabels()

	n = strs.FormatName(n, 5)
	if label, ok := labels[n]; ok {
		return label
	}

	return strs.FormatName(n, 6)
}

//是否有错误，提供f则为此字段是否有错误，空则为整个模型
func (m *Mode) HasErrors(f string) bool {
	if len(m.errors) == 0 {
		return false
	}
	if f == "" {
		return true
	} else {
		_, ok := m.errors[f]
		return ok
	}
}

//获取错误，提供f则为此字段错误，空则为整个模型
func (m *Mode) GatErrors(f string) map[string][]string {
	if m.errors == nil {
		return nil
	}

	if f == "" {
		return m.errors
	} else {
		es, ok := m.errors[f]
		if ok {
			return map[string][]string{
				f: es,
			}
		}
	}

	return nil
}

//获取第一个错误，提供f则仅为此字段错误，空则为整个模型
func (m *Mode) GatFirstError(f string) string {
	if m.errors == nil {
		return ""
	}

	if f == "" {
		for _, es := range m.errors {
			return es[0]
		}
	} else {
		es, ok := m.errors[f]
		if ok {
			return es[0]
		}
	}

	return ""
}

//增加一个字段的错误
func (m *Mode) AddError(f, e string) {
	if m.errors == nil {
		m.errors = make(map[string][]string)
	}

	m.errors[f] = append(m.errors[f], e)
}

//addErrors 增加一组错误
func (m *Mode) addErrors(items map[string][]string) {
	if m.errors == nil {
		m.errors = make(map[string][]string)
	}

	for f, item := range items {
		m.errors[f] = append(m.errors[f], item...)
	}
}

//清除某个字段的错误
func (m *Mode) clearErrors(f string) {
	if m.errors == nil {
		return
	}

	if f == "" {
		m.errors = make(map[string][]string)
	} else {
		delete(m.errors, f)
	}
}

/*
//从http请求中获取数据并设置到结构中
func (m *Mode) Load(data map[string]interface{}, isForm bool, formName string) bool {
	if isForm {
		if formName == "" {
			formName = m.FormName()
		}
	} else {
		formName = ""
	}

	scope := formName
	if scope == "" && data != nil {
		m.SetAttributes(data, true)
		return true
	} else if value, ok := data[scope]; ok {
		m.SetAttributes(value.(map[string]interface{}), true)
		return true
	}
	return false
}

//分析http请求中的相应数据
func (m *Mode) getReqData(req http.Request) map[string]interface{} {
	err := req.ParseMultipartForm(20 << 10)
	if err != nil {
		return nil
	}

	return nil
}
*/
