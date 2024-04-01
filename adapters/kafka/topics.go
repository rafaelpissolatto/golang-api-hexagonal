package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"golang-api-hexagonal/config"
)

// CreateKafkaTopics create the kafka topics
func CreateKafkaTopics(log *zap.SugaredLogger, config config.KafkaConfiguration, ctx context.Context, kafkaConfigMap *kafka.ConfigMap) {
	adminClient, err := kafka.NewAdminClient(kafkaConfigMap)
	if err != nil {
		log.Panicf("Error to create the kafka admin client: %s", err)
	}

	// Create the topics with 3 partitions and replication for 1 broker, and set 60 seconds of timeout. If you have more brokers in your cluster, you can set more than 1.
	results, err := adminClient.CreateTopics(
		ctx, []kafka.TopicSpecification{{
			Topic:             config.Consumer.Topics[0],
			NumPartitions:     3,
			ReplicationFactor: 1,
		}},
		kafka.SetAdminOperationTimeout(60000))
	if err != nil {
		log.Panicf("Error to create the kafka topics: %s", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Panicf("Error to create the kafka topic: %s, and error: %s", result.Topic, result.Error.String())
		}
	}

	log.Infof("Kafka topic created: %s", config.Consumer.Topics[0])
}
