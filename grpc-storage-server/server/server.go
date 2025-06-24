package server

import (
	"context"
	"fmt"
	pb "paper-storage-server/paper-storage-server/paperpb"
	"sync"

	"github.com/streadway/amqp"
)

type Server struct {
	pb.UnimplementedPaperStorageServiceServer
	ContentStore map[int32][]byte
	DetailStore  map[int32]map[string]string
	ContentMutex sync.Mutex
	DetailMutex  sync.Mutex
	ClientMutex  sync.Mutex
	Clients      map[string]chan string
	RabbitConn   *amqp.Connection
	ID           int
}

func (s *Server) AddPaper(ctx context.Context, req *pb.AddPaperArgs) (*pb.AddPaperResponse, error) {
	paper := req.Paper
	paper.PaperId = int32(s.ID)
	s.ID += 1

	s.ContentMutex.Lock()
	s.ContentStore[paper.PaperId] = paper.Content
	s.ContentMutex.Unlock()

	s.DetailMutex.Lock()
	s.DetailStore[paper.PaperId] = map[string]string{"title": paper.Title, "author": paper.Author, "format": paper.Format}
	s.DetailMutex.Unlock()

	// Notify via RabbitMQ about the new paper added
	err := s.notifyNewPaper(paper)
	if err != nil {
		return nil, fmt.Errorf("failed to send RabbitMQ message: %v", err)
	}

	// Broadcast notification to all clients
	message := fmt.Sprintf("New paper added: ID=%d, Title='%s'", paper.PaperId, paper.Title)
	s.BroadcastNotification(message)

	var response pb.AddPaperResponse
	response.Result = "paper content stored successfully!"

	return &response, nil
}

func (s *Server) FetchPaperContent(ctx context.Context, req *pb.FetchPaperContentArgs) (*pb.FetchPaperContentResponse, error) {
	var response *pb.FetchPaperContentResponse
	s.ContentMutex.Lock()
	defer s.ContentMutex.Unlock()
	if paperContent, ok := s.ContentStore[req.PaperId]; ok {
		response = &pb.FetchPaperContentResponse{
			Content: paperContent,
		}
		return response, nil
	}

	return response, fmt.Errorf("paper not found")
}

func (s *Server) GetPaperDetails(ctx context.Context, req *pb.GetPaperDetailsArgs) (*pb.GetPaperDetailsResponse, error) {
	var response pb.GetPaperDetailsResponse
	s.DetailMutex.Lock()
	defer s.DetailMutex.Unlock()
	if paperDetails, ok := s.DetailStore[req.PaperId]; ok {
		// Initialize the Paper field before assigning values
		response.Paper = &pb.Paper{
			Title:  paperDetails["title"],
			Author: paperDetails["author"],
			Format: paperDetails["format"],
		}
		return &response, nil
	}

	return &response, fmt.Errorf("paper not found")
}

func (s *Server) ListPapers(ctx context.Context, req *pb.ListPapersArgs) (*pb.ListPapersResponse, error) {
	var response pb.ListPapersResponse

	for paperID := range s.DetailStore {
		paper := pb.Paper{
			PaperId: paperID,
			Author:  s.DetailStore[paperID]["author"],
			Format:  s.DetailStore[paperID]["format"],
			Title:   s.DetailStore[paperID]["title"],
			Content: s.ContentStore[paperID],
		}
		response.Papers = append(response.Papers, &paper)
	}

	return &response, nil
}
