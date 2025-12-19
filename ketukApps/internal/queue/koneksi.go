package queue

import (
	"ketukApps/config"
	"log"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQClient *RabbitMQ

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQConnection(cfg *config.Config) (err error) {
	url := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(cfg.Queue.User, cfg.Queue.Password),
		Host:   cfg.Queue.Host + ":" + cfg.Queue.Port,
	}

	amqpConfig := amqp.Config{
		Heartbeat: 60 * time.Second,
		Locale:    "id_ID",
		
	}

	conn, err := amqp.DialConfig(url.String(), amqpConfig)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a RabbitMQ channel: %s", err)
	}

	RabbitMQClient = &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}
	log.Println("Connected to RabbitMQ")
	return nil
}
