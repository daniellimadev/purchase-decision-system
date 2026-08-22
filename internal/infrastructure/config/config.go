package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AWS      AWSConfig
	API      APIConfig
	Database DatabaseConfig
	Worker   WorkerConfig
	Business BusinessConfig
}

type AWSConfig struct {
	Region               string
	AccessKeyID          string
	SecretAccessKey      string
	SQSQueueURL          string
	SNSTopicARN          string
	SQSMaxMessages       int64
	SQSWaitTime          int64
	SQSVisibilityTimeout int64
}

type APIConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type WorkerConfig struct {
	Enabled      bool
	PollInterval time.Duration
}

type BusinessConfig struct {
	MaxPurchaseAmount    float64
	MinPurchaseAmount    float64
	AutoApproveThreshold float64
	AutoRejectThreshold  float64
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Warn("No .env.development file found, using environment variables")
	}

	config := &Config{
		AWS: AWSConfig{
			Region:               getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:          getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:      getEnv("AWS_SECRET_ACCESS_KEY", ""),
			SQSQueueURL:          getEnv("SQS_QUEUE_URL", ""),
			SNSTopicARN:          getEnv("SNS_TOPIC_ARN", ""),
			SQSMaxMessages:       getEnvAsInt64("SQS_MAX_MESSAGES", 10),
			SQSWaitTime:          getEnvAsInt64("SQS_WAIT_TIME_SECONDS", 20),
			SQSVisibilityTimeout: getEnvAsInt64("SQS_VISIBILITY_TIMEOUT", 30),
		},
		API: APIConfig{
			Port: getEnv("API_PORT", "8080"),
			Host: getEnv("API_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "purchase_decisions"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Worker: WorkerConfig{
			Enabled:      getEnvAsBool("WORKER_ENABLED", true),
			PollInterval: getEnvAsDuration("WORKER_POLL_INTERVAL", 5*time.Second),
		},
		Business: BusinessConfig{
			MaxPurchaseAmount:    getEnvAsFloat64("MAX_PURCHASE_AMOUNT", 10000.00),
			MinPurchaseAmount:    getEnvAsFloat64("MIN_PURCHASE_AMOUNT", 0.01),
			AutoApproveThreshold: getEnvAsFloat64("AUTO_APPROVE_THRESHOLD", 500.00),
			AutoRejectThreshold:  getEnvAsFloat64("AUTO_REJECT_THRESHOLD", 9500.00),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsFloat64(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
