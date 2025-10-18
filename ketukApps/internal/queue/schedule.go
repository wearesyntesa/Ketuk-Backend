package queue

import (
	"encoding/json"
	"fmt"
	"ketukApps/internal/models"
	"ketukApps/internal/services"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SchduleWorker(name string, ticketService *services.TicketService, scheduleService *services.ScheduleService) error {
	for {
		msgs, err := ConsumerSchedule(name)
		if err != nil {
			log.Printf("Failed to start consumer: %s", err)
			return err
		}

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Process the message here
			// Parse message body to get both schedule and ticket data
			requestData, err := parseBodyToJSON(d.Body)
			if err != nil {
				log.Printf("Failed to parse message body: %s", err)
				log.Printf("Raw message body (hex): %x", d.Body)
				log.Printf("Raw message body (string): %q", string(d.Body))
				d.Nack(false, false)
				continue
			}
			log.Printf("Parsed RequestData: %+v", requestData)

			// Step 1: Create schedule_ticket first
			scheduleTicket := &models.ScheduleTicket{
				Title:       requestData.Title,
				StartDate:   requestData.StartDate,
				EndDate:     requestData.EndDate,
				UserID:      int(requestData.UserID),
				Kategori:    requestData.Category,
				Description: requestData.Description,
			}

			savedSchedule, err := scheduleService.CreateScheduleTicket(scheduleTicket)
			if err != nil {
				log.Printf("Failed to save schedule_ticket to database: %s", err)
				d.Nack(false, false)
				continue
			}
			log.Printf("Successfully saved schedule_ticket to database with ID: %d", savedSchedule.IDSchedule)

			// Step 2: Create ticket with the schedule ID as foreign key
			ticket := &models.Ticket{
				UserID:      requestData.UserID,
				Title:       requestData.Title,
				Description: requestData.Description,
				Status:      requestData.Status,
				IDSchedule:  &savedSchedule.IDSchedule,
			}

			savedTicket, err := ticketService.CreateFromModel(ticket)
			if err != nil {
				log.Printf("Failed to save ticket to database: %s", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("Successfully saved ticket to database with ID: %d, linked to schedule ID: %d", savedTicket.ID, *savedTicket.IDSchedule)

			// Acknowledge the message after successful processing
			d.Ack(false)
		}
	}
}

func ConsumerSchedule(name string) (<-chan amqp.Delivery, error) {
	q, err := RabbitMQClient.Channel.QueueDeclare(
		name,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := RabbitMQClient.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
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

func WorkerSchedule() {
	log.Println("Worker started, waiting for messages...")

}

// ScheduleTicketMessage represents the message format from RabbitMQ
// It contains all the data needed to create both schedule_ticket and ticket
type ScheduleTicketMessage struct {
	UserID      uint            `json:"userId"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Category    models.Category `json:"category"`
	StartDate   time.Time       `json:"startDate"`
	EndDate     time.Time       `json:"endDate"`
}

func parseBodyToJSON(body []byte) (*ScheduleTicketMessage, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("empty message body")
	}
	var message ScheduleTicketMessage
	err := json.Unmarshal(body, &message)
	if err != nil {
		// Try to identify the problematic character
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			problemChar := ""
			if int(jsonErr.Offset) < len(body) {
				problemChar = string(body[jsonErr.Offset])
			}
			return nil, fmt.Errorf("JSON syntax error at offset %d (char: %q): %w", jsonErr.Offset, problemChar, err)
		}
		return nil, fmt.Errorf("failed to parse message body: %w", err)
	}
	return &message, nil
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
