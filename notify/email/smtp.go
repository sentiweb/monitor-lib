package email

import (
	"context"
	"fmt"

	"gopkg.in/mail.v2"
)

type SmtpSender struct {
	dsn    string // Use to text representation,
	dialer *mail.Dialer
}

func NewSmtpSender(host string, port int, user string, password string) *SmtpSender {
	var dsn string
	if user != "" {
		dsn = fmt.Sprintf("%s:%d", host, port)
	} else {
		dsn = fmt.Sprintf("%s:*****@%s:%d", user, host, port)
	}
	d := mail.NewDialer(host, port, user, password)
	return &SmtpSender{dialer: d, dsn: "smtp://" + dsn}
}

func (s *SmtpSender) Send(ctx context.Context, msg *mail.Message) error {
	return s.dialer.DialAndSend(msg)
}

func (s *SmtpSender) Start() error {
	// Nothing to do
	return nil
}

func (s *SmtpSender) String() string {
	// Nothing to do
	return s.dsn
}

func (p *SmtpSender) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}
