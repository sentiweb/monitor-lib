package notifier

import(
	"fmt"
	"context"
	"github.com/sentiweb/monitor-lib/notify/types"
)

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
	id := fmt.Sprintf("%s", n.Service())
	m, ok := c.notifs[id]
	if(!ok) {
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
