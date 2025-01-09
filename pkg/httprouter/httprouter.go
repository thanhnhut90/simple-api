package httprouter

import (
	"encoding/json"
	"github.com/thanhnhut90/simple-api/pkg/database"
	"net/http"
	"strconv"
)

// HTTPRouter struct holds the DB dependency
type HTTPRouter struct {
	DB *database.DB
}

// Struct to hold the string data
type StringData struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// Create a new string in the database
func (r *HTTPRouter) CreateString(w http.ResponseWriter, req *http.Request) {
	var newString StringData
	err := json.NewDecoder(req.Body).Decode(&newString)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert the string into the database
	id, err := r.DB.CreateString(newString.Value)
	if err != nil {
		http.Error(w, "Failed to insert string", http.StatusInternalServerError)
		return
	}

	newString.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newString)
}

// Get all strings from the database
func (r *HTTPRouter) GetStrings(w http.ResponseWriter, req *http.Request) {
	strings, err := r.DB.GetStrings()
	if err != nil {
		http.Error(w, "Failed to fetch strings", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(strings)
}

// Get a string by ID
func (r *HTTPRouter) GetStringByID(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[len("/api/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	s, err := r.DB.GetStringByID(id)
	if err != nil {
		http.Error(w, "String not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// Update a string by ID
func (r *HTTPRouter) UpdateString(w http.ResponseWriter, req *http.Request) {
	var updatedString StringData
	err := json.NewDecoder(req.Body).Decode(&updatedString)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = r.DB.UpdateString(updatedString.ID, updatedString.Value)
	if err != nil {
		http.Error(w, "Failed to update string", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedString)
}

// Delete a string by ID
func (r *HTTPRouter) DeleteString(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = r.DB.DeleteString(id)
	if err != nil {
		http.Error(w, "Failed to delete string", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Setup all the HTTP routes
func (r *HTTPRouter) SetupRoutes() {
	http.HandleFunc("/api", r.GetStrings)
	http.HandleFunc("/api/create", r.CreateString)
	http.HandleFunc("/api/update", r.UpdateString)
	http.HandleFunc("/api/delete", r.DeleteString)
	http.HandleFunc("/api/", r.GetStringByID)
}
