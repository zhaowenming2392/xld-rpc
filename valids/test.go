package valids

import (
	"fmt"

	"helpers.zhaowenming.cn/types"
)

type Love struct {
	Name    string
	Age_man int
	Sex     bool
	Mode
}

func (l Love) AttributeLabels() map[string]string {
	return map[string]string{
		"name":   "姓名",
		"ageMan": "年龄",
		"sex":    "性别",
	}
}

//自定义验证函数，不能为特殊字符串
func T(v interface{}) string {
	v = v.(string)
	if v == "我啊" {
		return "不能乱写哦"
	}

	return ""
}

//TODO 写第二个版本，将[]string转成[]*Validator，以便后期通用
func (l Love) Rules() []*Validator {
	return []*Validator{
		{
			Attributes: []string{"name"},
			Name:       "func",
			SkipEmpty:  types.BoolPtr(false),
			Params: map[string]interface{}{
				"func": T,
			},
		},
		{
			Attributes: []string{"name"},
			Name:       "required",
			Params: map[string]interface{}{
				"requiredValue": 110,
			},
		},
		{
			Attributes: []string{"name"},
			Name:       "str",
			Params: map[string]interface{}{
				"min": 2,
				"max": 4,
			},
		},
		{
			Attributes: []string{"age_man"},
			Name:       "default",
			Params: map[string]interface{}{
				"value": 18,
			},
		},
		{
			Attributes: []string{"age-Man"},
			Name:       "num",
			On:         []string{"create"},
			Params: map[string]interface{}{
				"min": 10,
				"max": 99,
			},
		},
		{
			Attributes: []string{"sex"},
			Name:       "boolean",
		},
	}
}

func Test() {
	l := Love{
		Name:    "",
		Age_man: 6,
		Sex:     true,
	}
	l.SetMode(&l)
	l.SetScenario("create")
	ok := l.Validate(nil, true)
	if !ok {
		fmt.Println(l.GatErrors(""))
	}

	fmt.Printf("%+v \n", l)
}
