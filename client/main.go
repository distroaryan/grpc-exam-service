package main

import (
	"log/slog"
	"os"

	"github.com/distroaryan/grpc-exam-service/client/clients"
	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	utils.InitLogger(true)

	if len(os.Args) < 2 {
		slog.Error("Usage: go run client/main.go [unary|server|client|bi]")
		return
	}

	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}
	defer conn.Close()

	client := exam.NewExamServiceClient(conn)

	switch os.Args[1] {
	case "unary":
		clients.Unary(client)
	case "server":
		clients.Server_stream(client)
	case "client":
		clients.Client_Stream(client)
	case "bidirectional":
		clients.BidirectionalStream(client)
	default:
		slog.Error("Unknown command. Use one of: unary, server, client, bi")
	}

}
