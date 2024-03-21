package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/shaneajeffery/udacity-go-crm-backend/db"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.DbConn(context.Background(), os.Getenv("POSTGRES_DB_URL"))

	mux := http.NewServeMux()

	mux.Handle("/customers", &CustomersHandler{})
	mux.Handle("/customers/", &CustomersHandler{})

	http.ListenAndServe(":8080", mux)
}

type CustomersHandler struct{}

func (c *CustomersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		CustomerRegex = regexp.MustCompile(`^/customers/*$`)
		// Looking for UUID for the Customer ID.
		CustomerRegexWithID = regexp.MustCompile(`^/customers/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
	)

	switch {
	case r.Method == http.MethodPost && CustomerRegex.MatchString(r.URL.Path):
		c.createCustomer(w, r)
		return
	case r.Method == http.MethodGet && CustomerRegex.MatchString(r.URL.Path):
		c.getCustomers(w, r)
		return
	case r.Method == http.MethodGet && CustomerRegexWithID.MatchString(r.URL.Path):
		c.getCustomer(w, r)
		return
	case r.Method == http.MethodPut && CustomerRegexWithID.MatchString(r.URL.Path):
		c.updateCustomer(w, r)
		return
	case r.Method == http.MethodDelete && CustomerRegexWithID.MatchString(r.URL.Path):
		c.deleteCustomer(w, r)
		return
	default:
		return
	}
}

func (c *CustomersHandler) getCustomers(w http.ResponseWriter, _ *http.Request) {
	customers, _ := db.GetDbConn().GetCustomers(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func (c *CustomersHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/customers/")

	customer, _ := db.GetDbConn().GetCustomer(ctx, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customer)
}

func (c *CustomersHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create customer route"))
}

func (c *CustomersHandler) updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update customer route"))
}

func (c *CustomersHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete customer route"))
}
