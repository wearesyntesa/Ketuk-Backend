package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Host     string
	LogLevel string
	Database DatabaseConfig
	Queue    QueueConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type QueueConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	VHost             string
	MaxChannels       int
	HeartbeatInterval int
	ConnectionTimeout int
	Prefetch          int
	Exchanges         ExchangesConfig
	Queues            QueuesConfig
}

type ExchangesConfig struct {
	Direct     string
	Topic      string
	Fanout     string
	DeadLetter string
}

type QueuesConfig struct {
	TicketsCreated string
	TicketsUpdated string
	TicketsDeleted string
	UsersCreated   string
	UsersUpdated   string
	Notifications  string
	Emails         string
	DeadLetters    string
}



func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:     getEnv("PORT", "8080"),
		Host:     getEnv("HOST", "localhost"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "mydb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Queue: QueueConfig{
			Host:              getEnv("QUEUE_HOST", "localhost"),
			Port:              getEnv("QUEUE_PORT", "5672"),
			User:              getEnv("QUEUE_USER", "user"),
			Password:          getEnv("QUEUE_PASSWORD", "password"),
			Name:              getEnv("QUEUE_NAME", "schedule"),
			VHost:             getEnv("QUEUE_VHOST", "/"),
			MaxChannels:       getEnvInt("QUEUE_MAX_CHANNELS", 20),
			HeartbeatInterval: getEnvInt("QUEUE_HEARTBEAT", 60),
			ConnectionTimeout: getEnvInt("QUEUE_CONNECT_TIMEOUT", 30),
			Prefetch:          getEnvInt("QUEUE_PREFETCH", 1),
			Exchanges: ExchangesConfig{
				Direct:     getEnv("QUEUE_EXCHANGE_DIRECT", "ketuk.direct"),
				Topic:      getEnv("QUEUE_EXCHANGE_TOPIC", "ketuk.topic"),
				Fanout:     getEnv("QUEUE_EXCHANGE_FANOUT", "ketuk.fanout"),
				DeadLetter: getEnv("QUEUE_EXCHANGE_DLX", "ketuk.dlx"),
			},
			Queues: QueuesConfig{
				TicketsCreated: getEnv("QUEUE_TICKETS_CREATED", "tickets.created"),
				TicketsUpdated: getEnv("QUEUE_TICKETS_UPDATED", "tickets.updated"),
				TicketsDeleted: getEnv("QUEUE_TICKETS_DELETED", "tickets.deleted"),
				UsersCreated:   getEnv("QUEUE_USERS_CREATED", "users.created"),
				UsersUpdated:   getEnv("QUEUE_USERS_UPDATED", "users.updated"),
				Notifications:  getEnv("QUEUE_NOTIFICATIONS", "notifications"),
				Emails:         getEnv("QUEUE_EMAILS", "emails"),
				DeadLetters:    getEnv("QUEUE_DEAD_LETTERS", "dead.letters"),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
