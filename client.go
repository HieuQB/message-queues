package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// Client service client that sends request to target service
type Client interface {
	Request(url string, topic string, xid string, payload []byte) error
}

// RequestPayload payload to send to target service
type RequestPayload struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

type clientImpl struct {
	httpClient *http.Client
}

func (c *clientImpl) Request(url string, topic string, xid string, payload []byte) error {

	log := logrus.WithField("func", "Client.Request").WithField("request_id", xid)

	requestPayload := RequestPayload{Topic: topic, Payload: payload}
	b, _ := json.Marshal(requestPayload)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("user-agent", "x-consumer/v1")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-request-id", xid)
	// todo: check error & implement retry on error
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to request")
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("failed to read response body")
		} else {
			log.WithField("response_code", resp.StatusCode).WithField("body", string(body))
		}
	}

	return nil
}

// NewClient ...
func NewClient() Client {
	c := &clientImpl{httpClient: &http.Client{}}

	return c
}
