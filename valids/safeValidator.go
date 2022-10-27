package valids

//SafeValidator 安全验证器
//
//该验证器并不进行数据验证。而是把一个属性标记为 安全属性。
type SafeValidator struct {
	*Validator
}

//创建安全验证器
func NewSafeValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	//创建
	dv := SafeValidator{}
	//设置默认值，包括Validator中某些属性的默认值
	dv.Validator = v

	return &dv
}

func (sv *SafeValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	return "", nil
}

func (sv *SafeValidator) ValidateAttributes(m ModelInterface, fs []string) {
}

func (sv *SafeValidator) ValidateAttribute(m ModelInterface, f string) {
}
