package evhub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/sirupsen/logrus"
	consumerPb "ms-consumer/proto/consumer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

// evhub implements gRPC server
type evhub struct {
	br      broker.Broker
	name    string
	version string

	started time.Time
}

type healthCheckResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// PublishEvent publish an event to broker
func (e *evhub) PublishEvent(ctx context.Context, req *consumerPb.PublishEventRequest) (*consumerPb.PublishEventResponse, error) {

	// validate
	if req.GetBody() == "" || req.GetTopic() == "" {
		return nil, status.New(codes.Unavailable, "bad request").Err()
	}

	message := &broker.Message{
		Header: make(map[string]string),
		Body:   nil,
	}

	if req.GetMessageId() == "" {
		req.MessageId = uuid.New().String()
	}

	message.Header["x-message-id"] = req.GetMessageId()
	message.Body, _ = json.Marshal(req.GetBody())

	err := e.br.Publish(req.GetTopic(), message)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %s", err.Error())
	}
	rsp := &consumerPb.PublishEventResponse{
		MessageId: req.GetMessageId(),
	}
	return rsp, nil
}

// HttpHandler ...
func (e *evhub) HttpHandler(c *gin.Context) {
	req := &consumerPb.PublishEventRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.MessageId == "" && c.GetHeader("x-request-id") != "" {
		req.MessageId = c.GetHeader("x-request-id")
	}

	resp, err := e.PublishEvent(c, req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Info(log.WithField("request", req))
	c.JSON(http.StatusCreated, resp)
}

// HealthCheckHandler ...
func (e *evhub) HealthCheckHandler(c *gin.Context) {
	rsp := healthCheckResponse{
		Name:    e.name,
		Version: e.version,
		Uptime:  time.Now().Sub(e.started).String(),
	}
	c.JSON(http.StatusOK, rsp)
}

// EvHubService ...
type EvHubService interface {
	consumerPb.ConsumerServer

	HealthCheckHandler(c *gin.Context)
	HttpHandler(c *gin.Context)
}

// NewEvHubService initializes new event hub service
func NewEvHubService(br broker.Broker, name, version string) EvHubService {
	return &evhub{
		br:      br,
		name:    name,
		version: version,
		started: time.Now(),
	}
}
