package main

import (
	"auth/internal/database/postgres"
	"fmt"
	"net/http"
)

func main() {
	authService, err := postgres.NewAuthService()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		authService.Login(w, r)
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		authService.Verify(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
