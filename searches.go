package splunk

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type searches struct{}

var (
	Searches = &searches{}
)

type (
	SearchesItem struct {
		Name    string    `json:"name"`
		Updated time.Time `json:"updated"`
		Author  string    `json:"author"`
		Content struct {
			SplunkHecTarget  string `json:"action.forward_alert_to_splunk_hec.param.splunk_hec_target"`
			ActionWebhookURL string `json:"action.webhook.param.url"`
			Actions          string `json:"actions"`
			CronSchedule     string `json:"cron_schedule"`
			Disabled         bool   `json:"disabled"`
			IsScheduled      bool   `json:"is_scheduled"`
			Search           string `json:"search"`
		} `json:"content"`
	}

	SearchRsp struct {
		Updated time.Time      `json:"updated"`
		Entry   []SearchesItem `json:"entry"`
		Paging  struct {
			Total int `json:"total"`
		} `json:"paging"`
	}
)

func (s *searches) List() (*SearchRsp, error) {
	if client == nil {
		return nil, ErrClientNotInit
	}
	_rsp := &SearchRsp{}
	rsp, err := client.R().
		SetDebug(clientDebug).
		SetResult(&_rsp).
		SetQueryParamsFromValues(url.Values{
			"f": []string{
				"search",
				"actions",
				"action.webhook",
				"action.webhook.param.url",
				"action.forward_alert_to_splunk_hec",
				"action.forward_alert_to_splunk_hec.param.splunk_hec_target",
				"cron_schedule",
				"is_scheduled",
				"disabled",
			},
		}).
		SetQueryParams(map[string]string{
			"output_mode": "json",
			"count":       "500",
		}).
		Get("/services/saved/searches")
	if err != nil {
		return nil, err
	}
	if !rsp.IsSuccess() {
		err = errors.Errorf("%s %s", rsp.Status(), rsp.String())
		return nil, err
	}
	return _rsp, nil
}
