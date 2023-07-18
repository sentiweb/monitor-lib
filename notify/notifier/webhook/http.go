package webhook

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sentiweb/monitor-lib/notify/common"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/utils"
)

// HTTPNotifier implements a notifier service sending the notification using an http request
// It's used to send notification to external webhooks
// HTTPNotifier embeds the common logic (handling loop, routing)
type HTTPNotifier struct {
	service  types.WebhookNotifierService
	poolSize uint
	timeout  time.Duration
	tags     map[string]struct{}
	out      chan types.Notification
	client   utils.HTTPClient
}

type HTTPNotifierOption func(*HTTPNotifier)

// NewHTTPNotifier create a new HTTP Notifier service connected to specific Webhook Service
func NewHTTPNotifier(service types.WebhookNotifierService, options ...HTTPNotifierOption) *HTTPNotifier {
	h := &HTTPNotifier{
		service:  service,
		timeout:  time.Second * 10,
		poolSize: 10,
	}
	for _, o := range options {
		o(h)
	}
	return h
}

// WithTags defines tags for which the service accept the notification
func WithTags(tags []string) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.tags = common.TagsToMap(tags)
	}
}

// WithTimeout defines the client
func WithTimeout(timeout time.Duration) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.timeout = timeout
	}
}

// WithPoolSize defines the size of the internal notification channel
// It defines the number of notifications waiting to be sent.
// If set to 1, the service will block the incoming notification (from another goroutines) until each notification is sent.
func WithPoolSize(size uint) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.poolSize = size
	}
}

// Create the http Client,
var HttpFactory utils.HTTPClientFactory = &utils.DefaultHTTPFactory{}

// Accepts checks if types.Notifier can send this types.Notification
func (c *HTTPNotifier) Accepts(n types.Notification) bool {
	return common.CheckTags(c.tags, n.Tags())
}

func (c *HTTPNotifier) String() string {
	return fmt.Sprintf("HTTPNotifier<%d,%d,%s,%s, %s>", c.timeout, c.poolSize, c.tags, c.service, c.client)
}

// Start HTTP Notifier service
func (c *HTTPNotifier) Start(ctx context.Context) error {
	err := c.service.Start()
	if err != nil {
		return err
	}
	c.client = HttpFactory.NewClient(c.timeout, nil)
	c.out = make(chan types.Notification, c.poolSize)
	go func(ctx context.Context) {
		err := c.handle(ctx)
		if err != nil {
			log.Println("Stopping HTTP types.Notifier handler", err)
		}
	}(ctx)
	return nil
}

// Send a notification to the notifiers
func (c *HTTPNotifier) Send(ctx context.Context, n types.Notification) error {
	c.out <- n
	return nil
}

// handle listen to notification and route it to the underlying http service to actually send it
func (c *HTTPNotifier) handle(ctx context.Context) error {
	debug := utils.DebugFromContext(ctx)
	fmt.Printf("Starting HTTPNotifier (debug %t)\n", debug)
	for {
		select {
		case n := <-c.out:
			if debug {
				log.Printf("Sending Notification of %s", n)
			}
			err := c.service.Send(ctx, c.client, n)
			if err != nil {
				log.Println("Error during sending of types.Notification", err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}

	}
}

// MarshalYAML
// Accept Mashshal to Yaml (for output representation)
func (c *HTTPNotifier) MarshalYAML() (interface{}, error) {
	var m struct {
		W struct {
			Service types.WebhookNotifierService
			Pool    uint          `yaml:"pool"`
			Timeout time.Duration `yaml:"timeout"`
			Tags    []string      `yaml:"tags,omitempty"`
		} `yaml:"webhook"`
	}
	m.W.Service = c.service
	m.W.Pool = c.poolSize
	m.W.Timeout = c.timeout
	m.W.Tags = common.MapToTags(c.tags)
	return m, nil
}
