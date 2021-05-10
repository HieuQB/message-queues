// This event handler is using for listen to broker messages and deliver messages to subscribed services that registered.
// firstly, we support http protocol, gRPC could come later
// the URL handlers should be able to handle requests like
// POST /consume
// content-type: application/json
//
// {
//   "event": "order.created",
//   "entity": "order",
//   "attributes": {"id": "order-id", "status": "created"}
// }
//
// the attributes will be used for notification trigger
//

package consumer

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/sirupsen/logrus"
	proto "ms-consumer/proto/consuming"
	"time"
)

var (
	_ EventHandler = &serviceEventHandler{}
)

// serviceEventHandler an service event handler
type serviceEventHandler struct {
	queue string
	topic string
	url   string

	client Client

	gClient proto.ConsumingServiceClient
}

// Handle ...
func (h *serviceEventHandler) Handle(p broker.Event) error {
	// check x-id header
	if _, exists := p.Message().Header["x-id"]; !exists {
		p.Message().Header["x-id"] = uuid.New().String()
	}

	var err error
	ctxLog := log.WithField("id", p.Message().Header["x-id"])

	if h.gClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		reqPayload := &proto.ConsumeRequest{
			Id:    p.Message().Header["x-id"],
			Topic: p.Topic(),
			Body:  p.Message().Body,
		}
		ctxLog.
			WithField("topic", p.Topic()).
			WithField("body", string(reqPayload.Body)).
			WithField("proto", "gRPC").
			Debug("send data to target")
		_, err = h.gClient.Consume(ctx, reqPayload)
	} else {
		ctxLog.
			WithField("topic", p.Topic()).
			WithField("body", string(p.Message().Body)).
			WithField("proto", "http").
			Debug("send data to target")
		err = h.client.Request(h.url, p.Topic(), p.Message().Header["x-id"], p.Message().Body)
	}

	if err != nil {
		// todo: handle err while requesting to target service
		log.Printf("failed to consume message: %v\n", err)
	} else {

	}
	_ = p.Ack()

	return nil
}

func (h *serviceEventHandler) QueueName() string {
	return h.queue
}

func (h *serviceEventHandler) AutoAck() bool {
	return false
}

func (h *serviceEventHandler) Topic() string {
	return h.topic
}

// NewServiceHandler ...
func NewServiceHandler(queue, topic, url string, client Client, gClient proto.ConsumingServiceClient) EventHandler {

	h := &serviceEventHandler{
		queue:   queue,
		topic:   topic,
		url:     url,
		client:  client,
		gClient: gClient,
	}

	return h
}
