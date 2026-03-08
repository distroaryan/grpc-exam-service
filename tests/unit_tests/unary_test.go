package unittests

import (
	"context"
	"testing"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"github.com/stretchr/testify/assert"
)

func TestGetExamResult(t *testing.T) {
	s := servers.NewExamServiceServer()
	tests := []struct {
		name string
		request *exam.GetExamResultRequest
		response *exam.GetExamResultResponse
	}{
		{
			name: "valid student id and exam id",
			request: &exam.GetExamResultRequest{
				StudentId: "123",
				ExamId: "math101",
			},
			response: &exam.GetExamResultResponse{
				StudentName: "John Doe",
				SubjectName: "Math 101",
				MarkObtained: 95,
				TotalMarks: 100,
				Grade: "A+",
			},
		},
		{
			name: "invalid student id and exam id",
			request: &exam.GetExamResultRequest{
				StudentId: "qwerty",
				ExamId: "masqom",
			},
			response: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.GetExamResult(context.Background(), tt.request)

			if tt.response == nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.response, resp)
			}
		})
	}
}