package main

import "github.com/kelseyhightower/envconfig"

// Config stores the application configurations.
type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	UIBuildPath string `envconfig:"UI_BUILD_PATH" default:"ui/build"`

	AdminUsername string `envconfig:"ADMIN_USERNAME" default:"admin"`
	AdminPassword string `envconfig:"ADMIN_PASSWORD" default:"admin"`

	Database DatabaseConfig
	Kafka    KafkaConfig
	Elastic  ElasticConfig
}

// DatabaseConfig stores database configurations.
type DatabaseConfig struct {
	URL            string `envconfig:"DATABASE_URL" required:"true"`
	Driver         string `envconfig:"DATABASE_DRIVER" default:"postgres"`
	MigrationsPath string `envconfig:"DATABASE_MIGRATIONS_PATH" required:"true" default:"file://migrations/postgresql"`
}

//KafkaConfig stores kafka configurations.
type KafkaConfig struct {
	Host     string `envconfig:"KAFKA_HOST" default:"kafka" required:"true"`
	Port     string `envconfig:"KAFKA_PORT" default:"9092" required:"true"`
	User     string `envconfig:"KAFKA_USER"`
	Password string `envconfig:"KAFKA_PASSWORD"`
}

//ElasticConfig stores elastic configurations.
type ElasticConfig struct {
	Host     string `envconfig:"ELASTIC_HOST" default:"elastic" required:"true"`
	Port     string `envconfig:"ELASTIC_PORT" default:"9200" required:"true"`
	Username string `envconfig:"ELASTIC_USERNAME"`
	Password string `envconfig:"ELASTIC_PASSWORD"`
	Index    string `envconfig:"ELASTIC_INDEX"`
}

// ReadConfig populates configurations from environment variables.
func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
