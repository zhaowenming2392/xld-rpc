package valids

import "helpers.zhaowenming.cn/maps"

//UrlValidator 网址验证器
//
//验证网址
type UrlValidator struct {
	ValidSchemes  []string //有效的URI协议方案列表
	DefaultScheme string   //默认的URI方案。如果输入不包含scheme部分，则默认的scheme将被添加到它前面（从而更改输入）。默认值为null，这意味着URL必须包含方案部分。
	*MatchValidator
}

//创建网址验证器
func NewUrlValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	//创建
	uv := UrlValidator{}
	//设置默认值，包括Validator中某些属性的默认值
	uv.Validator = v
	uv.Pattern = `/^{schemes}:\/\/(([A-Z0-9][A-Z0-9_-]*)(\.[A-Z0-9][A-Z0-9_-]*)+)(?::\d{1,5})?(?:$|[?\/#])/i`
	uv.Message = "{attribute} 不是正确的网址！"

	uv.ValidSchemes = []string{
		"http", "https",
	}
	uv.DefaultScheme = "http"

	//批量设置参数
	err := maps.SetMapToStruct(params, &uv)
	if err != nil {
		panic("NewUrlValidator err=" + err.Error())
	}

	return &uv
}
