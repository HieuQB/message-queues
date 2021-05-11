package consumer

import (
	"context"
	"errors"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"ms-consumer/mocks"
	"testing"
)

func TestInitBroker(t *testing.T) {
	csmr := NewConsumer("amqp://guest:guest@local:5672/")
	impl, ok := csmr.(*consumerImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl)
	assert.NotNil(t, csmr.GetBroker())
}

type consumerTestSuite struct {
	suite.Suite
	cmer       Consumer
	mockBroker *mocks.Broker
	url        string
}

var defaultHandler = makeHandler("test", "", true)

func makeHandler(topic, queueName string, autoAck bool) *mocks.EventHandler {
	handler := &mocks.EventHandler{}
	handler.On("Topic").Return(topic)
	handler.On("QueueName").Return(queueName)
	handler.On("AutoAck").Return(autoAck)
	return handler
}

func TestConsumer(t *testing.T) {
	suite.Run(t, new(consumerTestSuite))
}

func (ts *consumerTestSuite) SetupTest() {
	ts.mockBroker = &mocks.Broker{}

	ts.mockBroker.On("Disconnect").Return(nil)

	// mock init broker function
	newBroker = func(opts ...broker.Option) broker.Broker {
		return ts.mockBroker
	}

	ts.url = "amqp://admin:admin@rabbitmq:5672/"

	ts.cmer = NewConsumer(ts.url)
}

func (ts *consumerTestSuite) TearDownTest() {
	ts.cmer.Stop()
}

func (ts *consumerTestSuite) mockSuccessConnectBroker() {
	rightURLMatcher := mock.MatchedBy(func(fn broker.Option) bool {
		options := &broker.Options{}
		fn(options)

		return options.Addrs[0] == ts.url
	})

	ts.mockBroker.On("Init", rightURLMatcher).Return(nil)
	ts.mockBroker.On("Connect").Return(nil)
}

func (ts *consumerTestSuite) TestAddHandlers() {
	handlers := []EventHandler{
		defaultHandler,
		defaultHandler,
		defaultHandler,
	}
	ts.cmer.AddHandlers(handlers...)

	ts.Equal(len(handlers), len(ts.cmer.(*consumerImpl).handlers))
}

func (ts *consumerTestSuite) TestStartConsumerErrorInit() {
	ts.mockBroker.On("Init", mock.AnythingOfType("broker.Option")).Return(errors.New("init error"))
	err := ts.cmer.Start()
	ts.Error(err)
}

func (ts *consumerTestSuite) TestErrorConnect() {
	ts.mockBroker.On("Init", mock.AnythingOfType("broker.Option")).Return(nil)
	ts.mockBroker.On("Connect").Return(errors.New("connect error"))
	err := ts.cmer.Start()
	ts.Error(err)
}

func (ts *consumerTestSuite) TestAllDefault() {
	ts.cmer.AddHandlers(defaultHandler)

	ts.mockSuccessConnectBroker()

	// call consume to subscribe handler
	ts.mockBroker.On(
		"Subscribe",
		"test",
		mock.AnythingOfType("broker.Handler"),
		mock.AnythingOfType("broker.SubscribeOption"),
	).Return(nil, nil)

	err := ts.cmer.Start()
	ts.NoError(err)

	ts.Assert()
}

func (ts *consumerTestSuite) TestSubscribeError() {
	ts.cmer.AddHandlers(defaultHandler)

	ts.mockSuccessConnectBroker()

	// call consume to subscribe handler
	ts.mockBroker.On(
		"Subscribe",
		"test",
		mock.AnythingOfType("broker.Handler"),
		mock.AnythingOfType("broker.SubscribeOption"),
	).Return(nil, errors.New("subscribe error"))

	err := ts.cmer.Start()
	ts.Error(err)

	ts.Assert()
}

// TestDurableQueue should set durable queue when queue name is not empty
func (ts *consumerTestSuite) TestEnableAutoAck() {
	ts.cmer.AddHandlers(makeHandler("test", "testing-queue", false))
	ts.mockSuccessConnectBroker()

	options := &broker.SubscribeOptions{AutoAck: true, Context: context.Background()}
	optionsMatcher := mock.MatchedBy(func(fn broker.SubscribeOption) bool {
		fn(options)

		return true
	})
	// call consume to subscribe handler
	ts.mockBroker.On(
		"Subscribe",
		"test",
		mock.AnythingOfType("broker.Handler"),
		optionsMatcher,
		optionsMatcher,
		optionsMatcher,
		optionsMatcher,
	).Return(nil, nil)

	_ = ts.cmer.Start()

	ts.Assert()
	ts.False(options.AutoAck)
	ts.Equal("testing-queue", options.Queue)
	ts.NotEqual(context.Background(), options.Context)
}
