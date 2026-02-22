package servers

import (
	"io"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func (s *ExamServiceServer) SubmitExamResults(stream exam.ExamService_SubmitExamResultsServer) error {
	var (
		totalExams         int32
		totalMarksObtained int32
		totalPossibleMarks int32
		studentId          string
	)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				percentage := float32(totalMarksObtained) / float32(totalPossibleMarks) * 100
				var grade string
				switch {
				case percentage >= 90.0:
					grade = "A+"
				case percentage >= 85.0:
					grade = "A"
				case percentage >= 80.0:
					grade = "B+"
				case percentage >= 75.0:
					grade = "B"
				case percentage >= 70.0:
					grade = "C+"
				case percentage >= 65.0:
					grade = "C"
				case percentage >= 60.0:
					grade = "D+"
				case percentage >= 55.0:
					grade = "D"
				case percentage >= 50.0:
					grade = "E"
				default:
					grade = "F"
				}
				res := &exam.SubmitExamResultsResponse{
					StudentId:          studentId,
					TotalExams:         totalExams,
					TotalMarksObtained: totalMarksObtained,
					TotalPossibleMarks: totalPossibleMarks,
					Percentage:         percentage,
					Grade: grade,
				}
				return stream.SendAndClose(res)
			}
		}
		studentId = req.StudentId
		totalExams++
		totalMarksObtained += req.MarksObtained
		totalPossibleMarks += req.TotalMarks
	}
}
