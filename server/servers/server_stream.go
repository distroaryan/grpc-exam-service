package servers

import (
	"fmt"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func (s *ExamServiceServer) StreamExamResults(req *exam.StreamExamResultsRequest, stream exam.ExamService_StreamExamResultsServer) error {
	studentId := req.StudentId
	examIds := req.ExamIds

	found := false

	for _, examId := range examIds {
		key := fmt.Sprintf("%s_%s", studentId, examId)
		if val, isExist := s.examData[key]; isExist {
			stream.Send(val)
			found = true
		}

		if !found {
			return fmt.Errorf("exam results not found for student ID %s and exam IDs %v", studentId, examIds)
		}
	}
	return nil
}
