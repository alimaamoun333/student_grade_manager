package main

import (
	"fmt"
	"log"
)

const dataFilePath = "data/students.json"

func main() {
	storage := NewStorage(dataFilePath)
	manager := NewManager()

	// Load existing data from persistence file
	data, err := storage.Load()
	if err != nil {
		log.Printf("Warning: could not load initial student data (%v). Starting with empty records.\n", err)
	} else {
		manager.LoadData(data)
		fmt.Printf("Loaded %d student record(s) from %s\n", len(data), dataFilePath)
	}

	// Initialize and run the CLI menu loop
	menu := NewMenu(manager, storage)
	menu.Run()
}
