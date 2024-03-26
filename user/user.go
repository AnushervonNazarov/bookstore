package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"bookstore/db"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// User struct to represent a user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// UserCredentials struct to handle user login credentials
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Map to store user sessions
var sessions = make(map[string]int)

// Register a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userCred UserCredentials
	err := json.NewDecoder(r.Body).Decode(&userCred)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user := User{
		Username: userCred.Username,
		Password: userCred.Password,
	}

	err = ValidateRegister(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if username already exists
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&count)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Insert new user into database
	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	var retrievedUser User
	err = db.DB.QueryRow("SELECT id, password FROM users WHERE username = $1", user.Username).Scan(&retrievedUser.ID, &retrievedUser.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Compare stored password hash with provided password
	err = bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Generate session token
	sessionToken := uuid.New().String()
	sessions[sessionToken] = retrievedUser.ID

	// Set token in response header
	w.Header().Set("Authorization", "Bearer "+sessionToken)
}

func ValidateRegister(u *User) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)

	if u.Username == "" {
		return errors.New("no login provided")
	}

	if u.Password == "" {
		return errors.New("no password provided")
	}

	if len(strings.Split(u.Username, " ")) > 1 {
		return errors.New("login should be one word")
	}

	return nil
}
