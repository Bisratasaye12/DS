package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	clientNotification "paper-storage-server/cmd/paperclient/clientNotification"
	pb "paper-storage-server/paper-storage-server/paperpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("\n CLIENT \n")
		fmt.Println("Usage:")
		fmt.Println("  paperclient <server-address>")
		os.Exit(1)
	}

	serverAddress := os.Args[1]

	// Connect to the server
	client, _, err := createClient(serverAddress)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	// Generate a unique client ID
	clientID := uuid.New().String()
	notificationChan := make(chan string)

	// Make the client a consumer for notifications
	clientNotification.ConsumeNotifications(clientID, notificationChan)
	NotificationHandler(clientID, notificationChan)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\nEnter commands (type 'exit' to quit): \n")
	fmt.Println("  [add command]: add <author-name> <title> <filepath>")
	fmt.Println("  [list command]: list")
	fmt.Println("  [fetch command]: fetch <paper-id>")
	fmt.Println("  [detail command]: detail <paper-id>")
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		args := strings.Fields(input)
		if len(args) < 1 {
			fmt.Println("Invalid command. Try again.")
			continue
		}

		command := strings.ToLower(args[0])

		switch command {
		case "add":
			if len(args) != 4 {
				fmt.Println("Usage: add 'Author Name' 'Paper Title' paper.pdf")
				continue
			}
			author := args[1]
			title := args[2]
			filePath := args[3]

			err := addPaper(client, author, title, filePath)
			if err != nil {
				log.Printf("Error adding paper: %v", err)
			}
		case "list":
			err := listPapers(client)
			if err != nil {
				log.Printf("Error listing papers: %v", err)
			}
		case "detail":
			if len(args) != 2 {
				fmt.Println("Usage: detail <paper_id>")
				continue
			}

			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("Invalid paper ID: %v\n", err)
				continue
			}
			paperID := int32(id)

			err = getPaperDetails(client, paperID)
			if err != nil {
				log.Printf("Error fetching paper details: %v", err)
			}
		case "fetch":
			if len(args) != 2 {
				fmt.Println("Usage: fetch <paper_id>")
				continue
			}

			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("Invalid paper ID: %v\n", err)
				continue
			}
			paperID := int32(id)

			err = fetchPaperContent(client, paperID)
			if err != nil {
				log.Printf("Error fetching paper content: %v", err)
			}
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}

// createClient initializes a gRPC client using grpc.Dial
func createClient(serverAddress string) (pb.PaperStorageServiceClient, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		if state == connectivity.TransientFailure {
			return nil, nil, fmt.Errorf("gRPC connection is in transient failure state")
		}
		if !conn.WaitForStateChange(ctx, state) {
			return nil, nil, fmt.Errorf("timeout waiting for gRPC connection to become ready")
		}
	}

	client := pb.NewPaperStorageServiceClient(conn)
	return client, conn, nil
}

func addPaper(client pb.PaperStorageServiceClient, author, title, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	tmp := strings.Split(filePath, ".")
	format := strings.ToUpper(strings.TrimSpace(tmp[len(tmp) - 1]))
	fmt.Println("File path and format", filePath, format)
	req := &pb.AddPaperArgs{
		Paper: &pb.Paper{
			Author:  author,
			Title:   title,
			Format:  format,
			Content: content,
		},
	}

	res, err := client.AddPaper(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to add paper: %v", err)
	}

	fmt.Println(res.Result)
	return nil
}

func NotificationHandler(clientID string, notificationChan chan string) {
	go func() {
		for message := range notificationChan {
			fmt.Printf("[Notification] Client ID: %s, Message: %s\n", clientID, message)
		}
	}()
}

func listPapers(client pb.PaperStorageServiceClient) error {
	req := &pb.ListPapersArgs{}

	res, err := client.ListPapers(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to list papers: %v", err)
	}

	fmt.Println("Stored Papers:")
	for _, paper := range res.Papers {
		fmt.Printf("Paper ID: %v | Author: %v | Title: %v\n", paper.PaperId, paper.Author, paper.Title)
	}

	return nil
}

func getPaperDetails(client pb.PaperStorageServiceClient, paperID int32) error {
	req := &pb.GetPaperDetailsArgs{
		PaperId: paperID,
	}

	res, err := client.GetPaperDetails(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get paper details: %v", err)
	}

	fmt.Printf("Paper ID: %v\n", paperID)
	fmt.Printf("Author: %s\n", res.Paper.Author)
	fmt.Printf("Title: %s\n", res.Paper.Title)

	return nil
}

func fetchPaperContent(client pb.PaperStorageServiceClient, paperID int32) error {
	req := &pb.FetchPaperContentArgs{
		PaperId: paperID,
	}

	res, err := client.FetchPaperContent(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to fetch paper content: %v", err)
	}

	fmt.Printf("Paper Content for ID %v:\n%s\n", paperID, string(res.Content))
	return nil
}
