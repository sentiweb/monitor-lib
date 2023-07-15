package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"github.com/sentiweb/monitor-lib/utils"
	"github.com/sentiweb/monitor-lib/notify/types"
	"github.com/sentiweb/monitor-lib/notify/formatter"
	
)

type PayloadFunc func(n types.Notification, formatter types.Formatter) (interface{}, error)

// GenericHttpService send a json payload 
type GenericHttpService struct {
	URL string
	Method string
	Name string
	
	Payload PayloadFunc // Create the  payload content
	Check func(r *http.Response) error // Check the returned response
	Request func(req *http.Request)

	// Private field
	formatter types.Formatter
}

func NewGenericHttpService(name string, URL string, options ...func(*GenericHttpService)) *GenericHttpService {
	h := &GenericHttpService{
		Name: name,
		URL: URL,
		Method: http.MethodPost,
		Payload: defaultPayloadFunc,
		Request: defaultRequestFunc,
		Check: defaultCheckFunc,
	}
	for _, o := range options {
		o(h)
	}
	return h
}

func WithCheckFunc(check func(r *http.Response) error) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Check = check
	}
}

func WithRequestFunc(Request func(req *http.Request)) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Request = Request
	}
}

func WithMethod(method string) func(*GenericHttpService) {
	return func(h *GenericHttpService) {
		h.Method = method
	}
}

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

func (g *GenericHttpService) Send(ctx context.Context, client utils.HTTPClient, n types.Notification) (error) {
	payload, err := g.Payload(n, g.formatter)
	
	if err != nil {
		return err
	}

	var body []byte

	body, err = json.Marshal(payload)
	if(err != nil) {
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

func (g *GenericHttpService) Start() error {
	g.formatter = formatter.Get(g.Name)
	return nil
}