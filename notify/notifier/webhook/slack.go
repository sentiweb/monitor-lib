package webhook

import (
	"bytes"
	"errors"
	"net/http"
	"github.com/sentiweb/monitor-lib/notify/notifier/webhook/generic"
	"github.com/sentiweb/monitor-lib/notify/types"
)


type SlackMessage struct {
	Text string `json:"text"`
}

func newSlackNotifierService(URL string) types.WebhookNotifierService {
	return generic.NewGenericHttpService("slack", URL,
		generic.WithPayload(slackPayload),
		generic.WithCheckFunc(slackCheckResponse),
	)
}

func slackPayload(n types.Notification, f types.Formatter) (interface{}, error) {
	return SlackMessage{Text: f.Text(n)}, nil
}

func slackCheckResponse(response *http.Response) error {
	buf := new(bytes.Buffer)
	defer response.Body.Close()
	buf.ReadFrom(response.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}