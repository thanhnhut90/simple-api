package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB struct holds the connection and the mutex for thread safety
type DB struct {
	Conn *sql.DB
	Type string
	Mu   sync.Mutex
}

// StringData represents a string record
type StringData struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// Initialize the DB connection based on the DB type
func InitDB(dbType string) (*DB, error) {
	var conn *sql.DB
	var err error

	switch dbType {
	case "postgres":
		// Fetch Postgres credentials from environment variables
		dbUser := os.Getenv("POSTGRES_USER")
		dbPassword := os.Getenv("POSTGRES_PASSWORD")
		dbHost := os.Getenv("POSTGRES_HOST")
		dbPort := os.Getenv("POSTGRES_PORT")
		dbName := os.Getenv("POSTGRES_DB")

		// If any environment variable is missing, return an error
		if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
			return nil, fmt.Errorf("missing required PostgreSQL environment variables")
		}

		// Connection string for PostgreSQL
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
		conn, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to PostgreSQL: %v", err)
		}
		err = conn.Ping()
		if err != nil {
			return nil, fmt.Errorf("Failed to ping PostgreSQL database: %v", err)
		}
		fmt.Println("Connected to PostgreSQL successfully!")

	case "sqlite":
		// SQLite connection
		conn, err = sql.Open("sqlite3", "./stringdb.sqlite")
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to SQLite: %v", err)
		}
		err = conn.Ping()
		if err != nil {
			return nil, fmt.Errorf("Failed to ping SQLite database: %v", err)
		}
		fmt.Println("Connected to SQLite successfully!")

	default:
		return nil, fmt.Errorf("Unsupported DB_TYPE: %v. Use 'postgres' or 'sqlite'.", dbType)
	}

	return &DB{
		Conn: conn,
		Type: dbType,
	}, nil
}

// Create the strings table if it does not exist
func (db *DB) CreateTable(name string) error {
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT NOT NULL
	);
	`, name)
	_, err := db.Conn.Exec(createTableSQL)
	return err
}

// Insert a new string and return its ID
func (db *DB) CreateString(value string) (int, error) {
	db.Mu.Lock() // Lock to avoid race conditions
	defer db.Mu.Unlock()

	var id int
	var err error
	switch db.Type {
	case "postgres":
		err = db.Conn.QueryRow("INSERT INTO strings(value) VALUES($1) RETURNING id", value).Scan(&id)
	case "sqlite":
		err = db.Conn.QueryRow("INSERT INTO strings(value) VALUES(?) RETURNING id", value).Scan(&id)
	}
	return id, err
}

// Get all strings from the database
func (db *DB) GetStrings() ([]StringData, error) {
	db.Mu.Lock() // Lock to avoid race conditions
	defer db.Mu.Unlock()

	rows, err := db.Conn.Query("SELECT id, value FROM strings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strings []StringData
	for rows.Next() {
		var s StringData
		if err := rows.Scan(&s.ID, &s.Value); err != nil {
			return nil, err
		}
		strings = append(strings, s)
	}

	return strings, nil
}

// Get a string by its ID
func (db *DB) GetStringByID(id int) (StringData, error) {
	db.Mu.Lock() // Lock to avoid race conditions
	defer db.Mu.Unlock()

	var s StringData
	switch db.Type {
	case "postgres":
		err := db.Conn.QueryRow("SELECT id, value FROM strings WHERE id=$1", id).Scan(&s.ID, &s.Value)
		if err != nil {
			return s, err
		}
	case "sqlite":
		err := db.Conn.QueryRow("SELECT id, value FROM strings WHERE id=?", id).Scan(&s.ID, &s.Value)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

// Update a string by ID
func (db *DB) UpdateString(id int, value string) error {
	db.Mu.Lock() // Lock to avoid race conditions
	defer db.Mu.Unlock()

	_, err := db.Conn.Exec("UPDATE strings SET value=? WHERE id=?", value, id)
	return err
}

// Delete a string by ID
func (db *DB) DeleteString(id int) error {
	db.Mu.Lock() // Lock to avoid race conditions
	defer db.Mu.Unlock()

	_, err := db.Conn.Exec("DELETE FROM strings WHERE id=?", id)
	return err
}
