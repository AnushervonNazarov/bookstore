package main

import (
	"bookstore/internal/db"
	"bookstore/internal/user"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Book struct to represent a book
type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
}

func main() {
	var err error

	// Read configuration from file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer configFile.Close()

	var config Config
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config JSON:", err)
	}

	// Construct connection string
	connStr := "host=" + config.Database.Host +
		" port=" + strconv.Itoa(config.Database.Port) +
		" user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.DBName +
		" sslmode=" + config.Database.SSLMode

	// Connect to PostgreSQL
	db.InitDB(connStr)
	defer db.DB.Close()

	// Router for routes that require authentication
	authenticatedRouter := mux.NewRouter()
	authenticatedRouter.Use(AuthMiddleware)
	authenticatedRouter.HandleFunc("/books", getBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/books/{id}", getBook).Methods("GET")
	authenticatedRouter.HandleFunc("/books", addBook).Methods("POST")
	authenticatedRouter.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	authenticatedRouter.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Router for routes that don't require authentication
	nonAuthenticatedRouter := mux.NewRouter()
	nonAuthenticatedRouter.HandleFunc("/register", user.RegisterUser).Methods("POST")
	nonAuthenticatedRouter.HandleFunc("/login", user.LoginUser).Methods("POST")

	// Combine both routers
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/books").Handler(authenticatedRouter)
	mainRouter.PathPrefix("/").Handler(nonAuthenticatedRouter)

	log.Fatal(http.ListenAndServe(":8000", mainRouter))

}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, author, published_year FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	row := db.DB.QueryRow("SELECT id, title, author, published_year  FROM books WHERE id = $1", params["id"])

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Add a new book
func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("INSERT INTO books (title, author, published_year) VALUES ($1, $2, $3)", book.Title, book.Author, book.PublishedYear)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE books SET title = $1, author = $2, published_year = $3 WHERE id = $4", book.Title, book.Author, book.PublishedYear, params["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.DB.Exec("DELETE FROM books WHERE id = $1", params["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Middleware function to handle authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		token := strings.Split(authHeader, "Bearer ")
		if len(token) != 2 {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Here you would typically validate the token.
		// For this example, we assume a simple token validation.

		// If token is valid, proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
