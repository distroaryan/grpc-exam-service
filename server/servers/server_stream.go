package servers

import (
	"fmt"
	"time"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func (s *ExamServiceServer) StreamExamResults(req *exam.StreamExamResultsRequest, stream exam.ExamService_StreamExamResultsServer) error {
	studentId := req.StudentId
	examIds := req.ExamIds


	for _, examId := range examIds {
		key := fmt.Sprintf("%s_%s", studentId, examId)
		if val, isExist := s.examData[key]; isExist {
			stream.Send(val)
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
