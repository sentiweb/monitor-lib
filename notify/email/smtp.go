notify

import(
	"context"
	"gopkg.in/mail.v2"
)

type SmtpSender struct {
	dialer *mail.Dialer	
}

func NewSmtpSender(host string, port int, user string, password string) *SmtpSender {
	d := mail.NewDialer(host, port, user, password)
	return &SmtpSender{dialer: d}
}

func (s* SmtpSender) Send(ctx context.Context, msg *mail.Message) error {
	return s.dialer.DialAndSend(msg)
}

func (s* SmtpSender) Start() error {
	// Nothing to do
}