package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/models"
)

var (
	CustomerRegex = regexp.MustCompile(`^/customers/*$`)
	// Looking for UUID for the Customer ID.
	CustomerRegexWithID = regexp.MustCompile(`^/customers/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

var customerList = []models.Customer{
	{
		ID:        "741f965c-42d6-43b0-9389-648295925c0f",
		Name:      "Customer 1",
		Role:      "Admin",
		Email:     "customer1@customer.com",
		Phone:     "555-555-5555",
		Contacted: true,
	},
	{
		ID:        "4ab256fc-6367-444e-a205-e2bb2832874e",
		Name:      "Customer 2",
		Role:      "Customer",
		Email:     "customer2@customer.com",
		Phone:     "555-555-5556",
		Contacted: false,
	},
	{
		ID:        "4dac8255-0ae0-4216-abbb-15a85df4af36",
		Name:      "Customer 3",
		Role:      "Customer",
		Email:     "customer3@customer.com",
		Phone:     "555-555-5557",
		Contacted: false,
	},
}

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

func (c *CustomersHandler) getCustomers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customerList)
}

func (c *CustomersHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	matches := CustomerRegexWithID.FindStringSubmatch(r.URL.Path)

	// If the regex fails to get the URL base + the customerId arg, then throw err.
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	foundCustomerIndex := IndexOf(matches[1], customerList)

	if foundCustomerIndex == -1 {
		NotFoundHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customerList[foundCustomerIndex])
}

func (c *CustomersHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	// Generate unique ID for customer.
	newCustomer.ID = uuid.New().String()

	customerList = append(customerList, newCustomer)

	w.WriteHeader(http.StatusCreated)
}

func (c *CustomersHandler) updateCustomer(w http.ResponseWriter, r *http.Request) {
	matches := CustomerRegexWithID.FindStringSubmatch(r.URL.Path)

	// If the regex fails to get the URL base + the customerId arg, then throw err.
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	var updatedCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	foundCustomerIndex := IndexOf(matches[1], customerList)

	if foundCustomerIndex == -1 {
		NotFoundHandler(w, r)
		return
	}

	customerList[foundCustomerIndex].Name = updatedCustomer.Name
	customerList[foundCustomerIndex].Role = updatedCustomer.Role
	customerList[foundCustomerIndex].Email = updatedCustomer.Email
	customerList[foundCustomerIndex].Phone = updatedCustomer.Phone
	customerList[foundCustomerIndex].Contacted = updatedCustomer.Contacted

	w.WriteHeader(http.StatusOK)
}

func (c *CustomersHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	matches := CustomerRegexWithID.FindStringSubmatch(r.URL.Path)

	// If the regex fails to get the URL base + the customerId arg, then throw err.
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	foundCustomerIndex := IndexOf(matches[1], customerList)

	if foundCustomerIndex == -1 {
		NotFoundHandler(w, r)
		return
	}

	// Re-slicing the slice to remove the foundCustomerIndex so that the item is actually deleted.
	customerList = append(customerList[:foundCustomerIndex], customerList[foundCustomerIndex+1:]...)

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

// This function will check to see if the Customer ID exists
// in the passed Customer List.
// If it does exist, then the slice index will be returned.
// If not, then -1 will be returned.
func IndexOf(element string, data []models.Customer) int {
	for k, v := range data {
		if element == v.ID {
			return k
		}
	}

	return -1
}
