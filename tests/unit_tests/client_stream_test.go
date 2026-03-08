package unittests

import (
	"io"
	"math"
	"testing"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
	"github.com/distroaryan/grpc-exam-service/server/servers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type fakeClientStream struct {
	grpc.ServerStream
	requests []*exam.SubmitExamResultsRequest
	index    int
	response *exam.SubmitExamResultsResponse
}

func (f *fakeClientStream) Recv() (*exam.SubmitExamResultsRequest, error) {
	if f.index >= len(f.requests) {
		return nil, io.EOF
	}
	req := f.requests[f.index]
	f.index++
	return req, nil
}

func (f *fakeClientStream) SendAndClose(resp *exam.SubmitExamResultsResponse) error {
	f.response = resp
	return nil
}

func TestSubmitExamResults(t *testing.T) {
	s := servers.NewExamServiceServer()

	t.Run("submit_multiple_results", func(t *testing.T) {
		stream := &fakeClientStream{
			requests: []*exam.SubmitExamResultsRequest{
				{StudentId: "321", ExamId: "Math", MarksObtained: 80, TotalMarks: 100},
				{StudentId: "321", ExamId: "Physics", MarksObtained: 90, TotalMarks: 100},
				{StudentId: "321", ExamId: "Biology", MarksObtained: 85, TotalMarks: 100},
			},
		}
	
		err := s.SubmitExamResults(stream)
		assert.NoError(t, err)
		assert.NotNil(t, stream.response)
		assert.Equal(t, "321", stream.response.StudentId)
		assert.Equal(t, int32(3), stream.response.TotalExams)
		assert.Equal(t, int32(255), stream.response.TotalMarksObtained)
		assert.Equal(t, int32(300), stream.response.TotalPossibleMarks)
		assert.InDelta(t, 85.0, stream.response.Percentage, 0.01)
	})

	t.Run("no_results_submitted", func(t *testing.T) {
		stream := &fakeClientStream{requests: []*exam.SubmitExamResultsRequest{}}
		err := s.SubmitExamResults(stream)
		assert.NoError(t, err)
		assert.NotNil(t, stream.response)
		assert.Equal(t, int32(0), stream.response.TotalExams)
		assert.True(t, math.IsNaN(float64(stream.response.Percentage)), "Expected NaN when no exam results are submitted")
	})
}
