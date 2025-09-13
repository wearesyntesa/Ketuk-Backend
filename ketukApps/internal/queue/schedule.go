package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SchduleWorker(name string) error {
	for {
		msgs, err := ConsumerSchedule(name)
		if err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}

		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
			// Process the message here
			// Acknowledge the message if not using auto-ack
			// d.Ack(false)
		}
	}
	return nil
}



func ConsumerSchedule(name string) (<-chan amqp.Delivery, error) {
	q, err := RabbitMQClient.Channel.QueueDeclare(
		name,
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := RabbitMQClient.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return msgs, nil
}

func CloseRabbitMQ() {
	if RabbitMQClient != nil {
		if RabbitMQClient.Channel != nil {
			RabbitMQClient.Channel.Close()
		}
		if RabbitMQClient.Conn != nil {
			RabbitMQClient.Conn.Close()
		}
	}
}