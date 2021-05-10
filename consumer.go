// to listen an event, you need to make a Handler that specify some settings for topic.
//
//  package main
//
//  import (
//    "fmt"
//    "gitlab.com/jx-adt-www/central/consumer"
//  )
//
//  type EchoHandler struct{}
//
//  // QueueName returns empty if you want to use generated random queue name
//  func (h* EchoHandler) QueueName() string {
//    return ""
//  }
//
//  // AutoAck set autoAck mode, if disable this (return false)
//  // you need to ack on received events manually
//  func (h *EchoHandler) AutoAck() bool { return false }
//
//  // Handle process the incoming revent
//  func (h *EchoHandler) Handle(p consumer.Event) {
//     fmt.Printf("received message: %s", p.Message().Body)
//     _ = p.Ack()
//  }
//
//  func main() {
//    consumer := consumer.NewConsumer("amqp://guest:guest@localhost:5672/")
//    consumer.AddHandlers(new(EchoHandler))
//    if err := consumer.Start(); err != nil {
//        fmt.Printf("failed to start consuming: %v", err)
//        return
//    }
//    fmt.Println("listening ...")
//  }
//
package consumer

import (
	libbroker "github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
)

var newBroker = rabbitmq.NewBroker

//go:generate mockery -name Broker
//go:generate mockery -name EventHandler
//go:generate mockery -name Event

// Broker ...
type Broker interface {
	libbroker.Broker
}

// Event ...
type Event libbroker.Event

// Consumer describes an consumer
type Consumer interface {
	Start() error
	AddHandlers(handlers ...EventHandler)
	Stop()
	GetBroker() Broker
}

// EventHandler describe an event handler
type EventHandler interface {
	Handle(p libbroker.Event) error
	Topic() string
	QueueName() string
	AutoAck() bool
}

type consumerImpl struct {
	url      string
	handlers []EventHandler
	broker   Broker
}

// Start starts running the consumer
func (c *consumerImpl) Start() error {
	err := c.broker.Init(libbroker.Addrs(c.url))
	if err != nil {
		return err
	}

	err = c.broker.Connect()
	if err != nil {
		return err
	}

	return c.consume()
}

// AddHandlers add handlers to consumer
func (c *consumerImpl) AddHandlers(handlers ...EventHandler) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *consumerImpl) consume() error {
	for _, h := range c.handlers {
		options := []libbroker.SubscribeOption{
			libbroker.DisableAutoAck(),
		}

		queueName := h.QueueName()
		if len(queueName) > 0 {
			options = append(options, rabbitmq.DurableQueue(), libbroker.Queue(queueName))
		}
		if !h.AutoAck() {
			options = append(options, libbroker.DisableAutoAck())
		}

		if _, err := c.broker.Subscribe(h.Topic(), h.Handle, options...); err != nil {
			return err
		}
	}
	return nil
}

func (c *consumerImpl) GetBroker() Broker {
	return c.broker
}

// Stop ...
func (c *consumerImpl) Stop() {
	_ = c.broker.Disconnect()
}

// NewConsumer ...
func NewConsumer(brokerURL string) Consumer {
	return &consumerImpl{
		url:      brokerURL,
		handlers: nil,
		broker:   newBroker(),
	}
}
