package unittests

import (
	"io"
	"testing"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type fakeBidirectionalStream struct {
	grpc.ServerStream
	requests  []*exam.GetExamResultRequest
	index     int
	responses []*exam.GetExamResultResponse
}

func (f *fakeBidirectionalStream) Recv() (*exam.GetExamResultRequest, error) {
	if f.index >= len(f.requests) {
		return nil, io.EOF
	}
	resp := f.requests[f.index]
	f.index++
	return resp, nil
}

func (f *fakeBidirectionalStream) Send(resp *exam.GetExamResultResponse) error {
	f.responses = append(f.responses, resp)
	return nil
}

func TestLiveExamQuery(t *testing.T) {
	t.Run("correct examIds and studentIds", func(t *testing.T) {
		s := servers.NewExamServiceServer()
		stream := &fakeBidirectionalStream{
			requests: []*exam.GetExamResultRequest{
				{
					StudentId: "123",
					ExamId:    "math101",
				},
				{
					StudentId: "123",
					ExamId:    "phy101",
				},
				{
					StudentId: "123",
					ExamId:    "hist101",
				},
			},
			index: 0,
		}

		expectedResponse := []*exam.GetExamResultResponse{
			{
				StudentName:  "John Doe",
				SubjectName:  "Math 101",
				MarkObtained: 95,
				TotalMarks:   100,
				Grade:        "A+",
			},
			{
				StudentName:  "John Doe",
				SubjectName:  "Physics 101",
				MarkObtained: 81,
				TotalMarks:   100,
				Grade:        "A",
			},
			{
				StudentName:  "John Doe",
				SubjectName:  "History 101",
				MarkObtained: 76,
				TotalMarks:   100,
				Grade:        "B",
			},
		}

		err := s.LiveExamQuery(stream)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, stream.responses)
	})

	t.Run("incorrect examIds", func(t *testing.T) {
		s := servers.NewExamServiceServer()
		stream := &fakeBidirectionalStream{
			requests: []*exam.GetExamResultRequest{
				{
					StudentId: "123",
					ExamId:    "math101",
				},
				{
					StudentId: "123",
					ExamId:    "phy101",
				},
				{
					StudentId: "123",
					ExamId:    "bio101",
				},
			},
			index: 0,
		}

		expectedResponse := []*exam.GetExamResultResponse{
			{
				StudentName:  "John Doe",
				SubjectName:  "Math 101",
				MarkObtained: 95,
				TotalMarks:   100,
				Grade:        "A+",
			},
			{
				StudentName:  "John Doe",
				SubjectName:  "Physics 101",
				MarkObtained: 81,
				TotalMarks:   100,
				Grade:        "A",
			},
			{
				StudentName:  "N/A",
				SubjectName:  "bio101",
				MarkObtained: 0,
				TotalMarks:   0,
				Grade:        "Not Found",
			},
		}

		err := s.LiveExamQuery(stream)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, stream.responses)
	})
}
