package webhook

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/notify/notifier/tags"
	"github.com/sentiweb/monitor-lib/utils"
)

type HTTPNotifier struct {
	service types.WebhookNotifierService
	poolSize uint
	timeout time.Duration
	tags    map[string]interface{}
	out     chan types.Notification
	client  utils.HTTPClient
}

func arrayToMap(tags []string) map[string]interface{} {
	m := make(map[string]interface{}, len(tags))
	for _, s := range tags {
		m[s] = true
	}
	return m
}

func NewHTTPNotifier(service types.WebhookNotifierService, options ...func(*HTTPNotifier)) *HTTPNotifier {
	h := &HTTPNotifier{
		service: service,
		timeout: time.Second * 10,
		poolSize: 10,
	}
	for _, o := range options {
		o(h)
	}
	return h
}

func WithTags(tags []string) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.tags = arrayToMap(tags)
	}
}

func WithTimeout(timeout time.Duration) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.timeout = timeout 
	}
}

func WithPoolSize(size uint) func(*HTTPNotifier) {
	return func(h *HTTPNotifier) {
		h.poolSize = size 
	}
}

// Create the http Client,
var HttpFactory utils.HTTPClientFactory = &utils.DefaultHTTPFactory{}

// Accepts checks if types.Notifier can send this types.Notification
func (c *HTTPNotifier) Accepts(n types.Notification) bool {
	return tags.CheckTags(c.tags, n.Tags())
}

func (c *HTTPNotifier) String() string {
	return fmt.Sprintf("HTTPNotifier<%d,%d,%s,%s, %s>", c.timeout, c.poolSize, c.tags, c.service, c.client)
}

// Start HTTP types.Notifier service
func (c *HTTPNotifier) Start(ctx context.Context) error {
	err := c.service.Start()
	if(err != nil) {
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

func (c *HTTPNotifier) Send(ctx context.Context, n types.Notification) error {
	c.out <- n
	return nil
}

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