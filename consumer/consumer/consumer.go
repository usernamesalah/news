package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/usernamesalah/news/consumer/pkg/models"
	"github.com/usernamesalah/news/consumer/pkg/services"
)

// KafkaConsumer hold sarama consumer
type KafkaConsumer struct {
	Consumer      sarama.Consumer
	ElasticSearch services.ElasticService
	NewsService   services.NewsService
}

// NewAPI returns an initialized API type.
func NewConsumer(newsService services.NewsService, elasticService services.ElasticService, consumer sarama.Consumer) *KafkaConsumer {
	return &KafkaConsumer{
		NewsService:   newsService,
		ElasticSearch: elasticService,
		Consumer:      consumer,
	}
}

// Consume function to consume message from apache kafka
func (c *KafkaConsumer) Consume(topics []string, signals chan os.Signal) {
	chanMessage := make(chan *sarama.ConsumerMessage, 256)

	for _, topic := range topics {
		partitionList, err := c.Consumer.Partitions(topic)
		if err != nil {
			log.Printf("Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(c.Consumer, topic, partition, chanMessage)
		}
	}

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			log.Printf("New Message from kafka, message: %v", string(msg.Value))
			_ = c.processMessage(msg.Value)
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			log.Printf("Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}

}

// processMessage from kafka
func (c *KafkaConsumer) processMessage(message []byte) error {
	ctx := context.Background()

	var news models.News
	err := json.Unmarshal(message, &news)
	if err != nil {
		return err
	}

	newsDB, err := c.NewsService.CreateNews(ctx, news)
	if err != nil {
		return err
	}

	newsElastic, err := c.ElasticSearch.Create(ctx, newsDB)
	if err != nil {
		return err
	}
	fmt.Println("success insert to elastic : ", newsElastic)

	return nil
}
