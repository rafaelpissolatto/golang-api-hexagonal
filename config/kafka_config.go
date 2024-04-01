package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	protocol      = "security.protocol"
	plaintext     = "plaintext"
	saslSsl       = "sasl_ssl"
	saslPlaintext = "sasl_plaintext"
	plain         = "PLAIN"
	Producer      = "producer"
	Consumer      = "consumer"
	Topic         = "topic"
)

// NewKafkaConfigMap config the kafka connection properties
func NewKafkaConfigMap(log *zap.SugaredLogger, config KafkaConfiguration, configType string) *kafka.ConfigMap {
	var kafkaConf = &kafka.ConfigMap{
		"bootstrap.servers": config.Servers,
		"message.max.bytes": 1000000,
	}
	if configType == Producer {
		_ = kafkaConf.SetKey("retries", 5)
		_ = kafkaConf.SetKey("retry.backoff.ms", 1000)
		// Acks property controls how many partition replicas must acknowledge the receipt of a record before a producer can consider a particular write operation as successful.
		// acks = -1, the producer waits for the ack. Having the messages replicated to all the partition replicas.
		_ = kafkaConf.SetKey("acks", -1)
	}
	if configType == Consumer && config.ConsumerEnabled {
		_ = kafkaConf.SetKey("group.id", config.Consumer.Group)
		_ = kafkaConf.SetKey("auto.offset.reset", "latest")
		_ = kafkaConf.SetKey("heartbeat.interval.ms", 3000)
		_ = kafkaConf.SetKey("session.timeout.ms", 30000)
		_ = kafkaConf.SetKey("max.poll.interval.ms", 120000)
		_ = kafkaConf.SetKey("max.partition.fetch.bytes", 256000)
	}

	switch config.SecurityProtocol {
	case plaintext:
		_ = kafkaConf.SetKey(protocol, plaintext)
	case saslSsl:
		_ = kafkaConf.SetKey(protocol, saslSsl)
		_ = kafkaConf.SetKey("ssl.ca.location", "conf/ca-cert.pem")
		setSSLProperties(kafkaConf, &config)
	case saslPlaintext:
		_ = kafkaConf.SetKey(protocol, saslPlaintext)
		setSSLProperties(kafkaConf, &config)
	default:
		log.Panic(kafka.NewError(kafka.ErrUnknownProtocol, "Unknown kafka protocol", true))
	}

	return kafkaConf
}

func setSSLProperties(kafkaConf *kafka.ConfigMap, config *KafkaConfiguration) {
	_ = kafkaConf.SetKey("sasl.mechanism", plain)
	_ = kafkaConf.SetKey("sasl.username", config.User)
	_ = kafkaConf.SetKey("sasl.password", config.Pass)
}
