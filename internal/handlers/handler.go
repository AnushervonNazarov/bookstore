package handlers

import (
	"bookstore/internal/db"
	"bookstore/internal/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Get all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, author, published_year FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
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
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	row := db.DB.QueryRow("SELECT id, title, author, published_year  FROM books WHERE id = $1", params["id"])

	var book models.Book
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
func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
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
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book models.Book
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
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.DB.Exec("DELETE FROM books WHERE id = $1", params["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
