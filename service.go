package consumer

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"ms-consumer/evhub"
	consumerPb "ms-consumer/proto/consumer"
	proto "ms-consumer/proto/consuming"

	log "github.com/sirupsen/logrus"
)

func init() {
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:  time.RFC3339,
	})
}

func logE(e error, tag, msg string) {
	log.WithError(e).Errorf("[%s] %s", tag, msg)
}

// Service ...
type Service interface {
	Start(stop chan struct{}) error
}

var (
	readFile = ioutil.ReadFile
)

type serviceImpl struct {
	config            AppConfig
	consumer          Consumer
	servicesConsumers ServicesConsumers

	gs  *grpc.Server
	hub evhub.EvHubService

	hs *http.Server

	name    string
	version string
}

// Start starts service
func (s *serviceImpl) Start(stop chan struct{}) error {
	client := NewClient()

	var grpcConns []*grpc.ClientConn
	defer func() {
		for _, c := range grpcConns {
			_ = c.Close()
		}
	}()

	// register services consumers
	for _, svc := range s.servicesConsumers.Services {
		for _, l := range svc.Listeners {

			localLog := log.WithFields(log.Fields{"service": svc.ServiceID, "topic": l.Topic})

			localLog.Debug("register service topic")

			// create new grpc client for this service (if configured)
			var gClient proto.ConsumingServiceClient
			if svc.GrpcAddr != "" {

				localLog.
					WithField("svc.GrpcAddr", svc.GrpcAddr).
					Debug("[consumer][gRPC] create connection")
				// we rarely see error here, don't need to handle
				if conn, err := grpc.Dial(svc.GrpcAddr, grpc.WithInsecure()); err == nil {
					gClient = proto.NewConsumingServiceClient(conn)
					grpcConns = append(grpcConns, conn)
				}

			}

			h := NewServiceHandler(l.Queue, l.Topic, svc.URL, client, gClient)
			s.consumer.AddHandlers(h)
			localLog.Info("register service topic successfully")
		}
	}

	// then start consumer
	if err := s.consumer.Start(); err != nil {
		return fmt.Errorf("[consumer.start] %s", err.Error())
	}
	defer func() {
		if s.consumer.GetBroker() != nil {
			_ = s.consumer.GetBroker().Disconnect()
		}
	}()

	s.hub = evhub.NewEvHubService(s.consumer.GetBroker(), s.name, s.version)

	// start gRPC server
	go func() {
		s.startGrpc()
	}()

	go func() {
		s.startHttp()
	}()

	<-stop

	s.gs.Stop()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	_ = s.hs.Shutdown(ctxTimeout)

	return nil
}

func (s *serviceImpl) startGrpc() {

	tag := "grpc.start"
	s.gs = grpc.NewServer()

	consumerPb.RegisterConsumerServer(s.gs, s.hub)
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.GrpcPort)
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		logE(err, tag, "failed to listen")
		return
	}
	log.WithField("addr", addr).Infof("[%s] %s", tag, "listening")
	if err := s.gs.Serve(lis); err != nil {
		logE(err, tag, "failed to serve")
	}
}

func (s *serviceImpl) startHttp() {

	tag := "http.start"

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.GET("/", s.hub.HealthCheckHandler)
	router.POST("/events", s.hub.HttpHandler)

	s.hs = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.config.HttpPort),
		Handler: router,
	}

	log.WithField("addr", s.hs.Addr).Infof("[%s] %s", tag, "listening")
	if err := s.hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logE(err, tag, "failed to serve")
	}
}

// NewService ...
func NewService(name, version string) Service {
	// load config
	config := AppConfig{}
	_ = env.Parse(&config)

	s := &serviceImpl{}
	s.config = config

	s.consumer = NewConsumer(config.BrokerURL)

	// load services consumers configuration
	raw, err := readFile(config.ServicesYamlPath)
	if err != nil {
		logE(err, "config", "failed to load config")
	}
	if err := yaml.Unmarshal(raw, &s.servicesConsumers); err != nil {
		logE(err, "services.unmarshal", "failed to unmarshal services")
	}

	s.name = name
	s.version = version

	return s
}
