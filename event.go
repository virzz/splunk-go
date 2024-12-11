package splunk

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var eventClient *EventClient

type EventConfig struct {
	Host   string `json:"host" yaml:"host" mapstructure:"host"`
	Index  string `json:"index" yaml:"index" mapstructure:"index"`
	Source string `json:"source" yaml:"source" mapstructure:"source"`
	Token  string `json:"token" yaml:"token" mapstructure:"token"`
}

type EventReq struct {
	Timestamp_   int64  `json:"timestamp"`
	Name_        string `json:"name,omitempty"`
	Description_ string `json:"description,omitempty"`
	App_         string `json:"app,omitempty"`
	Label_       string `json:"label,omitempty"`
	Events_      []any  `json:"events,omitempty"`
	Event_       any    `json:"event,omitempty"`
}

func (e *EventReq) Name(v string) *EventReq {
	e.Name_ = v
	return e
}

func (e *EventReq) Description(v string) *EventReq {
	e.Description_ = v
	return e
}

func (e *EventReq) Events(v any) *EventReq {
	e.Events_ = append(e.Events_, v)
	return e
}

func (e *EventReq) Event(v any) *EventReq {
	e.Event_ = v
	return e
}
func (e *EventReq) App(v string) *EventReq {
	e.App_ = v
	return e
}
func (e *EventReq) Label(v string) *EventReq {
	e.Label_ = v
	return e
}

func NewEventReq() *EventReq {
	return &EventReq{Events_: []any{}, Timestamp_: time.Now().Unix()}
}

type EventClient struct {
	client        *resty.Client
	index, source string
	debug         bool
	headers       map[string]string
}

func (ec *EventClient) Index(v string) *EventClient {
	ec.index = v
	return ec
}
func (ec *EventClient) Source(v string) *EventClient {
	ec.source = v
	return ec
}
func (ec *EventClient) Debug(v bool) *EventClient {
	ec.debug = v
	return ec
}
func (ec *EventClient) Headers(v map[string]string) *EventClient {
	ec.headers = v
	return ec
}

func (ec *EventClient) Send(req *EventReq) error {
	_req := ec.client.R().
		SetQueryParam("source", ec.source).
		SetQueryParam("index", ec.index).
		SetBody(req).
		SetDebug(ec.debug).
		SetHeaders(ec.headers)
	rsp, err := _req.Post("/services/collector/raw")
	if err != nil {
		log.Error("Failed to send event to splunk", "err", err.Error())
		return err
	}
	if !rsp.IsSuccess() {
		err = errors.Errorf("%s %s", rsp.Status(), rsp.String())
		log.Error("Failed to send event to splunk [not success]", "err", err.Error())
		return err
	}
	return nil
}

func NewEventClient(host, token string) *EventClient {
	return &EventClient{
		client: resty.New().
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
			SetRetryWaitTime(1*time.Second).
			SetRetryCount(3).
			SetBaseURL(host).
			SetAuthScheme("Splunk").
			SetAuthToken(token).
			SetHeader("Content-Type", "application/json").
			SetQueryParam("output_mode", "json"),
	}
}

func InitEvent(cfg *EventConfig) *EventClient {
	eventClient = NewEventClient(cfg.Host, cfg.Token).
		Index(cfg.Index).
		Source(cfg.Source)
	return eventClient
}

func Event(req *EventReq) error {
	return eventClient.Send(req)
}

func RawEvent(req any) error {
	return eventClient.Send(NewEventReq().Event(req))
}
func RawEvents(req any) error {
	return eventClient.Send(NewEventReq().Events(req))
}
