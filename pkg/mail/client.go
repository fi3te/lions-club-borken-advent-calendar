package mail

import (
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/config"
	"github.com/wneessen/go-mail"
)

func SendMail(
	cfg config.EmailConfig,
	subject string,
	htmlContent string) error {

	message := mail.NewMsg()
	if err := message.FromFormat(cfg.Sender.Name, cfg.Sender.Address); err != nil {
		return err
	}
	if err := message.To(cfg.Recipients...); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, htmlContent)

	client, err := mail.NewClient(cfg.Smtp.Host,
		mail.WithSSL(),
		mail.WithSSLPort(false),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Smtp.Username),
		mail.WithPassword(cfg.Smtp.Password),
	)
	if err != nil {
		return err
	}
	return client.DialAndSend(message)
}
