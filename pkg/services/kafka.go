package services

import (
	"log"

	"github.com/Shopify/sarama"
)

// KafkaService service interface.
type KafkaService interface {
	SendMessage(topic, msg string) error
}

// KafkaProducer hold kafka producer session
type kafkaProducer struct {
	Producer sarama.SyncProducer
}

// NewKafkaService returns an initialized KafkaService implementation.
func NewKafkaService(Producer sarama.SyncProducer) KafkaService {
	return &kafkaProducer{Producer: Producer}
}

// SendMessage function to send message into kafka
func (p *kafkaProducer) SendMessage(topic, msg string) error {

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := p.Producer.SendMessage(kafkaMsg)
	if err != nil {
		return err
	}

	log.Printf("Send message success, Topic %v, Partition %v, Offset %d", topic, partition, offset)
	return nil
}
