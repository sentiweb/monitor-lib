package notifier

import (
	"context"

	"github.com/sentiweb/monitor-lib/notify/types"
)

// MemoryNotifier stores notification in memory
// Mainly for testing purposes.
type MemoryNotifier struct {
	notifs map[string][]types.Notification
}

func NewMemoryNotifier() *MemoryNotifier {
	n := make(map[string][]types.Notification)
	return &MemoryNotifier{notifs: n}
}

func (c *MemoryNotifier) Accepts(n types.Notification) bool {
	return true
}

func (c *MemoryNotifier) Send(ctx context.Context, n types.Notification) error {
	id := n.ServiceName()
	m, ok := c.notifs[id]
	if !ok {
		m = make([]types.Notification, 0)
	}
	m = append(m, n)
	c.notifs[id] = m
	return nil
}

func (c *MemoryNotifier) Start(context.Context) error {
	return nil
}
func (c *MemoryNotifier) String() string {
	return "MemoryNotifier<>"
}

func (c *MemoryNotifier) Notifications(id string) []types.Notification {
	return c.notifs[id]
}

func (c *MemoryNotifier) MarshalYAML() (interface{}, error) {
	return "memory", nil
}
