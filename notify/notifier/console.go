package notifier

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/notify/formatter"
)

// ConsoleNotifier shows Notification on the console. It's mainly for testing.
type ConsoleNotifier struct {
	delay time.Duration
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
