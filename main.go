package main

import (
	"log"
	"net/http"

	"docker-image-scanner/handlers"
	"docker-image-scanner/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// ✅ Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// ✅ Setup router
	r := mux.NewRouter()

	// 🟢 Public routes
	r.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// 🔒 Protected routes
	s := r.PathPrefix("/api").Subrouter()
	s.Use(middleware.JWTMiddleware)
	s.HandleFunc("/scan", handlers.ScanDockerImageHandler).Methods("POST")

	// ✅ Start server
	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
