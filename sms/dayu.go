package sms

import (
	"errors"
	"strings"

	//"fmt"
	//"github.com/aliyun/alibaba-cloud-sdk-go" //整个阿里云go语言SDK
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

//Sms 短信
/*
{
	"Message": "OK",
	"RequestId": "873043ac-bcda-44db-9052-2e204c6ed20f",
	"BizId": "607300000000000000^0",
	"Code": "OK"
}
*/
type Dayu struct {
	accessKey    string
	accessSecret string
}

//NewSms New Sms
func NewDayuSms(accessKey, accessSecret string) *Dayu {
	return &Dayu{
		accessKey,
		accessSecret,
	}
}

//SmsTpl SmsTpl
type DayuSmsTpl struct {
	SignName      string
	TemplateCode  string
	TemplateParam string
}

func NewDayuTpl(SignName, TemplateCode, TemplateParam string) *DayuSmsTpl {
	return &DayuSmsTpl{
		SignName,
		TemplateCode,
		TemplateParam,
	}
}

//Send Send sms，smsOrder 创建短信发送订单记录返回订单号，用于记录短信发送
func (s *Dayu) Send(phones []string, tpl *DayuSmsTpl, smsOrder func(phones []string, tpl *DayuSmsTpl) string) error {
	if len(phones) > 1000 {
		return errors.New("手机号码不能超过1000个")
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", s.accessKey, s.accessSecret)
	if err != nil {
		//fmt.Print(err.Error())
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = strings.Join(phones, ",") //"15900000000,13500000000" //支持对多个手机号码发送短信，手机号码之间以英文逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
	request.SignName = tpl.SignName                  //"短信签名名称"                                  //短信签名名称
	request.TemplateCode = tpl.TemplateCode          //"SMS_152550005"                       //短信模板ID
	request.TemplateParam = tpl.TemplateParam        //"{\"name\":\"短信模板变量对应的实际值，JSON格式\"}" //短信模板变量对应的实际值，JSON格式。
	request.SmsUpExtendCode = ""                     //无特殊需要此字段的用户请忽略此字段

	if smsOrder != nil {
		request.OutId = smsOrder(phones, tpl) //外部流水扩展字段。
	}

	response, err := client.SendSms(request)
	if err != nil {
		//fmt.Print(err.Error())
		return err
	}
	//fmt.Printf("response is %#v\n", response)

	if response.Code == "OK" {
		return nil
	}

	return errors.New(response.Message)
}

//SendBatchSms 批量发送短信
/*
正常返回格式
{
	"Message":"OK", 成功或失败消息
	"RequestId":"2184201F-BFB3-446B-B1F2-C746B7BF0657",
	"BizId":"197703245997295588^0",
	"Code":"OK" 不等于OK即为失败
}
*/
func (s *Dayu) SendBatchSms(phones []string, tpls []*DayuSmsTpl, smsOrder func(phones []string, tpl *DayuSmsTpl) string) error {
	if len(phones) > 1000 {
		return errors.New("手机号码不能超过1000个")
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", s.accessKey, s.accessSecret)
	if err != nil {
		//fmt.Print(err.Error())
		return err
	}

	request := dysmsapi.CreateSendBatchSmsRequest()
	request.Scheme = "https"

	if len(phones) == len(tpls) {
		request.PhoneNumberJson = "[" + strings.Join(phones, ",") + "]"
		request.SignNameJson = "["
		request.TemplateCode = "["
		//request.TemplateParamJson = "["
		request.SmsUpExtendCodeJson = "["

		for _, tpl := range tpls {
			request.SignNameJson += "\"" + tpl.SignName + "\""
			request.TemplateCode += "\"" + tpl.TemplateCode + "\""
			request.TemplateParamJson += "\"" + tpl.TemplateParam + "\""
			//request.SmsUpExtendCodeJson += "\""+tpl.SignName+"\""
		}

		request.SignNameJson += "]"
		request.TemplateCode += "]"
		//request.TemplateParamJson += "]"
		request.SmsUpExtendCodeJson += "]"
	} else {
		if len(tpls) != 1 {
			return errors.New("模板要么对应手机号码，要么统一使用一组")
		}

		request.PhoneNumberJson = "[" + strings.Join(phones, ",") + "]"
		request.SignNameJson = "["
		request.TemplateCode = "["
		//request.TemplateParamJson = "["
		request.SmsUpExtendCodeJson = "["

		tpl := tpls[0] //采用第一个

		for i := 0; i < len(phones); i++ {
			request.SignNameJson += "\"" + tpl.SignName + "\""
			request.TemplateCode += "\"" + tpl.TemplateCode + "\""
			request.TemplateParamJson += "\"" + tpl.TemplateParam + "\""
			//request.SmsUpExtendCodeJson += "\""+tpl.SignName+"\""
		}

		request.SignNameJson += "]"
		request.TemplateCode += "]"
		//request.TemplateParamJson += "]"
		request.SmsUpExtendCodeJson += "]"

		// request.PhoneNumberJson = "[\"15900000000\",\"13500000000\"]"
		// request.SignNameJson = "[\"短信签名名称，JSON数组格式\",\"短信签名的个数必须与手机号码的个数相同、内容一一对应\"]"
		// request.TemplateCode = "SMS_152550005"
		// request.TemplateParamJson = "[{\"name\":\"短信模板变量对应的实际值，JSON格式\"},{\"name\":\"模板变量值的个数必须与手机号码、签名的个数相同、内容一一对应\"}]"
		// request.SmsUpExtendCodeJson = "[\"90999\",\"上行短信扩展码，JSON数组格式\"]"
	}

	response, err := client.SendBatchSms(request)
	if err != nil {
		//fmt.Print(err.Error())
		return err
	}
	//fmt.Printf("response is %#v\n", response)

	if response.Code == "OK" {
		return nil
	}

	return errors.New(response.Message)
}
