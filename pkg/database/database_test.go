package database_test

import (
	"testing"

	"github.com/thanhnhut90/simple-api/pkg/database"
	_ "github.com/mattn/go-sqlite3"
)

func TestDatabaseOperations(t *testing.T) {
	// Initialize the SQLite DB
	db, err := database.InitDB("sqlite")
	if err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer db.Conn.Close()

	// Create the table
	err = db.CreateTable("unittest")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert a string
	value := "Hello, World!"
	id, err := db.CreateString(value)
	if err != nil {
		t.Fatalf("Failed to insert string: %v", err)
	}
	if id <= 0 {
		t.Fatalf("Inserted string should have a valid ID, got: %d", id)
	}

	// Retrieve the string by ID
	retrievedString, err := db.GetStringByID(id)
	if err != nil {
		t.Fatalf("Failed to retrieve string: %v", err)
	}
	if retrievedString.Value != value {
		t.Errorf("Retrieved value mismatch. Expected: %s, Got: %s", value, retrievedString.Value)
	}

	// Update the string
	updatedValue := "Updated String"
	err = db.UpdateString(id, updatedValue)
	if err != nil {
		t.Fatalf("Failed to update string: %v", err)
	}

	// Retrieve the updated string
	updatedString, err := db.GetStringByID(id)
	if err != nil {
		t.Fatalf("Failed to retrieve updated string: %v", err)
	}
	if updatedString.Value != updatedValue {
		t.Errorf("Updated value mismatch. Expected: %s, Got: %s", updatedValue, updatedString.Value)
	}

	// Delete the string
	err = db.DeleteString(id)
	if err != nil {
		t.Fatalf("Failed to delete string: %v", err)
	}

	// Try to retrieve the deleted string (should return error)
	_, err = db.GetStringByID(id)
	if err == nil {
		t.Error("Should not be able to retrieve deleted string")
	}
}

