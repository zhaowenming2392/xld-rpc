package email

import (
	"crypto/tls"
	"log"
	"time"

	gomail "gopkg.in/gomail.v2"
)

//ChanEmail ChanEmail
type ChanEmail struct {
	host      string
	port      int
	username  string
	password  string
	isSSL     bool
	tlsConfig *tls.Config

	from string

	send  chan *gomail.Message
	close chan bool
}

//NewChanEmail 创建邮件发送通过
func NewChanEmail(ec *EmailConfig) *ChanEmail {
	//TODO SSL连接和验证未测试
	if ec.SSL && ec.TLSConfig == nil {
		ec.TLSConfig = &tls.Config{ServerName: ec.Host}
	}

	//使用自定义连接
	ce := &ChanEmail{
		host:      ec.Host,
		port:      ec.Port,
		username:  ec.Username,
		password:  ec.Password,
		isSSL:     ec.SSL,
		tlsConfig: ec.TLSConfig,

		from: ec.From,

		send:  make(chan *gomail.Message, 5),
		close: make(chan bool),
	}

	return ce
}

//Send 使用ChanEmail发送邮件
func (ce *ChanEmail) Send(subject, body string, tos []string, ccs []string, attachs []string) {
	m := gomail.NewMessage()
	m.SetHeader("From", ce.from)
	m.SetHeader("To", tos...)
	//m.SetAddressHeader("To", r.Address, r.Name) //构建含有名称的邮箱地址
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Cc", ccs...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	for _, at := range attachs {
		m.Attach(at)
	}

	ce.send <- m
}

//Run 运行ChanEmail协程
func (ce *ChanEmail) Run() {
	go func() {
		for {
			select {
			case m := <-ce.send:
				//发送
				go func() {
					d := gomail.Dialer{
						Host:      ce.host,
						Port:      ce.port,
						Username:  ce.username,
						Password:  ce.password,
						SSL:       ce.isSSL,
						TLSConfig: ce.tlsConfig,
					}

					var s gomail.SendCloser
					var err error
					open := false

					// Close the connection to the SMTP server if no email was sent in
					// the last 30 seconds.
					select {
					case <-time.After(30 * time.Second):
						if open {
							if err := s.Close(); err != nil {
								log.Println(err)
							}
							open = false
						}

						return
					default:
						if !open {
							if s, err = d.Dial(); err != nil {
								log.Println(err)
								return
							}
							open = true
						}

						if err := gomail.Send(s, m); err != nil {
							log.Println(err)
						}

						return
					}
				}()
			case <-ce.close:
				//关闭
				close(ce.send)
				close(ce.close)
				return
			}
		}
	}()
}

//Close ChanEmail
func (ce *ChanEmail) Close() {
	ce.close <- true
}
