package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Student represents a single student record with their ID, name, and grades.
type Student struct {
	ID     int                `json:"id"`
	Name   string             `json:"name"`
	Grades map[string]float64 `json:"grades"`
}

// NewStudent creates and returns a new Student instance with initialized grades map.
func NewStudent(id int, name string) *Student {
	return &Student{
		ID:     id,
		Name:   strings.TrimSpace(name),
		Grades: make(map[string]float64),
	}
}

// AddGrade adds or updates a grade for a specific subject.
// Grade must be between 0.0 and 100.0.
func (s *Student) AddGrade(subject string, grade float64) error {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return errors.New("subject name cannot be empty")
	}
	if grade < 0 || grade > 100 {
		return fmt.Errorf("invalid grade: %.2f (must be between 0 and 100)", grade)
	}
	if s.Grades == nil {
		s.Grades = make(map[string]float64)
	}
	s.Grades[subject] = grade
	return nil
}

// CalculateAverage calculates the average grade for this student.
// Returns 0.0 if the student has no recorded grades.
func (s *Student) CalculateAverage() float64 {
	if len(s.Grades) == 0 {
		return 0.0
	}

	var total float64
	for _, grade := range s.Grades {
		total += grade
	}
	return total / float64(len(s.Grades))
}

// FormattedSummary returns a user-friendly string summarizing the student and their grades.
func (s *Student) FormattedSummary() string {
	avg := s.CalculateAverage()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ID: %-5d | Name: %-20s | Average: %6.2f\n", s.ID, s.Name, avg))

	if len(s.Grades) == 0 {
		sb.WriteString("   No grades recorded yet.\n")
	} else {
		// Sort subjects alphabetically for consistent output
		subjects := make([]string, 0, len(s.Grades))
		for subj := range s.Grades {
			subjects = append(subjects, subj)
		}
		sort.Strings(subjects)

		sb.WriteString("   Grades:\n")
		for _, subj := range subjects {
			sb.WriteString(fmt.Sprintf("     - %-15s: %6.2f\n", subj, s.Grades[subj]))
		}
	}

	return sb.String()
}
