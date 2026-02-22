package servers

import "github.com/distroaryan/grpc-exam-service/proto/generate/exam"

type ExamServiceServer struct {
	exam.UnimplementedExamServiceServer
	examData map[string]*exam.GetExamResultResponse
}

func NewExamServiceServer() *ExamServiceServer {
	data := map[string]*exam.GetExamResultResponse{
		"123_math101": {
			StudentName:  "John Doe",
			SubjectName:  "Math 101",
			MarkObtained: 95,
			TotalMarks:   100,
			Grade:        "A+",
		},
		"123_phy101": {
			StudentName:  "John Doe",
			SubjectName:  "Physics 101",
			MarkObtained: 81,
			TotalMarks:   100,
			Grade:        "A",
		},
		"123_hist101": {
			StudentName:  "John Doe",
			SubjectName:  "History 101",
			MarkObtained: 76,
			TotalMarks:   100,
			Grade:        "B",
		},
		"456_phy101": {
			StudentName:  "Jane Smith",
			SubjectName:  "Physics 101",
			MarkObtained: 88,
			TotalMarks:   100,
			Grade:        "A",
		},
		"789_chem101": {
			StudentName:  "Alice Johnson",
			SubjectName:  "Chemistry 101",
			MarkObtained: 92,
			TotalMarks:   100,
			Grade:        "A+",
		},
		"101_bio101": {
			StudentName:  "Bob Brown",
			SubjectName:  "Biology 101",
			MarkObtained: 85,
			TotalMarks:   100,
			Grade:        "A",
		},
		"102_hist101": {
			StudentName:  "Charlie Davis",
			SubjectName:  "History 101",
			MarkObtained: 90,
			TotalMarks:   100,
			Grade:        "A+",
		},
	}
	return &ExamServiceServer{
		examData: data,
	}
}
