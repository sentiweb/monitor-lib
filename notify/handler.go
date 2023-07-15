package notify

import(
	"time"
	"log"
	"context"
	"github.com/sentiweb/monitor-lib/notify/types"
	
)

type NotificationHandler struct {
	timeout time.Duration
	debug bool
	channels []types.Notifier
}

func NewNotificationHandler(timeout time.Duration) *NotificationHandler {
	return &NotificationHandler{timeout: timeout}
}

func (h *NotificationHandler) setDebug(debug bool) {
	h.debug = debug
}

func (h *NotificationHandler) AddNotifier(channel types.Notifier) {
	h.channels = append(h.channels, channel)
}

func (h *NotificationHandler) Handle(ctx context.Context, input <-chan types.Notification) error {
	log.Printf("Starting Notification Hander with %s timeout", h.timeout)
	
	for _, channel := range h.channels {
		channel.Start(ctx)
	}

	for {
		select {
			case <-ctx.Done():
				return ctx.Err()

			case notif := <-input:
				log.Printf("<- Notification<%s>", notif)
				for _, channel := range h.channels {
					if(h.debug) {
						log.Println("Notifier ", channel)
					}
					if(channel.Accepts(notif)) {
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
			if(err != nil) {
				log.Printf("-> Error with %s : %s", notifier, err )
			} else {
				if(debug) {
					log.Printf("-> Notification<%s> sent by %s", n, notifier )
				}
			}
			return

		case <-ctx.Done():
			ch <- nil
			log.Printf("-> Notification<%s> stopped, %s", notifier, ctx.Err() )
			return		
	}
}

