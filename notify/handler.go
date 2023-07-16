package notify

import (
	"context"
	"log"
	"time"

	"github.com/sentiweb/monitor-lib/notify/types"
)

// NotificationHandler listen for notification from a channel
// and route them to a list of Notifiers.
// Routing is handled by tags. Each notifier can accept a list of tags
// Notifier can be added with AddNotifier() method
// Handle() method starts to listen from the provided channel
type NotificationHandler struct {
	timeout  time.Duration
	debug    bool
	channels []types.Notifier
}

// NewNotificationHandler creates a NotificationHandler
func NewNotificationHandler(timeout time.Duration) *NotificationHandler {
	return &NotificationHandler{timeout: timeout}
}

func (h *NotificationHandler) setDebug(debug bool) {
	h.debug = debug
}

func (h *NotificationHandler) AddNotifier(channel types.Notifier) {
	h.channels = append(h.channels, channel)
}

// Start the handler
// Start each notifier and run the listening loop in a go routine
func (h *NotificationHandler) Start(ctx context.Context, input <-chan types.Notification) error {
	for _, channel := range h.channels {
		err := channel.Start(ctx)
		if err != nil {
			return err
		}
	}
	go h.handle(ctx, input)
	return nil
}

// Handle starts the listening for notification from the input channel
func (h *NotificationHandler) handle(ctx context.Context, input <-chan types.Notification) {
	log.Printf("Starting Notification Hander with %s timeout", h.timeout)
	for {
		select {
		case <-ctx.Done():
			log.Println("NotificationHandler", ctx.Err())
			return

		case notif := <-input:
			log.Printf("<- Notification<%s>", notif)
			for _, channel := range h.channels {
				if h.debug {
					log.Println("Notifier ", channel)
				}
				if channel.Accepts(notif) {
					go doNotify(ctx, channel, notif, h.timeout, h.debug)
				}
			}
		}
	}
}

func doNotify(ctx context.Context, nc types.Notifier, n types.Notification, timeout time.Duration, debug bool) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	notifyWithContext(ctx, nc, n, debug)
}

func notifyWithContext(ctx context.Context, notifier types.Notifier, n types.Notification, debug bool) {
	ch := make(chan error, 1)

	go func() {
		ch <- notifier.Send(ctx, n)
	}()

	select {

	case err := <-ch:
		if err != nil {
			log.Printf("-> Error with %s : %s", notifier, err)
		} else {
			if debug {
				log.Printf("-> Notification<%s> sent by %s", n, notifier)
			}
		}
		return

	case <-ctx.Done():
		ch <- nil
		log.Printf("-> Notification<%s> stopped, %s", notifier, ctx.Err())
		return
	}
}
