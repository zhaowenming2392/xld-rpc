package valids

import (
	"fmt"
	"time"

	"helpers.zhaowenming.cn/maps"
)

//DateValidator 时间验证器
//
//验证该属性是否以正确的[[格式]]表示日期、时间或日期时间
//
//除了验证日期，它还可以将解析的时间戳导出为机器可读的格式可以使用[[timestampAttribute]]配置
type DateValidator struct {
	Kind int //时间类型，0=datetime，1=time，2=date

	TimeAttribute     string //某个属性
	AttributeFormat   string //某个属性时间格式，如果是空会设置时间戳给TimeAttribute
	AttributeTimeZone string //某个属性时间时区，如果是空和TimeZone相同

	format string

	Max interface{} //用户输入 string（满足格式的字符串）|int 时间戳
	Min interface{} //string|int

	max int64
	min int64

	TimeZone string //时区，注意这里按照统一时区了

	TooBig    string
	TooSmall  string
	MaxString string //在错误消息中显示的用户友好的下限值
	MinString string //在错误消息中显示的用户友好的上限值

	*Validator
}

func NewDateValidator(params map[string]interface{}, v *Validator) ValidatorInterface {
	dv := DateValidator{}
	dv.Validator = v

	err := maps.SetMapToStruct(params, &dv)
	if err != nil {
		panic("NewDateValidator err=" + err.Error())
	}

	return &dv
}

//Init 初始化，注意是在对具体验证器进行设置之后的初始化
func (dv *DateValidator) Init() {
	//父级初始化
	dv.Validator.Init()

	//初始化
	if dv.Message == "" {
		dv.Message = "{attribute} 的时间格式无效！"
	}

	//没有指定类型，直接使用默认
	switch dv.Kind {
	case 0: //默认
		dv.format = "2006-01-02 15:04:05"
	case 1:
		dv.format = "15:04:05"
	case 2:
		dv.format = "2006-01-02"
	default:
		panic("时间格式种类有误，只能小于3！")
	}

	if dv.TimeZone == "" {
		dv.TimeZone = "Asia/Shanghai"
	}
	if dv.AttributeTimeZone == "" {
		dv.AttributeTimeZone = dv.TimeZone
	}

	if dv.Min != nil && dv.TooSmall == "" {
		dv.TooSmall = "{attribute} 不能早于 {min}"
	}
	if dv.Max != nil && dv.TooBig == "" {
		dv.TooBig = "{attribute} 不能晚于 {max}"
	}
	if dv.MinString == "" {
		dv.MinString = fmt.Sprintf("%s", dv.Min)
	}
	if dv.MaxString == "" {
		dv.MaxString = fmt.Sprintf("%s", dv.Max)
	}

	//全部转成int格式
	if dv.Max != nil {
		dv.max = dv.parseDateValue(dv.Max, dv.TimeZone)
		if dv.max == 0 {
			panic("解析Max时间出错：不满足格式的字符串或非时间戳")
		}
	}

	if dv.Min != nil {
		dv.min = dv.parseDateValue(dv.Min, dv.TimeZone)
		if dv.min == 0 {
			panic("解析Min时间出错：不满足格式的字符串或非时间戳")
		}
	}
}

//parseDateValue 解析时间，时间必须满足格式或者为int/int64格式
func (dv *DateValidator) parseDateValue(v interface{}, zone string) int64 {
	var theV int64
	if t, ok := v.(string); ok {
		loc, err := time.LoadLocation(zone)
		if err != nil {
			panic("解析时区出错：" + err.Error())
		}
		tTime, err := time.ParseInLocation(dv.format, t, loc)
		if err != nil {
			panic("解析Min时间出错：" + err.Error())
		}
		theV = tTime.Unix()
	}

	if t, ok := v.(int); ok {
		theV = int64(t)
	}
	if t, ok := v.(int64); ok {
		theV = t
	}

	return theV
}

func (dv *DateValidator) ValidateValue(v interface{}) (msg string, formatParams map[string]string) {
	theV := dv.parseDateValue(v, dv.TimeZone)
	if theV == 0 {
		return dv.Message, nil
	}

	if dv.min != 0 && theV < dv.min {
		return dv.TooSmall, map[string]string{
			"min": dv.MinString,
		}
	}

	if dv.max != 0 && theV > dv.max {
		return dv.TooBig, map[string]string{
			"max": dv.MaxString,
		}
	}

	return "", nil
}

func (dv *DateValidator) ValidateAttribute(m ModelInterface, f string) {
	value := m.GetAttributeValue(f)

	if dv.IsEmpty(value) {
		//空，而且需要设置时间属性，则也都为空
		if dv.TimeAttribute != "" {
			m.SetAttributeValue(dv.TimeAttribute, value)
		}

		return
	}

	//解析成时间戳，然后先进行比较
	theV := dv.parseDateValue(value, dv.AttributeTimeZone)
	if theV == 0 {
		m.AddError(f, dv.formatMessage(dv.Message, nil))
	} else if dv.min != 0 && theV < dv.min {
		m.AddError(f, dv.formatMessage(dv.Message, map[string]string{
			"min": dv.MinString,
		}))
	} else if dv.max != 0 && theV > dv.max {
		m.AddError(f, dv.formatMessage(dv.Message, map[string]string{
			"max": dv.MaxString,
		}))
	} else if dv.TimeAttribute != "" {
		//满足条件，又需要赋值
		if dv.AttributeFormat == "" {
			//没有格式，赋值为时间戳
			// if int64(int(theV)) == theV {
			// 	//int
			// 	m.SetAttributeValue(f, int(theV))
			// } else {
			// 	m.SetAttributeValue(f, theV)
			// }

			m.SetAttributeValue(f, theV)
		} else {
			//按照格式赋值
			t := time.Unix(theV, 0)
			m.SetAttributeValue(f, t.Format(dv.AttributeFormat))
		}
	}
}
