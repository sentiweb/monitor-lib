package notifier

import (
	"context"
	"log"
	"testing"
	"time"

	utils_tests "github.com/sentiweb/monitor-lib/tests"

	"github.com/sentiweb/monitor-lib/notify/tests"
)

func TestMemoryNotifier(t *testing.T) {
	serviceName := "1"
	defer utils_tests.CaptureLog(t).Release()

	ctx := context.TODO()

	mn := NewMemoryNotifier()

	mn.Start(ctx)

	notif := tests.NewMockNotification(serviceName, "up", serviceName, time.Now())

	mn.Send(ctx, notif)

	aw := utils_tests.NewAwait(time.Second, 100*time.Millisecond)

	check := aw.Wait(func() bool {
		return mn.Size() == 1
	})

	if !check {
		t.Error("Notification not sent")
	}

	nn := mn.Notifications(serviceName)

	log.Println(nn)

	if len(nn) != 1 {
		t.Error("Expecting Memory notifier to have 1 notification")
	}

}
