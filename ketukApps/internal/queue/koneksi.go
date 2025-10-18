package queue

import (
	"ketukApps/config"
	"log"
	"net/url"

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

	conn, err := amqp.Dial(url.String())
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
