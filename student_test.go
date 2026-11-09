package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStudentOperations(t *testing.T) {
	student := NewStudent(1, "Alice Smith")

	if student.Name != "Alice Smith" || student.ID != 1 {
		t.Fatalf("unexpected student info: ID=%d, Name=%s", student.ID, student.Name)
	}

	// Average with no grades should be 0
	if avg := student.CalculateAverage(); avg != 0.0 {
		t.Errorf("expected average 0.0 with no grades, got %.2f", avg)
	}

	if err := student.AddGrade("Math", 95.0); err != nil {
		t.Fatalf("failed to add grade: %v", err)
	}
	if err := student.AddGrade("Science", 85.0); err != nil {
		t.Fatalf("failed to add grade: %v", err)
	}

	avg := student.CalculateAverage()
	if avg != 90.0 {
		t.Errorf("expected average 90.0, got %.2f", avg)
	}

	// Test invalid grade
	if err := student.AddGrade("History", 150.0); err == nil {
		t.Errorf("expected error when adding grade > 100")
	}
	if err := student.AddGrade("History", -5.0); err == nil {
		t.Errorf("expected error when adding grade < 0")
	}
	if err := student.AddGrade("", 80.0); err == nil {
		t.Errorf("expected error when subject is empty")
	}
}

func TestManagerOperations(t *testing.T) {
	manager := NewManager()

	s1 := *NewStudent(101, "John Doe")
	s2 := *NewStudent(102, "Jane Smith")

	if err := manager.AddStudent(s1); err != nil {
		t.Fatalf("failed to add s1: %v", err)
	}
	if err := manager.AddStudent(s2); err != nil {
		t.Fatalf("failed to add s2: %v", err)
	}

	// Duplicate ID check
	if err := manager.AddStudent(s1); err == nil {
		t.Errorf("expected error adding duplicate student ID")
	}

	// Invalid student checks
	if err := manager.AddStudent(*NewStudent(0, "Invalid ID")); err == nil {
		t.Errorf("expected error for non-positive ID")
	}
	if err := manager.AddStudent(*NewStudent(103, "   ")); err == nil {
		t.Errorf("expected error for empty name")
	}

	// Find by ID
	found, err := manager.FindStudent(101)
	if err != nil || found.Name != "John Doe" {
		t.Errorf("failed finding student 101: %v", err)
	}

	// Find non-existent ID
	if _, err := manager.FindStudent(999); err == nil {
		t.Errorf("expected error finding non-existent ID 999")
	}

	// List all students
	all := manager.GetAllStudents()
	if len(all) != 2 {
		t.Errorf("expected 2 students, got %d", len(all))
	}

	// Add grade via manager
	if err := manager.AddGradeToStudent(101, "Math", 80); err != nil {
		t.Fatalf("failed adding grade to student: %v", err)
	}
}

func TestStorageSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_students.json")

	storage := NewStorage(filePath)
	manager := NewManager()

	s := *NewStudent(1, "Test User")
	_ = s.AddGrade("English", 88.0)
	_ = manager.AddStudent(s)

	if err := storage.Save(manager.GetData()); err != nil {
		t.Fatalf("failed to save data: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file was not created")
	}

	loaded, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	if len(loaded) != 1 || loaded[1].Name != "Test User" || loaded[1].Grades["English"] != 88.0 {
		t.Fatalf("loaded data mismatch: %+v", loaded)
	}
}
