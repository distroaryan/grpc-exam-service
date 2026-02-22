package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/distroaryan/grpc-exam-service/proto/generate/exam"
)

func Unary(client exam.ExamServiceClient) {
	fmt.Println("Enter student ID and exam ID (e.g. 123, math101)")
	var studentId, examId string
	fmt.Scanf("%s %s", &studentId, &examId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &exam.GetExamResultRequest{
		StudentId: studentId,
		ExamId:    examId,
	}
	resp, err := client.GetExamResult(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Student Name: %s\n", resp.StudentName)
	fmt.Printf("Subject: %s\n", resp.SubjectName)
	fmt.Printf("Marks Obtained: %d out of %d\n", resp.MarkObtained, resp.TotalMarks)
	fmt.Printf("Grade: %s\n", resp.Grade)
	fmt.Println("Unary RPC call completed successfully.")
}
