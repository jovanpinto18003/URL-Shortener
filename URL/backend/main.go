package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var db *sql.DB

const (
	host = "http://localhost:8080"
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./url.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Database initialized successfully")
}

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	data := hasher.Sum(nil)
	hash := base64.URLEncoding.EncodeToString(data)
	return hash[:8]
}

func createURL(originalURL string) (string, error) {
	shortURL := generateShortURL(originalURL)

	stmt, err := db.Prepare("INSERT INTO urls (original_url, short_url) VALUES (?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(originalURL, shortURL)
	if err != nil {
		return "", err
	}

	fullShortURL := fmt.Sprintf("%s/redirect/%s", host, shortURL)
	return fullShortURL, nil
}

func getURL(shortURL string) (string, error) {
	var originalURL string

	err := db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("URL not found")
		}
		return "", err
	}

	return originalURL, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at root page")
	fmt.Fprintf(w, "Hello, world!")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL, err := createURL(data.URL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/redirect/"):]
	originalURL, err := getURL(shortURL)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	initDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", RootPageURL)
	mux.HandleFunc("/shorten", ShortURLHandler)
	mux.HandleFunc("/redirect/", redirectURLHandler)

	handler := cors.Default().Handler(mux)
	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
