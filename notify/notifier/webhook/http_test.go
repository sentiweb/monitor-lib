package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sentiweb/monitor-lib/notify/tests"

	"github.com/sentiweb/monitor-lib/notify/notifier/webhook/generic"
	utils_tests "github.com/sentiweb/monitor-lib/tests"
	"github.com/sentiweb/monitor-lib/utils"
)

// MockClient is the mock client
type MockClient struct {
	StatusCode  int
	Body        string
	RequestBody string
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Println(req)
	b, _ := io.ReadAll(req.Body)
	m.RequestBody = string(b)
	r := io.NopCloser(bytes.NewReader([]byte(m.Body)))
	return &http.Response{
			StatusCode: m.StatusCode,
			Body:       r,
		},
		nil
}

func (m *MockClient) String() string {
	return fmt.Sprintf("MockClient<%d>", m.StatusCode)
}

type MockClientFactory struct {
	client *MockClient
}

func (f *MockClientFactory) NewClient(timeout time.Duration, config *utils.HTTPClientParams) utils.HTTPClient {
	return f.client
}

var _ utils.HTTPClientFactory = &MockClientFactory{}

func TestHTTPNotifier(t *testing.T) {

	defer utils_tests.CaptureLog(t).Release()

	myFactory := &MockClientFactory{}

	client := &MockClient{
		StatusCode: 200,
		Body:       "ok",
	}

	myFactory.client = client

	HttpFactory = myFactory

	svc := generic.NewGenericHttpService("generic", "localhost")

	notifier := NewHTTPNotifier(svc)

	context := context.WithValue(context.TODO(), utils.ContextDebug, true)

	notifier.Start(context)

	fmt.Println(notifier)

	serviceName := "test12341"
	notif := tests.NewMockNotification(serviceName, "up", serviceName, time.Now())

	notifier.Send(context, notif)

	aw := utils_tests.NewAwait(time.Second, 100*time.Millisecond)

	check := aw.Wait(func() bool {
		return strings.Contains(client.RequestBody, serviceName)
	})

	t.Logf("Result %t, waited %s", check, aw.TimeWaited())

	if !check {
		t.Error("Notification not sent")
	}
}
