package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/sentiweb/monitor-lib/notify/formatter"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/utils"
)

// PayloadFunc is the Payload() function signature for
type PayloadFunc func(n types.Notification, formatter types.Formatter) (interface{}, error)

// GenericHttpService describes a generic service sending a notification to a HTTP endpoint as a json
// It's a base implementation for many webhooks accepting message as a json payload
// Specific part of the webhook can be provided by Payload(), Check(), and Request() functions
type GenericHttpService struct {
	URL    string
	Method string
	Name   string

	Payload PayloadFunc                  // Create the  payload content
	Check   func(r *http.Response) error // Check the returned response
	Request func(req *http.Request)

	// Private field
	formatter types.Formatter
}

// Creates a new GenericHttpService
// Name is internal, and it's used to get a formatter from the FormatterFactory
func NewGenericHttpService(name string, URL string, options ...func(*GenericHttpService)) *GenericHttpService {
	h := &GenericHttpService{
		Name:    name,
		URL:     URL,
		Method:  http.MethodPost,
		Payload: defaultPayloadFunc,
		Request: defaultRequestFunc,
		Check:   defaultCheckFunc,
	}
	for _, o := range options {
		o(h)
	}
	return h
}

// WithCheckFunc defines the Check function used to check the response of the webhook after the request is sent
func WithCheckFunc(check func(r *http.Response) error) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Check = check
	}
}

// WithRequestFunc defines the Request function, called before to send the request
// It can be used by a service to defines specific headers of the service
func WithRequestFunc(Request func(req *http.Request)) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Request = Request
	}
}

// WithMethod defines the HTTP method to use (default is POST request)
func WithMethod(method string) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Method = method
	}
}

// WithMethod defines PayLoad() function, called with the Notification and formatter
// Payload() creates the object, to be serialized into json to be sent to the webhook
func WithPayload(payload PayloadFunc) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Payload = payload
	}
}

func defaultPayloadFunc(n types.Notification, formatter types.Formatter) (interface{}, error) {
	title := formatter.Title(n)
	text := formatter.Text(n)
	return map[string]string{"title": title, "message": text}, nil
}

func defaultRequestFunc(req *http.Request) {
	// Do Nothing
}

func defaultCheckFunc(r *http.Response) error {
	return nil
}

// Send the notification to the webhook using an HTTP request
func (g *GenericHttpService) Send(ctx context.Context, client utils.HTTPClient, n types.Notification) error {
	payload, err := g.Payload(n, g.formatter)

	if err != nil {
		return err
	}

	var body []byte

	body, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	var req *http.Request

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, g.URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add(utils.HttpHeaderContentType, utils.MimeJson)

	g.Request(req)

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	err = g.Check(resp)
	return err
}

// Start the service
func (g *GenericHttpService) Start() error {
	g.formatter = formatter.Get(g.Name)
	return nil
}
