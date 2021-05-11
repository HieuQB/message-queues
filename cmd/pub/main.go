// this tool is just used for testing publishing message
package main

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/broker"
	"ms-consumer"
	"log"
	"os"
	"time"
)

type echoMessage struct {
	Time time.Time `json:"time"`
	Body string    `json:"body"`
}

func main() {
	brokerURL := os.Getenv("BROKER_URL")
	if brokerURL == "" {
		brokerURL = "amqp://admin:admin@rabbitmq:5672/"
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		topic = "echo"
	}

	message := os.Getenv("MESSAGE")
	if message == "" {
		message = "hello from publisher"
	}

	csmr := consumer.NewConsumer(brokerURL)
	if err := csmr.Start(); err != nil {
		log.Printf("failed to connect broker: %s", err)
		os.Exit(1)
	}
	defer func() {
		_ = csmr.GetBroker().Disconnect()
	}()

	msgPayload := echoMessage{
		Time: time.Now(),
		Body: message,
	}

	b, _ := json.Marshal(msgPayload)
	msg := &broker.Message{
		Header: nil,
		Body: b,
	}
	if err := csmr.GetBroker().Publish(topic, msg); err != nil {
		log.Printf("failed to publish message: %s", err)
	}
}
