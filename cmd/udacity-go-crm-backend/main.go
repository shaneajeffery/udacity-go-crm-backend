package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/db"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.DbConn(context.Background(), os.Getenv("POSTGRES_DB_URL"))

	mux := routes.NewRouter()

	fmt.Printf("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}