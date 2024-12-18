package clientnotification

import (
	"log"

	"github.com/streadway/amqp"
)

func ConsumeNotifications(clientID string, notificationChan chan string) {
	go func() {
		// Connect to RabbitMQ
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatalf("[Error] Failed to connect to RabbitMQ: %v", err)
		}
		defer conn.Close()

		// Create a channel
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("[Error] Failed to open a channel: %v", err)
		}
		defer ch.Close()

		// Declare the queue
		queueName := "papers"
		_, err = ch.QueueDeclare(
			queueName, // Name of the queue
			true,      // Durable
			false,     // Auto-delete
			false,     // Exclusive
			false,     // No-wait
			nil,       // Additional arguments
		)
		if err != nil {
			log.Fatalf("[Error] Failed to declare the queue: %v", err)
		}

		// Register as a consumer
		msgs, err := ch.Consume(
			queueName, // Name of the queue
			"",        // Consumer tag
			true,      // Auto-acknowledge
			false,     // Exclusive
			false,     // No-local
			false,     // No-wait
			nil,       // Additional arguments
		)
		if err != nil {
			log.Fatalf("[Error] Failed to register as a consumer: %v", err)
		}

		log.Printf("[Info] Listening for notifications for Client ID: %s", clientID)

		// Process messages
		for msg := range msgs {
			notificationChan <- string(msg.Body)
		}
	}()
}