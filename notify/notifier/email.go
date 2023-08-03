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
	"github.com/sentiweb/monitor-lib/datastruct/sets"
	"gopkg.in/mail.v2"
)

// EmailNotifier sends notification to email thought an EmailSender instance
// A pooling mechanism can be used to pool several notifications arriving closely (using a pooling time window) together in the same email
type EmailNotifier struct {
	to        []string
	subject   string
	sender    types.EmailSender
	pool      []types.Notification
	sending   chan bool
	delay     time.Duration
	tags      *sets.Set[string]
	formatter types.Formatter
	mu        sync.RWMutex
}

// NewEmailNotifier creates an Email notifier instance
func NewEmailNotifier(sender types.EmailSender, to []string, subject string, poolingDelay time.Duration, tags []string) *EmailNotifier {
	return &EmailNotifier{
		to:      to,
		subject: subject,
		sender:  sender,
		delay:   poolingDelay,
		tags:    common.TagsToMap(tags),
	}
}

// Accepts checks if the notifier should send this types.Notification
func (c *EmailNotifier) Accepts(n types.Notification) bool {
	return common.CheckTags(c.tags, n.Tags())
}

func (c *EmailNotifier) String() string {
	return fmt.Sprintf("EmailNotifier<to:%s, pool:%s, tags:%s>", c.to, c.delay, c.tags)
}

func (c *EmailNotifier) MarshalYAML() (interface{}, error) {
	var m struct {
		Email struct {
			To      []string          `yaml:"to"`
			Subject string            `yaml:"subject"`
			Sender  types.EmailSender `yaml:"sender"`
			Delay   time.Duration     `yaml:"delay"`
			Tags    []string          `yaml:"tags,omitempty"`
		} `yaml:"email"`
	}
	m.Email.To = c.to
	m.Email.Subject = c.subject
	m.Email.Delay = c.delay
	m.Email.Sender = c.sender
	m.Email.Tags = common.MapToTags(c.tags)
	return m, nil
}

// Start the notifier service
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

// Send Email Notification
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
