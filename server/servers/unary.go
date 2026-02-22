package servers

import (
	"context"
	"fmt"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func (s *ExamServiceServer) GetExamResult(ctx context.Context, req *exam.GetExamResultRequest) (*exam.GetExamResultResponse, error) {
	key := fmt.Sprintf("%s_%s", req.StudentId, req.ExamId)
	if val, exists := s.examData[key]; exists {
		return val, nil
	} else {
		return nil, fmt.Errorf("exam result not found for student ID %s and exam ID %s", req.StudentId, req.ExamId)
	}
}
