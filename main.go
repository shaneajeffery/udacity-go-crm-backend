package main

import (
	"net/http"
	"regexp"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/customers", &CustomersHandler{})
	mux.Handle("/customers/", &CustomersHandler{})

	http.ListenAndServe(":8080", mux)
}

type CustomersHandler struct{}

func (c *CustomersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		CustomerRegex       = regexp.MustCompile(`^/customers/*$`)
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
	w.Write([]byte("Get customers route"))
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
