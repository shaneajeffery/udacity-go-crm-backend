package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/routes"
)

func main() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := routes.NewRouter()

	fmt.Println("Server listening on http://localhost:8082")
	http.ListenAndServe(":8082", mux)
	if err != nil {
		panic(err)
	}
}
