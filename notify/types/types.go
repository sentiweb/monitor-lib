package types

import (
	"fmt"
	"context"
	"time"
	"gopkg.in/mail.v2"
	"github.com/sentiweb/monitor-lib/utils"
)

// EmailSender An object able to send an email
type EmailSender interface {
	Start() error // Allow the sender to initialize
	Send(ctx context.Context, msg *mail.Message) error
}

const (
	NotificationStatusUp = "up"
	NotificationStatusDown = "down"
)

// Notification Interface represents the message to be notified
type Notification interface {
	fmt.Stringer
	Status() string // Notification status type (use NotificationStatusUp & NotificationStatusDown)
	Label() string // Notification Message to send
	FromTime() time.Time
	Tags() []string // List of tags to handle routing
	Service() fmt.Stringer // Provide a name identifier of the entity for which the notification is raised
}

// Notifier is able to send Notification
type Notifier interface {
	// Send a notification
	Send(context.Context, Notification) error
	// If Notifier accept this notification
	Accepts(Notification) bool
	// Allow the notifier to start with a global context (transmitted by AlertHandler.Handle())
	Start(context.Context) error
}

// WebhookNotifierService describes a webservice to send a message (json)
type WebhookNotifierService interface {
	Send(ctx context.Context, client utils.HTTPClient, notif Notification) (error)
	Start() error
}


// Formatter create text content for a given Notifier
// This can be used to format text differently depending on notifier
// We only provide a generic formatter, which send the same text everywhere
type Formatter interface {
	Title(n Notification) string
	Text(n Notification) string
}

type FormatterFactory interface {
	Get(notifierName string) Formatter
}