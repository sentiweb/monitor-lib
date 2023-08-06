package notifier

import (
	"context"
	"fmt"
	"log"

	"github.com/sentiweb/monitor-lib/notify/types"
)

// MemoryNotifier stores notification in memory
// Mainly for testing purposes.
type MemoryNotifier struct {
	channel chan types.Notification
	count   int
	notifs  map[string][]types.Notification
}

func NewMemoryNotifier() *MemoryNotifier {
	n := make(map[string][]types.Notification)
	return &MemoryNotifier{notifs: n, channel: make(chan types.Notification, 10), count: 0}
}

func (c *MemoryNotifier) Accepts(n types.Notification) bool {
	return true
}

func (c *MemoryNotifier) Send(ctx context.Context, n types.Notification) error {
	c.channel <- n
	return nil
}

func (c *MemoryNotifier) handler(ctx context.Context) {
	for {
		select {
		case n := <-c.channel:
			id := n.ServiceName()
			m, ok := c.notifs[id]
			if !ok {
				m = make([]types.Notification, 0)
			}
			m = append(m, n)
			c.notifs[id] = m
			c.count++

		case <-ctx.Done():
			log.Println("Stopping Memory notifier service")
		}
	}
}

func (c *MemoryNotifier) Start(ctx context.Context) error {
	go c.handler(ctx)
	return nil
}

func (c *MemoryNotifier) String() string {
	return fmt.Sprintf("MemoryNotifier<%d>", c.count)
}

func (c *MemoryNotifier) Notifications(id string) []types.Notification {
	return c.notifs[id]
}

func (c *MemoryNotifier) MarshalYAML() (interface{}, error) {
	return "memory", nil
}

func (c *MemoryNotifier) Size() int {
	return c.count
}

func (c *MemoryNotifier) All() map[string][]types.Notification {
	return c.notifs
}
