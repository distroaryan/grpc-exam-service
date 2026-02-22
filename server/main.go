package main

import (
	"log/slog"
	"net"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"github.com/distroaryan/grpc-exam-service/utils"
	"google.golang.org/grpc"
)

func main() {
	utils.InitLogger(true)

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		slog.Error("failed to lister", "error", err)
	}

	// New gRPC server instance
	s := grpc.NewServer()

	// register the exam service
	exam.RegisterExamServiceServer(s, servers.NewExamServiceServer())

	if err := s.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}
