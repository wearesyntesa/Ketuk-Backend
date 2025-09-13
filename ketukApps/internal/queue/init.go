package queue

import (
	"fmt"
	"ketukApps/config"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection

func InitRabbitMQ(config *config.Config)(err error) {
	url := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(config.Queue.User, config.Queue.Password),
		Host:   config.Queue.Host + ":" + config.Queue.Port,
	}
	RabbitMQConn, err = amqp.Dial(url.String())
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer RabbitMQConn.Close()

	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		config.Queue.Name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	return nil
}

func GetRabbitMQConnection() *amqp.Connection {
	return RabbitMQConn
}

func CloseRabbitMQ() {
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
