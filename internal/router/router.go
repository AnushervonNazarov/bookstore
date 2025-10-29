package router

import (
	"bookstore/internal/handlers"
	"bookstore/internal/middleware"
	"bookstore/internal/user"
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

func RunRouters() {
	// Router for routes that require authentication
	authenticatedRouter := mux.NewRouter()
	authenticatedRouter.Use(middleware.AuthMiddleware)
	authenticatedRouter.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	authenticatedRouter.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	authenticatedRouter.HandleFunc("/books", handlers.AddBook).Methods("POST")
	authenticatedRouter.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	authenticatedRouter.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

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
