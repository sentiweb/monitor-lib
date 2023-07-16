package notifier

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sentiweb/monitor-lib/notify/common"
	"github.com/sentiweb/monitor-lib/notify/formatter"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/utils"
	"gopkg.in/mail.v2"
)

// Default types.Notification channel accepts all and send email
// Can pool types.Notifications by waiting some delay before sending all
type EmailNotifier struct {
	to        []string
	subject   string
	sender    types.EmailSender
	pool      []types.Notification
	sending   chan bool
	delay     time.Duration
	tags      map[string]struct{}
	formatter types.Formatter
	mu        sync.RWMutex
}

// NewEmailNotifier creates an Email types.Notifier instance
func NewEmailNotifier(sender types.EmailSender, to []string, subject string, delay time.Duration, tags map[string]struct{}) *EmailNotifier {
	return &EmailNotifier{
		to:      to,
		subject: subject,
		sender:  sender,
		delay:   delay,
		tags:    tags,
	}
}

// Accepts checks if the types.Notifier should send this types.Notification
func (c *EmailNotifier) Accepts(n types.Notification) bool {
	return common.CheckTags(c.tags, n.Tags())
}

func (c *EmailNotifier) String() string {
	return fmt.Sprintf("EmailNotifier<%s, %s>", c.to, c.delay)
}

// Start the types.Notifier service
func (c *EmailNotifier) Start(ctx context.Context) error {

	c.formatter = formatter.Get("email")

	err := c.sender.Start()
	if err != nil {
		return err
	}

	c.sending = make(chan bool)

	// Launch handler thread
	go func(ctx context.Context) {
		err := c.handle(ctx)
		if err != nil {
			log.Println("Stopping Email notifier handler", err)
		}
	}(ctx)

	return nil
}

// Send Email types.Notification
func (c *EmailNotifier) Send(ctx context.Context, n types.Notification) error {
	c.mu.Lock()
	c.pool = append(c.pool, n)
	c.mu.Unlock()
	c.sending <- true
	return nil
}

// Pooler handler
func (c *EmailNotifier) handle(ctx context.Context) error {
	log.Printf("%s start", c.String())
	timer := time.NewTimer(c.delay)
	debug := utils.DebugFromContext(ctx)
	waiting := false
	for {
		select {

		case <-timer.C:
			if waiting {
				if debug {
					log.Println(c, "sending pool")
				}
				c.mu.Lock()
				p := make([]types.Notification, len(c.pool))
				copy(p, c.pool)
				c.mu.Unlock()
				err := c.sendPool(ctx, p)
				if err != nil {
					log.Println("EmailNotifier sendPool said :", err)
				}
				waiting = false
			}

		case <-c.sending:
			if !waiting {
				waiting = true
				if debug {
					log.Println(c, "waiting ")
				}
				timer.Reset(c.delay)
			}

		case <-ctx.Done():
			return ctx.Err()
		}

	}
}

func (c *EmailNotifier) sendPool(ctx context.Context, pool []types.Notification) error {
	m := mail.NewMessage()
	m.SetHeader(utils.HeaderSubject, c.subject)
	for _, t := range c.to {
		m.SetHeader(utils.HeaderTo, t)
	}

	var b strings.Builder

	for _, n := range pool {
		b.WriteString("- ")
		b.WriteString(c.formatter.Text(n))
		b.WriteString("\n")
	}
	body := b.String()
	m.SetBody(utils.MimeTextPlain, body)
	err := c.sender.Send(ctx, m)
	if err == nil {
		c.pool = nil
	}
	return err
}
