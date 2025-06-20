package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/log"

	appError "github.com/fireops-software/fireops-edge-sevenio-notifier/error"
)

type SevenIoApi struct {
	ctx    context.Context
	logger log.ILogger
	url    string
	apiKey string

	maxRetries int
	timeout    time.Duration
}

type sevenIoSendPayload struct {
	To   string `json:"to"`
	Text string `json:"text"`
	From string `json:"from"`
}

func (s *SevenIoApi) SendSms(from string, to string, text string) chan async.ActionResult[any] {
	r := make(chan async.ActionResult[any])
	go func() {
		retries := 0
		var err error
		for retries < s.maxRetries {
			err = s.sendSms(
				sevenIoSendPayload{
					To:   to,
					From: from,
					Text: text,
				},
			)
			if err == nil {
				break
			}
			s.logger.Errorf("%v", err)
		}
		r <- async.ActionResult[any]{
			Result: true,
			Error:  err,
		}
	}()
	return r
}

func (s *SevenIoApi) sendSms(msg sevenIoSendPayload) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return appError.NewErrSevenIo("failed to marshal message as json - %v", err)
	}
	reqCtx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(
		reqCtx,
		"POST",
		fmt.Sprintf("%s/api/sms", s.url),
		bytes.NewBufferString(string(body)),
	)
	if err != nil {
		return appError.NewErrSevenIo("%v", err)
	}
	req.Header.Add("X-Api-Key", s.apiKey)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return appError.NewErrSevenIo("%v", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return appError.NewErrSevenIo("failed to send sms (StatusCode: %d)", resp.StatusCode)
	}
	return nil
}

func WithSevenIoMaxRetries(retries int) func(*SevenIoApi) {
	return func(sia *SevenIoApi) {
		sia.maxRetries = retries
	}
}

func WithSevenIoTimeout(timeout time.Duration) func(*SevenIoApi) {
	return func(sia *SevenIoApi) {
		sia.timeout = timeout
	}
}

func NewSevenIoApi(ctx context.Context, logger log.ILogger, apiUrl string, apiKey string, opts ...func(*SevenIoApi)) ISevenIoApi {
	s := &SevenIoApi{
		ctx:    ctx,
		logger: logger,
		url:    apiUrl,
		apiKey: apiKey,

		maxRetries: 3,
		timeout:    10 * time.Second,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}
