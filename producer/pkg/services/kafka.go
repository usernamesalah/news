package services

import (
	"log"

	"github.com/Shopify/sarama"
)

// KafkaService service interface.
type KafkaService interface {
	SendMessage(msg string) error
}

// KafkaProducer hold kafka producer session
type kafkaProducer struct {
	Producer sarama.AsyncProducer
	Topic    string
}

// NewKafkaService returns an initialized KafkaService implementation.
func NewKafkaService(Producer sarama.AsyncProducer, topic string) KafkaService {
	return &kafkaProducer{Producer: Producer, Topic: topic}
}

// SendMessage function to send message into kafka
func (p *kafkaProducer) SendMessage(msg string) error {
	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.StringEncoder(msg),
	}

	select {
	case p.Producer.Input() <- kafkaMsg:
		log.Printf("Produced topic : %s", p.Topic)
	case err := <-p.Producer.Errors():
		log.Printf("Fail producing topic : %s error : %v", p.Topic, err)
		return err
	}

	return nil
}
