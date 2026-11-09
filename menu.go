package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Menu handles the user interaction via command line interface.
type Menu struct {
	manager *Manager
	storage *Storage
	scanner *bufio.Scanner
}

// NewMenu creates a new Menu controller.
func NewMenu(manager *Manager, storage *Storage) *Menu {
	return &Menu{
		manager: manager,
		storage: storage,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run starts the main interaction loop.
func (m *Menu) Run() {
	fmt.Println("==================================================")
	fmt.Println("       WELCOME TO STUDENT GRADE MANAGER           ")
	fmt.Println("==================================================")

	for {
		m.displayOptions()
		fmt.Print("\nEnter your choice (1-5): ")
		choice := m.readLine()

		switch choice {
		case "1":
			m.handleListStudents()
		case "2":
			m.handleAddStudent()
		case "3":
			m.handleFindStudent()
		case "4":
			m.handleAddGrade()
		case "5", "exit", "quit", "q":
			m.handleExit()
			return
		default:
			fmt.Println("❌ Invalid choice. Please enter a number between 1 and 5.")
		}

		fmt.Println("\nPress Enter to return to the menu...")
		m.readLine()
	}
}

func (m *Menu) displayOptions() {
	fmt.Println("\n---------------- MENU OPTIONS ----------------")
	fmt.Println(" 1. 📋 List All Students & Grades")
	fmt.Println(" 2. ➕ Add New Student")
	fmt.Println(" 3. 🔍 Find Student by ID")
	fmt.Println(" 4. 📝 Add / Update Subject Grade")
	fmt.Println(" 5. 🚪 Save and Exit")
	fmt.Println("----------------------------------------------")
}

func (m *Menu) handleListStudents() {
	students := m.manager.GetAllStudents()
	if len(students) == 0 {
		fmt.Println("ℹ️  No students registered yet.")
		return
	}

	fmt.Println("\n================ STUDENT LIST ================")
	for _, student := range students {
		fmt.Print(student.FormattedSummary())
		fmt.Println("----------------------------------------------")
	}
}

func (m *Menu) handleAddStudent() {
	fmt.Println("\n--- Add New Student ---")
	fmt.Print("Enter Student ID (positive integer): ")
	idStr := m.readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Println("❌ Error: Invalid Student ID. Must be a positive integer.")
		return
	}

	fmt.Print("Enter Student Full Name: ")
	name := m.readLine()
	if strings.TrimSpace(name) == "" {
		fmt.Println("❌ Error: Student name cannot be empty.")
		return
	}

	student := *NewStudent(id, name)
	if err := m.manager.AddStudent(student); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully added student: %s (ID: %d)\n", student.Name, student.ID)
}

func (m *Menu) handleFindStudent() {
	fmt.Println("\n--- Find Student by ID ---")
	fmt.Print("Enter Student ID: ")
	idStr := m.readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("❌ Error: Invalid ID.")
		return
	}

	student, err := m.manager.FindStudent(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Println("\n--- Student Found ---")
	fmt.Print(student.FormattedSummary())
}

func (m *Menu) handleAddGrade() {
	fmt.Println("\n--- Add / Update Subject Grade ---")
	fmt.Print("Enter Student ID: ")
	idStr := m.readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("❌ Error: Invalid Student ID.")
		return
	}

	student, err := m.manager.FindStudent(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("Selected Student: %s (ID: %d)\n", student.Name, student.ID)
	fmt.Print("Enter Subject Name (e.g. Math, Physics, CS): ")
	subject := m.readLine()
	if strings.TrimSpace(subject) == "" {
		fmt.Println("❌ Error: Subject name cannot be empty.")
		return
	}

	fmt.Print("Enter Grade (0.0 - 100.0): ")
	gradeStr := m.readLine()
	grade, err := strconv.ParseFloat(gradeStr, 64)
	if err != nil || grade < 0 || grade > 100 {
		fmt.Println("❌ Error: Invalid grade. Must be a number between 0 and 100.")
		return
	}

	if err := m.manager.AddGradeToStudent(id, subject, grade); err != nil {
		fmt.Printf("❌ Error adding grade: %v\n", err)
		return
	}

	fmt.Printf("✅ Grade %.2f added for subject '%s' to student %s.\n", grade, subject, student.Name)
}

func (m *Menu) handleExit() {
	fmt.Println("\nSaving data before exit...")
	if err := m.storage.Save(m.manager.GetData()); err != nil {
		fmt.Printf("⚠️ Warning: Failed to auto-save data: %v\n", err)
	} else {
		fmt.Println("💾 Data saved successfully.")
	}
	fmt.Println("Goodbye! 👋")
}

func (m *Menu) readLine() string {
	if m.scanner.Scan() {
		return strings.TrimSpace(m.scanner.Text())
	}
	return ""
}
