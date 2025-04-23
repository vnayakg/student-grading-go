package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) ([]student, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	//skip csv header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %v", err)
	}

	var students []student
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading record: %v", err)
		}

		student, err := parseStudentRecord(record)
		if err != nil {
			fmt.Printf("skipping invalid record: %v", err)
			continue
		}

		students = append(students, student)
	}

	return students, nil
}

func parseStudentRecord(record []string) (student, error) {
	if len(record) < 7 {
		return student{}, fmt.Errorf("incomplete record")
	}

	test1, err := strconv.Atoi(record[3])
	if err != nil {
		return student{}, fmt.Errorf("invalid Test1 score: %v", err)
	}
	test2, err := strconv.Atoi(record[4])
	if err != nil {
		return student{}, fmt.Errorf("invalid Test2 score: %v", err)
	}
	test3, err := strconv.Atoi(record[5])
	if err != nil {
		return student{}, fmt.Errorf("invalid Test3 score: %v", err)
	}
	test4, err := strconv.Atoi(record[6])
	if err != nil {
		return student{}, fmt.Errorf("invalid Test4 score: %v", err)
	}

	return student{
		firstName:  record[0],
		lastName:   record[1],
		university: record[2],
		test1Score: test1,
		test2Score: test2,
		test3Score: test3,
		test4Score: test4,
	}, nil
}

func calculateGrade(students []student) []studentStat {
	var studentSet []studentStat

	for _, student := range students {
		finalScore := student.getFinalScore()
		grade := student.getGrade(finalScore)

		studentSet = append(studentSet, studentStat{student, finalScore, grade})
	}

	return studentSet
}

func (s *student) getFinalScore() float32 {
	return float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4
}

func (s *student) getGrade(finalScore float32) Grade {
	switch {
	case finalScore >= 70:
		return A
	case finalScore >= 50 && finalScore < 70:
		return B
	case finalScore >= 35 && finalScore < 50:
		return C
	default:
		return F
	}
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var topperScore float32 = 0
	var topperStudent studentStat

	for _, gradedStudent := range gradedStudents {
		if gradedStudent.finalScore > topperScore {
			topperScore = gradedStudent.finalScore
			topperStudent = gradedStudent
		}
	}
	return topperStudent
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	return nil
}
