package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/shaneajeffery/udacity-go-crm-backend/db"
)

var ctx = context.Background()

var (
	CustomerRegex = regexp.MustCompile(`^/customers/*$`)
	// Looking for UUID for the Customer ID.
	CustomerRegexWithID = regexp.MustCompile(`^/customers/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

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

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

type CustomersHandler struct{}

func (c *CustomersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

func (c *CustomersHandler) getCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := db.GetDbConn().GetCustomers(ctx)

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func (c *CustomersHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	matches := CustomerRegexWithID.FindStringSubmatch(r.URL.Path)

	// If the regex fails to get the URL base + the customerId arg, then throw err.
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	customer, err := db.GetDbConn().GetCustomer(ctx, matches[1])

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

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
	matches := CustomerRegexWithID.FindStringSubmatch(r.URL.Path)

	// If the regex fails to get the URL base + the customerId arg, then throw err.
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Check to see if the user exists before we try to delete it.
	_, err := db.GetDbConn().GetCustomer(ctx, matches[1])

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	if err := db.GetDbConn().DeleteCustomer(ctx, matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
