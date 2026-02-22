package clients

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func Server_stream(client exam.ExamServiceClient) {
	req := &exam.StreamExamResultsRequest{
		StudentId: "123",
		ExamIds:   []string{"math101", "phy101", "hist101"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.StreamExamResults(ctx, req)
	if err != nil {
		log.Fatalf("error calling StreamExamResults: %v", err)
	}

	fmt.Println("Streaming exam results:")

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("All results received")
				break
			}
			log.Fatalf("error receiving exam result: %v", err)
		}
		fmt.Printf("- %s: %s (%d/%d), Grade: %s\n",
			res.StudentName, res.SubjectName, res.MarkObtained, res.TotalMarks, res.Grade)
	}
}
