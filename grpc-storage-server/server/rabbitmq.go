package server

import (
	"fmt"
	"log"
	pb "paper-storage-server/paper-storage-server/paperpb"

	"github.com/streadway/amqp"
)

func (s *Server) notifyNewPaper(paper *pb.Paper) error {
	if s.RabbitConn == nil {
		return fmt.Errorf("RabbitMQ connection is nil")
	}

	ch, err := s.RabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	queueName := "papers"
	_, err = ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable (keeps queue even if RabbitMQ restarts)
		false,     // Auto-delete (deletes when no consumers)
		false,     // Exclusive (only used by this connection)
		false,     // No-wait (doesn't wait for server response)
		nil,       // Additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.Publish(
		"",       // Exchange
		"papers", // Queue name
		false,    // Mandatory
		false,    // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("New paper of type %v added to the paper store: %s by %s", paper.Format, paper.Title, paper.Author)),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf("Sent message: New paper added: %s by %s", paper.Title, paper.Author)
	return nil
}

// RegisterClient adds a client to the server's notification list
func (s *Server) RegisterClient(clientID string) chan string {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	// Create a notification channel for the client
	clientChan := make(chan string)
	s.Clients[clientID] = clientChan
	log.Printf("[Client Registered] Client ID: %s", clientID)
	return clientChan
}

// BroadcastNotification sends a notification to all connected clients
func (s *Server) BroadcastNotification(message string) {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	for clientID, clientChan := range s.Clients {
		// Send the message in a non-blocking way

		go func(clientID string, clientChan chan string) {
			select {
			case clientChan <- message:
				log.Printf("[Notification Sent] To: %s, Message: %s", clientID, message)
			default:
				log.Printf("[Notification Skipped] Client ID: %s (Channel Full)", clientID)
			}
		}(clientID, clientChan)

	}
}
