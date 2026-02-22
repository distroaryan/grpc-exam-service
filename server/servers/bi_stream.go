package servers

import (
	"fmt"
	"io"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func (s *ExamServiceServer) LiveExamQuery(stream exam.ExamService_LiveExamQueryServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		key := fmt.Sprintf("%s_%s", req.StudentId, req.ExamId)
		result, ok := s.examData[key]
		if !ok {
			res := &exam.GetExamResultResponse{
				StudentName:  "N/A",
				SubjectName:  req.ExamId,
				MarkObtained: 0,
				TotalMarks:   0,
				Grade:        "Not Found",
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			continue
		}
		if err := stream.Send(result); err != nil {
			return err
		}
	}
}
