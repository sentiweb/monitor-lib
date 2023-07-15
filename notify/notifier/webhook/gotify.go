package webhook

import (
	"fmt"
	
	"net/url"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/notify/notifier/webhook/generic"
)

type GotifyMessage struct {
	Title string `json:"title"`
	Message string  `json:"message"`
}

func NewGotifyService(URL string, token string) types.WebhookNotifierService {
	
	params := url.Values{}
    params.Add("token", token)

	url := fmt.Sprintf("%s/message?%s", URL, params.Encode()) 
	g := generic.NewGenericHttpService("gotify", url, 
					generic.WithPayload(gotifyPayload),
				)	
	return g
}

func gotifyPayload(n types.Notification, formatter types.Formatter) (interface{}, error) {
	message := formatter.Text(n)
	title := formatter.Title(n)
	return GotifyMessage{Message: message, Title: title}, nil
}