package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rigsecurity/internal/github_cel/interfaces"
	"time"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := interfaces.NewGitAccessInformationV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &interfaces.AccessInformationRequest{
		GitType:           interfaces.GitType_Github,
		AccessControlType: interfaces.AccessControlType_CEL,
		RepositoriesUrl:   []string{"https://github.com/testalexgreenman/testme"},
		Organizations:     "https://github.com/testalexgreenman",
		ScanText:          `userName == "alex148148717" `,
		GitToken:          "token",
	}

	resp, err := client.GetAccessInformation(ctx, req)
	if err != nil {
		log.Fatalf("could not get access information: %v", err)
	}

	log.Printf("Access: %v", resp)
}
