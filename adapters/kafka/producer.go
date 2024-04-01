package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"time"
)

const eventNameKey = "event.name"

// MessageProducer message kafka producer
type MessageProducer struct {
	log      *zap.SugaredLogger
	producer *kafka.Producer
}

// NewKafkaProducer create the kafka producer connection
func NewKafkaProducer(log *zap.SugaredLogger, kafkaConfigMap *kafka.ConfigMap) *MessageProducer {
	producer, err := kafka.NewProducer(kafkaConfigMap)
	if err != nil {
		log.Panicf("Error to create the kafka producer: %s", err)
	}

	log.Infof("Kafka Producer Connecteded: %v", producer.Len())
	return &MessageProducer{
		log:      log,
		producer: producer,
	}
}

// Close the kafka producer
func (mp *MessageProducer) Close() {
	mp.producer.Close()
}

// ProduceMessage produce the message in the kafka topic
func (mp *MessageProducer) ProduceMessage(topicName, value, eventName, traceID string) {
	deliveryChan := make(chan kafka.Event, 10000)

	err := mp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Key:            []byte(traceID),
		Timestamp:      time.Now().UTC(),
		Headers:        []kafka.Header{{Key: eventNameKey, Value: []byte(eventName)}},
	}, deliveryChan)
	if err != nil {
		mp.log.With("traceId", traceID).Errorf("Internal error to send the kafka message in topic: %v, with error: %v", topicName, err)
	}

	// synchronous writes
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		mp.log.With("traceId", traceID).Errorf("Internal error to send the kafka message in partition topic: %v, with error: %v", topicName, err)
	} else {
		mp.log.With("traceId", traceID).Infof("Message delivered to topic: %v, partition: %v, at offset: %v",
			topicName, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
