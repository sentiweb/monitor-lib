package notifier

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sentiweb/monitor-lib/notify/email"
	"github.com/sentiweb/monitor-lib/notify/tests"
	utils_tests "github.com/sentiweb/monitor-lib/tests"
	"github.com/sentiweb/monitor-lib/utils"
)

func TestEmailNotifier(t *testing.T) {

	defer utils_tests.CaptureLog(t).Release()

	sender := email.NewMemorySender()

	notifier := NewEmailNotifier(sender, []string{"me@example.com"}, "Notification", 0, nil)

	context := context.WithValue(context.TODO(), utils.ContextDebug, true)

	err := notifier.Start(context)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(notifier)

	serviceName := "test12341"

	notif := tests.NewMockNotification(serviceName, "up", serviceName, time.Now())

	notifier.Send(context, notif)

	aw := utils_tests.NewAwait(time.Second, 100*time.Millisecond)

	check := aw.Wait(func() bool {
		return len(sender.Messages()) == 1
	})

	t.Logf("Result %t, waited %s", check, aw.TimeWaited())
}
