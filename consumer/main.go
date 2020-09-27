package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"

	"github.com/usernamesalah/news/consumer/consumer"
	"github.com/usernamesalah/news/consumer/pkg/services"
)

// @title News Api
// @version 1.0.0
// @description API documentation for News Api

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

func main() {

	log.Println("Reading the configuration from environment variables ...")
	cfg, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	log.Println("Migrating the database ...")
	m, err := migrate.New(cfg.Database.MigrationsPath, cfg.Database.URL)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	log.Println("Initializing the database connection ...")
	db, err := sqlx.Connect(cfg.Database.Driver, cfg.Database.URL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println("Initializing the kafka connection ...")
	kafka, err := KafkaInit(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := kafka.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	log.Println("Initializing the elasticsearch connection ...")
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(
			fmt.Sprintf("http://%s:%s", cfg.Elastic.Host, cfg.Elastic.Port),
		),
		elastic.SetBasicAuth(cfg.Elastic.Username, cfg.Elastic.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		panic(err)
	}

	log.Println("Initializing services ...")
	newsService := services.NewNewsService(db)
	elasticService := services.NewElasticService(elasticClient, cfg.Elastic.Index)
	kafkaConsumer := consumer.NewConsumer(newsService, elasticService, kafka)

	signals := make(chan os.Signal, 1)
	kafkaConsumer.Consume([]string{cfg.Kafka.Topic}, signals)
}

// KafkaInit for starting Kafka Consumer
func KafkaInit(cfg Config) (sarama.Consumer, error) {

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	if cfg.Kafka.User != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = cfg.Kafka.User
		kafkaConfig.Net.SASL.Password = cfg.Kafka.Password
	}

	kafkaConfig.Consumer.Return.Errors = true

	host := strings.Join([]string{cfg.Kafka.Host, cfg.Kafka.Port}, ":")

	kafkaConsumer, err := sarama.NewConsumer([]string{host}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return kafkaConsumer, nil
}
