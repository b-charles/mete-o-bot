package main

import (
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/log"
	"github.com/b-charles/pigs/smartconfig"
	"gopkg.in/gomail.v2"
)

type UserMail struct {
	Name    string
	Address string
}

type MailConfig struct {
	From *UserMail
	To   []*UserMail
}

type SmtpConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func init() {
	smartconfig.Configure("mail", &MailConfig{})
	smartconfig.Configure("smtp", &SmtpConfig{})
}

type MailSender struct {
	Logger log.Logger  `inject:""`
	Mail   *MailConfig `inject:""`
	Smtp   *SmtpConfig `inject:""`
}

func (self *MailSender) send(msg *Message) {

	self.Logger.InfoLog("mail", "Dial SMTP server...")
	dialer := gomail.NewDialer(self.Smtp.Host, self.Smtp.Port, self.Smtp.User, self.Smtp.Password)
	sender, err := dialer.Dial()
	if err != nil {
		self.Logger.ErrorLog("dial.error", err.Error())
		return
	}

	message := gomail.NewMessage()
	for _, recipient := range self.Mail.To {

		self.Logger.Info().Set("sendingTo.name", recipient.Name).Log()

		message.SetAddressHeader("From", self.Mail.From.Address, self.Mail.From.Name)
		message.SetAddressHeader("To", recipient.Address, recipient.Name)
		message.SetHeader("Subject", msg.Title)
		message.SetBody("text/html", msg.Body)

		if err := gomail.Send(sender, message); err != nil {
			self.Logger.Error().SetAll(map[string]any{
				"sendingTo.name":    recipient.Name,
				"sendingTo.address": recipient.Address,
				"sendingTo.error":   err.Error(),
			}).Log()
		}

		message.Reset()

	}

}

func init() {
	ioc.Put(&MailSender{})
}
