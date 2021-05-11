package consumer

// Listener ...
type Listener struct {
	Queue string `yaml:"queue"`
	Topic string `yaml:"topic"`
}

// ServiceConfig ...
type ServiceConfig struct {
	ServiceID string     `yaml:"service_id"`
	URL       string     `yaml:"url"`
	Listeners []Listener `yaml:"listeners"`
	GrpcAddr  string     `yaml:"grpc_addr"`
}

// ServicesConsumers ...
type ServicesConsumers struct {
	Services []ServiceConfig `yaml:"services"`
}

// AppConfig ...
type AppConfig struct {
	HttpPort         int    `env:"HTTP_PORT" envDefault:"80"`
	BrokerURL        string `env:"BROKER_URL" envDefault:"amqp://admin:admin@localhost:5672/"`
	ServicesYamlPath string `env:"SERVICES_YAML_PATH" envDefault:"config/config.yml"`

	GrpcPort int `env:"GRPC_PORT" envDefault:"8081"`
}
