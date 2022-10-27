package email

import (
	"crypto/tls"
	"log"

	gomail "gopkg.in/gomail.v2"
)

type EmailConfig struct {
	// Host represents the host of the SMTP server.
	Host string `json:"host"`
	// Port represents the port of the SMTP server.
	Port int `json:"port"`
	// Username is the username to use to authenticate to the SMTP server.
	Username string `json:"username"`
	// Password is the password to use to authenticate to the SMTP server.
	Password string `json:"password"`
	// SSL defines whether an SSL connection is used. It should be false in
	// most cases since the authentication mechanism should use the STARTTLS
	// extension instead.
	SSL bool `json:"ssl"`
	// TSLConfig represents the TLS configuration used for the TLS (when the
	// STARTTLS extension is used) or SSL connection.
	TLSConfig *tls.Config `json:"-"`

	//邮件来自
	From string
}

//Email Email
type Email struct {
	f string
	d *gomail.Dialer
}

//NewEmail New Email
func NewEmail(ec *EmailConfig) *Email {
	//TODO SSL连接和验证未测试
	if ec.SSL && ec.TLSConfig == nil {
		ec.TLSConfig = &tls.Config{ServerName: ec.Host}
	}

	//使用自定义连接
	return &Email{
		d: &gomail.Dialer{
			Host:      ec.Host,
			Port:      ec.Port,
			Username:  ec.Username,
			Password:  ec.Password,
			SSL:       ec.SSL,
			TLSConfig: ec.TLSConfig,
		},
		f: ec.From,
	}
}

//SendOneEmail 发送单个邮件
func (e *Email) SendOneEmail(subject, body string, tos []string, ccs []string, attachs []string) error {
	m := gomail.NewMessage()
	//m.SetHeader("From", e.f)//不含名称
	//m.SetHeader("From", m.FormatAddress(e.d.Username, e.f))//含名称
	m.SetAddressHeader("From", e.d.Username, e.f) //含名称，其实和上面是一样的

	m.SetHeader("To", tos...) //发送到一组目标，他们收到相同内容的同一个邮件
	//m.SetAddressHeader("To", r.Address, r.Name) //构建含有名称的邮箱地址，可以采用tos map[string地址]string名称

	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")//抄送
	m.SetHeader("Cc", ccs...)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	for _, at := range attachs {
		m.Attach(at)
	}

	// Send the email to Bob, Cora and Dan.
	if err := e.d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

//SendListEmail 发送一列邮件
func (e *Email) SendListEmail(subject, body string, tos []string, ccs []string, attachs []string) error {
	//x509：未知授权机构签署的证书
	//TODO 未测试使用&tls.Config{ServerName: host}能否实现ssl，从而不用绕过验证
	e.d.TLSConfig = &tls.Config{InsecureSkipVerify: true} //绕过服务器证书链和主机名的验证

	s, err := e.d.Dial()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	for _, r := range tos {
		m.SetHeader("From", e.f)
		m.SetHeader("To", r)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %v", r)
		}
		m.Reset()
	}

	return nil
}
