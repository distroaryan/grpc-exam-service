package unittests

import (
	"testing"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type fakeServerStream struct {
	grpc.ServerStream
	response []*exam.GetExamResultResponse
}

func (f *fakeServerStream) Send(resp *exam.GetExamResultResponse) error {
	f.response = append(f.response, resp)
	return nil
}

func TestStreamExamResults(t *testing.T) {
	t.Run("correct examIds and studentIds", func(t *testing.T) {
		req := &exam.StreamExamResultsRequest{
			StudentId: "123",
			ExamIds:   []string{"math101", "phy101", "hist101"},
		}

		s := servers.NewExamServiceServer()
		stream := &fakeServerStream{}

		err := s.StreamExamResults(req, stream)
		assert.NoError(t, err)
		assert.Len(t, stream.response, 3)

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

		assert.Equal(t, expectedResponse, stream.response)
	})

	t.Run("incorrect examIds", func(t *testing.T) {
		req := &exam.StreamExamResultsRequest{
			StudentId: "123",
			ExamIds: []string{"math101", "phy101", "bio101"},
		}
		s := servers.NewExamServiceServer()
		stream := &fakeServerStream{}

		err := s.StreamExamResults(req, stream)
		assert.NoError(t, err)
		assert.Len(t, stream.response, 2)

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
		}

		assert.Equal(t, expectedResponse, stream.response)
		
	})
}
