package main

import (
	"log"
	"net/http"
	"os"

	"cms-server/internal/database"
	"cms-server/internal/handlers"
	"cms-server/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	database.ConnectMongo()

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Register public routes
	registerPublicRoutes(r)

	// Register private routes
	registerPrivateRoutes(r)

	// Start the server
	startServer(r)
}

func registerPublicRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutUserHandler).Methods("POST")
	r.HandleFunc("/contents", handlers.GetContentsHandler).Methods("GET")
	r.HandleFunc("/stacks", handlers.GetStacksHandler).Methods("GET")
}

func registerPrivateRoutes(r *mux.Router) {
	r.Handle("/content", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateContentHandler))).Methods("POST")
	r.Handle("/content", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetContentHandler))).Methods("GET")
	r.Handle("/content/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.EditContentHandler))).Methods("PUT")
	r.Handle("/content/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteContentHandler))).Methods("DELETE")

	r.Handle("/stacks", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateStackHandler))).Methods("POST")
	r.Handle("/stacks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.EditStackHandler))).Methods("PUT")
	r.Handle("/stacks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteStackHandler))).Methods("DELETE")
}

func startServer(r *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
