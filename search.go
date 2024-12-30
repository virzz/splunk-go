package splunk

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/virzz/vlog"
)

type (
	QuerySid struct {
		Sid string `json:"sid"`
	}

	QueryResults []map[string]string
	QueryRsp     struct {
		Results QueryResults `json:"results"`
	}

	JobStatusEntryCount struct {
		IsDone      bool  `json:"isDone"`
		DiskUsage   int64 `json:"diskUsage"`
		ResultCount int64 `json:"resultCount"`
	}
	JobStatus struct {
		Entry []struct {
			Content JobStatusEntryCount `json:"content"`
		} `json:"entry"`
	}
	OutputMode string
)

const (
	OutputModeCSV   OutputMode = "csv"
	OutputModeJSON  OutputMode = "json"
	defaultTryTimes            = 5
)

type search struct {
}

var (
	Search = &search{}
)

func (s *search) Query(query string, id ...string) (string, error) {
	if client == nil {
		return "", ErrClientNotInit
	}
	data := map[string]string{"search": query}
	if len(id) > 0 {
		data["id"] = id[0]
	}
	res := QuerySid{}
	rsp, err := client.R().
		SetDebug(clientDebug).
		SetResult(&res).
		SetQueryParam("output_mode", "json").
		SetFormData(data).
		Post("/services/search/jobs")
	if err != nil {
		return "", err
	}
	if !rsp.IsSuccess() {
		err = errors.Errorf("%s %s", rsp.Status(), rsp.String())
		return "", err
	}
	if res.Sid == "" {
		return "", errors.New("Splunk Search Job ID is empty")
	}
	vlog.Info("Splunk Search Job Start", "sid", res.Sid)
	return res.Sid, nil
}

func (s *search) Status(sid string, times ...int) (r *JobStatusEntryCount, err error) {
	if client == nil {
		return nil, ErrClientNotInit
	}
	var (
		res = JobStatus{}
		rsp *resty.Response
	)
	tryTimes := defaultTryTimes
	if len(times) > 0 {
		tryTimes = times[0]
	}
	tryCount := 0
	for {
		select {
		case <-gCtx.Done():
			return nil, gCtx.Err()
		default:
			rsp, err = client.R().
				SetDebug(clientDebug).
				SetResult(&res).
				SetQueryParam("output_mode", "json").
				SetPathParam("sid", sid).
				Get("/services/search/jobs/{sid}")
			if err == nil && rsp.IsSuccess() && len(res.Entry) > 0 &&
				res.Entry[0].Content.IsDone &&
				res.Entry[0].Content.ResultCount > 0 {
				r = &res.Entry[0].Content
				vlog.Info("Splunk Search Job Done", "sid", sid, "count", r.ResultCount, "size", r.DiskUsage)
				break
			}
			tryCount++
		}
		if r != nil {
			break
		}
		vlog.Info("Try to get job status", "count", tryCount)
		time.Sleep(1 * time.Second)
		if tryTimes > 0 && tryCount >= tryTimes {
			return nil, errors.New("Splunk Search Job Timeout")
		}
	}
	return r, nil
}

func (s *search) Results(sid string) (*QueryResults, error) {
	if client == nil {
		return nil, ErrClientNotInit
	}
	item := &QueryRsp{}
	rsp, err := client.R().
		SetDebug(clientDebug).
		SetResult(&item).
		SetQueryParam("output_mode", "json").
		SetPathParam("sid", sid).
		Get("/services/search/jobs/{sid}/results")
	if err != nil {
		vlog.Error("Failed to get results", "err", err.Error())
		return nil, err
	}
	if !rsp.IsSuccess() {
		err = errors.Errorf("%s %s", rsp.Status(), rsp.String())
		vlog.Error("Failed to get result", "err", err)
		return nil, err
	}
	return &item.Results, nil
}

func (s *search) Download(sid, filename string, outputMode OutputMode) error {
	if client == nil {
		return ErrClientNotInit
	}
	rsp, err := client.R().
		SetDebug(clientDebug).
		SetQueryParams(map[string]string{
			"output_mode": string(outputMode),
			"count":       "0",
		}).
		SetPathParam("sid", sid).
		SetOutput(filename).
		Get("/services/search/jobs/{sid}/results")
	if err != nil {
		vlog.Error("Failed to get result", "err", err.Error())
		return err
	}
	if !rsp.IsSuccess() {
		err = errors.Errorf("%s %s", rsp.Status(), rsp.String())
		vlog.Error("Failed to get result", "err", err.Error())
		return err
	}
	vlog.Infof("Save result to %s", filename)
	return nil
}

func (s *search) QueryAndResults(query, filename string, times int, outputMode OutputMode) (*QueryResults, error) {
	sid, err := s.Query(query)
	if err != nil {
		return nil, err
	}
	stats, err := s.Status(sid, 10)
	if err != nil {
		return nil, err
	}
	if !stats.IsDone {
		return nil, ErrJobNotDone
	}
	return s.Results(sid)
}

func (s *search) QueryAndDownload(query, filename string, times int, outputMode OutputMode) error {
	sid, err := s.Query(query)
	if err != nil {
		return err
	}
	stats, err := s.Status(sid, times)
	if err != nil {
		return err
	}
	if !stats.IsDone {
		return ErrJobNotDone
	}
	err = s.Download(sid, filename, outputMode)
	if err != nil {
		return err
	}
	return nil
}
