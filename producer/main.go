package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/asaskevich/govalidator"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/usernamesalah/news/producer/api/v1"
	_ "github.com/usernamesalah/news/producer/api/v1/doc"
	"github.com/usernamesalah/news/producer/pkg/services"
)

// @title Palembang Digital API
// @version 1.0.0
// @description API documentation for palembangdigital.org

// @contact.name Palembang Digital
// @contact.url https://palembangdigital.org
// @contact.email support@palembangdigital.org

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
	kafkaProducer, err := KafkaInit(cfg)
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			return
		}
	}()

	log.Println("Initializing services ...")
	newsService := services.NewNewsService(db)
	kafkaService := services.NewKafkaService(kafkaProducer)

	log.Println("Initializing the web server ...")
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.Validator = &requestValidator{}

	// Utility endpoints
	e.GET("/docs/api/v1/indehtml", echoSwagger.WrapHandler)
	e.GET("/docs/api/v1/doc.json", echoSwagger.WrapHandler)
	e.GET("/docs/api/v1/*", echoSwagger.WrapHandler)
	e.GET("/ping", ping)

	// Serve API
	api := api.NewAPI(newsService, kafkaService, cfg.AdminUsername, cfg.AdminPassword)
	api.Register(e.Group("/api/v1", middleware.Logger()))

	// Start server
	s := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	e.Logger.Fatal(e.StartServer(s))
}

type requestValidator struct{}

func (rv *requestValidator) Validate(i interface{}) (err error) {
	_, err = govalidator.ValidateStruct(i)
	return
}

// ping write pong to http.ResponseWriter.
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// KafkaInit for starting Kafka Producer
func KafkaInit(cfg Config) (sarama.SyncProducer, error) {

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	if cfg.Kafka.User != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = cfg.Kafka.User
		kafkaConfig.Net.SASL.Password = cfg.Kafka.Password
	}

	kafkaProducer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return kafkaProducer, nil
}
