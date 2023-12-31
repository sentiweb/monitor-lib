package notify

import (
	"time"

	"github.com/sentiweb/monitor-lib/notify/types"
)

type NotifierBuilder struct {
	handler *NotificationHandler
}

func NewNotifierBuilder(globalTimeout time.Duration) *NotifierBuilder {
	notifHandler := NewNotificationHandler(globalTimeout)
	return &NotifierBuilder{handler: notifHandler}
}

func (builder *NotifierBuilder) Get() *NotificationHandler {
	return builder.handler
}

func (builder *NotifierBuilder) AddNotifier(notifier types.Notifier) {
	builder.handler.AddNotifier(notifier)
}
