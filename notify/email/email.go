package email

import (
	"context"
	"errors"
	"log"

	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/utils"
	"gopkg.in/mail.v2"
)

// It's an implementation of types.EmailSender internal types.EmailSender implementations
type BaseEmailSender struct {
	from     string
	fromName string
	sender   types.EmailSender // Internal email sender
}

func New(from string, fromName string, options ...func(*BaseEmailSender)) *BaseEmailSender {
	svr := &BaseEmailSender{
		from:     from,
		fromName: fromName,
	}
	for _, o := range options {
		o(svr)
	}
	return svr
}

func WithSmtp(host string, port int, username string, password string) func(*BaseEmailSender) {
	if port == 0 {
		port = 25
	}
	return func(s *BaseEmailSender) {
		sender := NewSmtpSender(host, p, username, password)
		s.sender = sender
	}
}

func WithFake(path string) func(*BaseEmailSender) {
	return func(s *BaseEmailSender) {
		s.sender = NewFileSender(path)
	}
}

func (o *BaseEmailSender) Start() error {
	if o.sender == nil {
		return errors.New("Email sender must be defined")
	}
	if !utils.IsEmailValid(o.from) {
		log.Printf("Bad email format for '%s'", o.from)
		return utils.ErrBadEmail
	}
	return nil
}

func (o *BaseEmailSender) Send(ctx context.Context, msg *mail.Message) error {
	msg.SetHeader(utils.HeaderFrom, msg.FormatAddress(o.from, o.fromName))
	return o.sender.Send(ctx, msg)
}
