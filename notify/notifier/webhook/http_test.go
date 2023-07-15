package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"net/http"
	"testing"
	"time"
	"github.com/sentiweb/monitor-lib/notify/notifier/webhook/generic"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/utils"
	"github.com/sentiweb/monitor-lib/tests"
)

// MockClient is the mock client
type MockClient struct {
	StatusCode int
	Body       string
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

	defer tests.CaptureLog(t).Release()

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

	notif := types.NewMockNotification("up", "test12341", time.Now())

	notifier.Send(context, notif)

	aw := tests.NewAwait(time.Second, 100*time.Millisecond)

	check := aw.Wait(func() bool {
		return strings.Contains(client.RequestBody, "test12341")
	})

	t.Logf("Result %t, waited %s", check, aw.TimeWaited())
}
