package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "testgrpc"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Kênh để nhận tin nhắn
	messageChan := make(chan string)
	go func() {
		for msg := range messageChan {
			fmt.Println(msg)
		}
	}()

	// Thêm client vào server
	defer func() {
		close(messageChan)
	}()

	// Gọi hàm nhận tin nhắn từ server
	go receiveMessage(c, messageChan)

	for {
		var name string
		fmt.Print("Client 2 - Enter your name (or type 'exit' to quit): ")
		fmt.Scanln(&name)

		if name == "exit" {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Client 2 - Greeting: %s", resp.Message)
		// Gửi tin nhắn
		if err := sendMessage(c, name); err != nil {
			log.Fatalf("could not send message: %v", err)
		}
	}
}

func sendMessage(c pb.GreeterClient, from string) error {
	var message string
	fmt.Print("Client 2 - Enter your message: ")
	fmt.Scanln(&message)

	_, err := c.SendMessage(context.Background(), &pb.MessageRequest{From: from, Message: message})
	return err
}

func receiveMessage(c pb.GreeterClient, ch chan string) {
	stream, err := c.ReceiveMessage(context.Background() /* context */, &pb.MessageRequest{})
	if err != nil {
		log.Fatalf("could not receive message: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("could not receive message: %v", err)
		}
		ch <- resp.Message
	}
}
