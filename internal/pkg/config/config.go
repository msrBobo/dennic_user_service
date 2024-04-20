package config

import (
	"os"
	"strings"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	OTLPCollector struct {
		Host string
		Port string
	}

	MongoDb struct {
		MongoURI      string
		MongoDatabase string
	}

	Kafka struct {
		Address []string
		Topic   struct {
			InvestorCreate string
		}
	}

	PostService struct {
		PostServiceHost string
		PostServicePort string
	}

	CommentService struct {
		CommentServiceHost string
		CommentServicePort string
	}
}

func New() *Config {
	var c Config

	// general configuration
	c.APP = getEnv("APP", "app")
	c.Environment = getEnv("ENVIRONMENT", "develop")
	c.LogLevel = getEnv("LOG_LEVEL", "debug")
	c.RPCPort = getEnv("RPC_PORT", ":50025")
	c.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	c.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	c.DB.Port = getEnv("POSTGRES_PORT", "5544")
	c.DB.User = getEnv("POSTGRES_USER", "postgres")
	c.DB.Password = getEnv("POSTGRES_PASSWORD", "1234")
	c.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	c.DB.Name = getEnv("POSTGRES_DATABASE", "dennic_user_service")

	// otlp collector configuration
	c.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "localhost")
	c.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":2020")

	// kafka configuration
	c.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:29092"), ",")
	c.Kafka.Topic.InvestorCreate = getEnv("KAFKA_TOPIC_INVESTOR_CREATE", "investor.created")

	c.MongoDb.MongoURI = getEnv("MONGO_URI", "mongodb://localhost:27018")
	c.MongoDb.MongoDatabase = getEnv("MONGO_DATABASE", "userdb")

	c.PostService.PostServiceHost = getEnv("POST_SERVICE_HOST", "localhost")
	c.PostService.PostServicePort = getEnv("POST_SERVICE_PORT", "3030")

	c.CommentService.CommentServiceHost = getEnv("COMMENT_SERVICE_HOST", "localhost")
	c.CommentService.CommentServicePort = getEnv("COMMENT_SERVICE_PORT", "4040")

	return &c
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
