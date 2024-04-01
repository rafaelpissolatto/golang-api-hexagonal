package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"golang-api-hexagonal/config"
	"golang-api-hexagonal/core/domain"
	"os"
	"os/signal"
	"syscall"
)

// NewKafkaConsumer create the kafka consumer connection
func NewKafkaConsumer(log *zap.SugaredLogger, kafkaConfigMap *kafka.ConfigMap) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(kafkaConfigMap)
	if err != nil {
		log.Panicf("Error to create the kafka producer: %s", err)
	}

	log.Infof("Kafka Consumer Connecteded")
	return consumer
}

// CloseConsumer Close the kafka consumer
func CloseConsumer(consumer *kafka.Consumer) {
	_ = consumer.Close()
}

// ConsumeMessages consumes all kafka messages from topics
func ConsumeMessages(log *zap.SugaredLogger, config config.KafkaConfiguration, consumer *kafka.Consumer) {
	err := consumer.SubscribeTopics(config.Consumer.Topics, nil)
	if err != nil {
		log.Panicf("Error to subscribe to the kafka topics: %s", err)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-signChan:
			log.Infof("Caught termination signal: %v", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Infof("Topic: %v, with message: %v, with header: %v", e.TopicPartition, string(e.Value), e.Headers)
				if e.Headers[0].Key == domain.ProductEventName {
					log.Infof("Product creation event")
					// Parse the value to a golang struct and develop some business rule.
				}
			case kafka.Error:
				log.Errorf("Error to read the message: %v, with code: %v", e.String(), e.Code())
			default:
				log.Infof("Unkwon message: %v", e)
			}
		}
	}
	close(signChan)
}
