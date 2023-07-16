package notify

import (
	"context"
	"testing"
	"time"

	"github.com/sentiweb/monitor-lib/notify/notifier"
	"github.com/sentiweb/monitor-lib/notify/types"

	"github.com/sentiweb/monitor-lib/notify/tests"
	utils_tests "github.com/sentiweb/monitor-lib/tests"
	"github.com/sentiweb/monitor-lib/utils"
)

func TestNotificationHanlder(t *testing.T) {

	defer utils_tests.CaptureLog(t).Release()

	handler := NewNotificationHandler(time.Second)

	memNotifier := notifier.NewMemoryNotifier()
	handler.AddNotifier(memNotifier)

	console := notifier.NewConsoleNotifier(0)
	handler.AddNotifier(console)

	ctx := context.WithValue(context.TODO(), utils.ContextDebug, true)

	notifChan := make(chan types.Notification)

	err := handler.Start(ctx, notifChan)
	if err != nil {
		t.Error(err)
	}

	serviceName := "test12341"

	notif := tests.NewMockNotification(serviceName, "up", serviceName, time.Now())

	notifChan <- notif

	aw := utils_tests.NewAwait(time.Second, 100*time.Millisecond)

	check := aw.Wait(func() bool {
		n := memNotifier.Notifications(serviceName)
		return len(n) == 1
	})

	t.Logf("Result %t, waited %s", check, aw.TimeWaited())
	if !check {
		t.Error("Notification not sent")
	}
}
