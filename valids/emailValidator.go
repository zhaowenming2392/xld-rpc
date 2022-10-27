package valids

//EmailValidator 邮箱地址验证器
//
//验证邮箱地址
type EmailValidator struct {
	*MatchValidator
}

//创建邮箱验证器
func NewEmailValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	//创建
	ev := EmailValidator{}
	//设置默认值，包括Validator中某些属性的默认值
	ev.Validator = v
	//ev.Pattern = "/^[a-zA-Z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$/"
	//TODO \.改成了\\. 格式对吗？这个正则这么复杂？
	ev.Pattern = "/^[a-zA-Z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\\.[a-zA-Z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$/"
	ev.Message = "{attribute} 不是邮箱地址！"

	return &ev
}
