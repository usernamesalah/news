package consumer_test

import (
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/usernamesalah/news/consumer/consumer"
)

func TestConsume(t *testing.T) {
	consumers := mocks.NewConsumer(t, nil)
	defer func() {
		if err := consumers.Close(); err != nil {
			t.Error(err)
		}
	}()

	consumers.SetTopicMetadata(map[string][]int32{
		"news": {0},
	})

	kafka := &consumer.KafkaConsumer{
		Consumer: consumers,
	}

	consumers.ExpectConsumePartition("news", 0, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world")})

	signals := make(chan os.Signal, 1)
	go kafka.Consume([]string{"news"}, signals)
	timeout := time.After(2 * time.Second)
	for {
		select {
		case <-timeout:
			signals <- os.Interrupt
			return
		}
	}
}
