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
	Producer sarama.AsyncProducer
}

// NewKafkaService returns an initialized KafkaService implementation.
func NewKafkaService(Producer sarama.AsyncProducer) KafkaService {
	return &kafkaProducer{Producer: Producer}
}

// SendMessage function to send message into kafka
func (p *kafkaProducer) SendMessage(topic, msg string) error {
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	select {
	case p.Producer.Input() <- kafkaMsg:
		log.Printf("Produced topic : %s", topic)
	case err := <-p.Producer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topic, err)
		return err
	}

	return nil
}
