package main

import (
	"context"
	"log"
	"net"
	"sync"
	pb "testgrpc"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
	clients     map[pb.Greeter_ReceiveMessageServer]struct{}
	mu          sync.Mutex
	messageChan chan string
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", req.Name)
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	s.messageChan <- req.Message
	return &pb.MessageResponse{Message: "OK"}, nil
}

func (s *server) ReceiveMessage(req *pb.MessageRequest, stream pb.Greeter_ReceiveMessageServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()
	go func() {
		for {
			select {
			case msg := <-s.messageChan:
				err := stream.Send(&pb.MessageResponse{Message: msg})
				if err != nil {
					s.mu.Lock()
					delete(s.clients, stream)
					s.mu.Unlock()
					return
				}
			}
		}
	}()
	// Gửi dữ liệu đến kênh messageChan
	s.messageChan <- "Hello from server!"
	// Gửi dữ liệu đến client định kỳ
	ticker := time.NewTicker(time.Second)
	sendCount := 0
	for range ticker.C {
		if sendCount < 5 {
			s.messageChan <- "Hello from server!"
			sendCount++
		} else {
			ticker.Stop()
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, &server{clients: make(map[pb.Greeter_ReceiveMessageServer]struct{}), messageChan: make(chan string)})

	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
