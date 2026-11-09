package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Manager manages the collection of students and operations on them.
type Manager struct {
	Students map[int]Student `json:"students"`
}

// NewManager creates and returns a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		Students: make(map[int]Student),
	}
}

// AddStudent adds a new student to the collection.
// Returns an error if a student with the same ID already exists or name is empty.
func (m *Manager) AddStudent(student Student) error {
	if student.ID <= 0 {
		return errors.New("student ID must be a positive integer")
	}
	if strings.TrimSpace(student.Name) == "" {
		return errors.New("student name cannot be empty")
	}
	if _, exists := m.Students[student.ID]; exists {
		return fmt.Errorf("student with ID %d already exists", student.ID)
	}
	if student.Grades == nil {
		student.Grades = make(map[string]float64)
	}
	m.Students[student.ID] = student
	return nil
}

// FindStudent looks up a student by their ID.
func (m *Manager) FindStudent(id int) (Student, error) {
	student, exists := m.Students[id]
	if !exists {
		return Student{}, fmt.Errorf("student with ID %d not found", id)
	}
	return student, nil
}

// GetAllStudents returns a slice of all students sorted by ID.
func (m *Manager) GetAllStudents() []Student {
	list := make([]Student, 0, len(m.Students))
	for _, student := range m.Students {
		list = append(list, student)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	return list
}

// AddGradeToStudent adds or updates a grade for a student.
func (m *Manager) AddGradeToStudent(id int, subject string, grade float64) error {
	student, exists := m.Students[id]
	if !exists {
		return fmt.Errorf("student with ID %d not found", id)
	}

	if err := student.AddGrade(subject, grade); err != nil {
		return err
	}

	m.Students[id] = student
	return nil
}

// LoadData replaces the current manager students map.
func (m *Manager) LoadData(data map[int]Student) {
	if data == nil {
		m.Students = make(map[int]Student)
	} else {
		m.Students = data
	}
}

// GetData returns the current map of students.
func (m *Manager) GetData() map[int]Student {
	return m.Students
}
