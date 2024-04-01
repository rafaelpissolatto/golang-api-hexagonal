package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Configurations Configurations from config file
type Configurations struct {
	Server   ServerConfigurations   `yaml:"server"`
	Service  ServiceConfigurations  `yaml:"service"`
	DB       DatabaseConfigurations `yaml:"database"`
	Kafka    KafkaConfiguration     `yaml:"kafka"`
	Redis    RedisConfiguration     `yaml:"redis"`
	Oauth    Oauth                  `yaml:"oauth"`
	Policies PoliciesConfiguration  `yaml:"policies"`
}

// ServerConfigurations Server configurations
type ServerConfigurations struct {
	Port string `yaml:"port"`
}

// ServiceConfigurations Service configurations
type ServiceConfigurations struct {
	Name string `yaml:"name"`
}

// DatabaseConfigurations Database configurations
type DatabaseConfigurations struct {
	// DNS     string `yaml:"dns"`
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db-name"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Insecure bool   `yaml:"insecure"`
	Pool     int    `yaml:"pool"`
	Timeout  int    `yaml:"timeout"`
}

// KafkaConfiguration kafka connection and producer configuration
type KafkaConfiguration struct {
	SecurityProtocol string                     `yaml:"security-protocol"`
	Servers          string                     `yaml:"servers"`
	User             string                     `yaml:"user"`
	Pass             string                     `yaml:"pass"`
	ClientName       string                     `yaml:"client-name"`
	Producer         KafkaProducerConfiguration `yaml:"producer"`
	ConsumerEnabled  bool                       `yaml:"consumer-enabled"`
	Consumer         KafkaConsumerConfiguration `yaml:"consumer"`
}

// RedisConfiguration redis connection configuration
type RedisConfiguration struct {
	Localhost        bool   `yaml:"localhost"`
	URL              string `yaml:"url"`
	User             string `yaml:"user"`
	Pass             string `yaml:"pass"`
	DB               int    `yaml:"db"`
	PublicKeyFile    string `yaml:"public-key-file"`
	PrivateKeyFile   string `yaml:"private-key-file"`
	CaCertFile       string `yaml:"ca-cert-file"`
	TimeOutInSeconds int64  `yaml:"time-out-in-seconds"`
}

// KafkaProducerConfiguration kafka producer configuration
type KafkaProducerConfiguration struct {
	ProductTopic string `yaml:"product-topic-event"`
}

// KafkaConsumerConfiguration kafka consumer configuration
type KafkaConsumerConfiguration struct {
	Group      string   `yaml:"group"`
	Topics     []string `yaml:"topics"`
	MaxRecords int      `yaml:"max-records"`
}

// Oauth secret key
type Oauth struct {
	Secret string `yaml:"secret"`
}

// PoliciesConfiguration policies configuration
type PoliciesConfiguration struct {
	Path string `yaml:"path"`
}

// LoadConfigFile Load the yml config file and environment variables
func LoadConfigFile(log *zap.SugaredLogger) *Configurations {
	configFile, err := os.ReadFile("./resources/config.yml")
	if err != nil {
		log.Fatalf("Failed to load the config file: %v", err)
	}

	// expand environment variables
	confContent := []byte(os.ExpandEnv(string(configFile)))

	configurations := &Configurations{}
	err = yaml.Unmarshal(confContent, configurations)
	if err != nil {
		log.Fatalf("Failed to unmarshall the config file: %v", err)
	}

	return configurations
}
