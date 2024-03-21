package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/shaneajeffery/udacity-go-crm-backend/db"
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

type CustomersHandler struct {
	pgInstance *pgxpool.Pool
}

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

func (c *CustomersHandler) getCustomers(w http.ResponseWriter, r *http.Request) {
	customers, _ := db.GetDbConn(context.Background()).GetCustomers(context.Background())

	fmt.Print("hello")

	for _, customer := range customers {
		fmt.Printf("%#v\n", customer)
	}
}

func (c *CustomersHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get customer route"))

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
