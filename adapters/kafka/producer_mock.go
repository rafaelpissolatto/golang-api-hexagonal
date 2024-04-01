package kafka

// MessageProducerMock kafka message producer mock
type MessageProducerMock struct{}

var (
	ProduceMessageFunc func(topicName, value, eventName, traceID string)
)

// ProduceMessage is the produce message mock for ProduceMessage func
func (mp *MessageProducerMock) ProduceMessage(topicName, value, eventName, traceID string) {
	ProduceMessageFunc(topicName, value, eventName, traceID)
}
