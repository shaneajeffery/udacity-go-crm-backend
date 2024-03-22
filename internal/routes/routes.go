package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/shaneajeffery/udacity-go-crm-backend/internal/db"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/models"
)

var ctx = context.Background()

var (
	CustomerRegex = regexp.MustCompile(`^/customers/*$`)
	// Looking for UUID for the Customer ID.
	CustomerRegexWithID = regexp.MustCompile(`^/customers/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/customers", &CustomersHandler{})
	mux.Handle("/customers/", &CustomersHandler{})

	return mux
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
	This API allows access to a Customers CRM API.  Your available endpoints are:

	GET    /customers       -- Returns all customers.
	GET    /customers/:id   -- Returns specified customer.
	POST   /customers       -- Creates new customer.
	PUT    /customers/:id   -- Updates specified customer.
	DELETE /customers/      -- Deletes specified customer.
	`))
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
	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	err := db.GetDbConn().CreateCustomer(ctx, customer)

	if err != nil {
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *CustomersHandler) updateCustomer(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Update customer rout"))
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
		fmt.Println(err)
		NotFoundHandler(w, r)
		return
	}

	if err := db.GetDbConn().DeleteCustomer(ctx, matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
