package stresstest

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// init() runs once when the test package is loaded. We use it to spin up an in-memory gRPC server.
func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	examServer := servers.NewExamServiceServer()
	exam.RegisterExamServiceServer(s, examServer)

	// Start serving traffic in the background
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

// bufDialer connects to the in-memory listener instead of actual network hardware
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// mustNoError is a helper to assert no error occurred during a benchmark.
// It fails the benchmark immediately with a descriptive message if err != nil.
// op describes the operation that failed (e.g. "Dial", "Send", "CloseAndRecv").
func mustNoError(b *testing.B, err error, op string) {
	b.Helper()
	if err != nil {
		b.Fatalf("[STRESS FAIL] operation=%q error=%s", op, fmt.Sprintf("%v", err))
	}
}

// BenchmarkGetExamResult is a stress test for the Unary endpoint.
func BenchmarkGetExamResult(b *testing.B) {
	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	mustNoError(b, err, "Dial bufnet")
	defer conn.Close()

	client := exam.NewExamServiceClient(conn)
	req := &exam.GetExamResultRequest{
		StudentId: "123",
		ExamId:    "math101",
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetExamResult(ctx, req)
			mustNoError(b, err, "GetExamResult RPC")
		}
	})
}

// BenchmarkStreamExamResults is a stress test for the Server Streaming endpoint.
func BenchmarkStreamExamResults(b *testing.B) {
	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	mustNoError(b, err, "Dial bufnet")
	defer conn.Close()

	client := exam.NewExamServiceClient(conn)
	req := &exam.StreamExamResultsRequest{
		StudentId: "123",
		ExamIds:   []string{"math101", "phy101", "hist101"},
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stream, err := client.StreamExamResults(ctx, req)
			mustNoError(b, err, "StreamExamResults open stream")

			for {
				_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				mustNoError(b, err, "StreamExamResults Recv")
			}
		}
	})
}

// BenchmarkSubmitExamResults is a stress test for the Client Streaming endpoint.
func BenchmarkSubmitExamResults(b *testing.B) {
	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	mustNoError(b, err, "Dial bufnet")
	defer conn.Close()

	client := exam.NewExamServiceClient(conn)

	requests := []*exam.SubmitExamResultsRequest{
		{StudentId: "321", ExamId: "Math", MarksObtained: 80, TotalMarks: 100},
		{StudentId: "321", ExamId: "Physics", MarksObtained: 90, TotalMarks: 100},
		{StudentId: "321", ExamId: "Biology", MarksObtained: 85, TotalMarks: 100},
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stream, err := client.SubmitExamResults(ctx)
			mustNoError(b, err, "SubmitExamResults open stream")

			for _, r := range requests {
				mustNoError(b, stream.Send(r), "SubmitExamResults Send")
			}

			_, err = stream.CloseAndRecv()
			mustNoError(b, err, "SubmitExamResults CloseAndRecv")
		}
	})
}

// BenchmarkLiveExamQuery is a stress test for the Bidirectional Streaming endpoint.
func BenchmarkLiveExamQuery(b *testing.B) {
	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	mustNoError(b, err, "Dial bufnet")
	defer conn.Close()

	client := exam.NewExamServiceClient(conn)

	requests := []*exam.GetExamResultRequest{
		{StudentId: "123", ExamId: "math101"},
		{StudentId: "123", ExamId: "phy101"},
		{StudentId: "123", ExamId: "hist101"},
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		stream, err := client.LiveExamQuery(ctx)
		mustNoError(b, err, "LiveExamQuery open stream")

		for pb.Next() {
			for _, r := range requests {
				mustNoError(b, stream.Send(r), "LiveExamQuery Send")

				_, err = stream.Recv()
				mustNoError(b, err, "LiveExamQuery Recv")
			}
		}
	})
}
