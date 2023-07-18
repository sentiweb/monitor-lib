package notifier

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sentiweb/monitor-lib/notify/formatter"
	"github.com/sentiweb/monitor-lib/notify/types"
)

// ConsoleNotifier shows Notification on the console. It's mainly for testing.
type ConsoleNotifier struct {
	delay     time.Duration
	formatter types.Formatter
}

func NewConsoleNotifier(delay time.Duration) *ConsoleNotifier {
	return &ConsoleNotifier{delay: delay}
}

func (c *ConsoleNotifier) Accepts(n types.Notification) bool {
	return true
}

func (c *ConsoleNotifier) Start(context.Context) error {
	c.formatter = formatter.Get("console")
	return nil
}

func (c *ConsoleNotifier) Send(ctx context.Context, n types.Notification) error {
	body := c.formatter.Text(n)
	time.Sleep(c.delay)
	log.Println(body)
	return nil
}

func (c *ConsoleNotifier) String() string {
	return fmt.Sprintf("ConsoleNotifier<%s>", c.delay)
}

func (p *ConsoleNotifier) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}
