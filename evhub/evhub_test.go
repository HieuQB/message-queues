package evhub

import (
	"context"
	"errors"
	"testing"

	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ms-consumer/mocks"
	proto "ms-consumer/proto/consumer"
)

func TestPublishEventSuccess(t *testing.T) {

	// publish success case
	br := &mocks.Broker{}
	hub := NewEvHubService(br, "consumer", "?")
	messageMatcher := func(message *broker.Message) bool {
		assert.Equal(t, "sample-message-id", message.Header["x-message-id"])
		return true
	}
	br.On("Publish", "user.profile_updated", mock.MatchedBy(messageMatcher)).Once().Return(nil)

	req := &proto.PublishEventRequest{
		Topic:     "user.profile_updated",
		Body:      `{"id": "user-1", "name": "Kingsley"}`,
		MessageId: "sample-message-id",
	}
	rsp, err := hub.PublishEvent(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, rsp.GetMessageId())
}

func TestPublishEventAutoFillXMessageID(t *testing.T) {
	br := &mocks.Broker{}
	hub := NewEvHubService(br, "consumer", "?")
	messageMatcher := func(message *broker.Message) bool {
		assert.Contains(t, message.Header, "x-message-id", "should append x-message-id if missing")

		return true
	}
	br.On("Publish", "user.profile_updated", mock.MatchedBy(messageMatcher)).Once().Return(nil)

	req := &proto.PublishEventRequest{
		Topic: "user.profile_updated",
		Body:  `{"id": "user-1", "name": "Kingsley"}`,
	}
	rsp, err := hub.PublishEvent(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, rsp.GetMessageId())
}

func TestPublishEventFailedOnPublishBrokerError(t *testing.T) {
	br := &mocks.Broker{}
	hub := NewEvHubService(br, "consumer", "?")
	br.On("Publish", "user.profile_updated", mock.AnythingOfType(`*broker.Message`)).Once().Return(errors.New("error"))

	req := &proto.PublishEventRequest{
		Topic: "user.profile_updated",
		Body:  `{"id": "user-1", "name": "Kingsley"}`,
	}
	rsp, err := hub.PublishEvent(context.Background(), req)
	assert.Errorf(t, err, "error")
	assert.Empty(t, rsp)
}

func TestPublishEventFailedValidation(t *testing.T) {
	req := &proto.PublishEventRequest{}
	_, err := NewEvHubService(nil, "", "").PublishEvent(nil, req)
	assert.Error(t, err)
}
