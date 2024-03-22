package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shaneajeffery/udacity-go-crm-backend/internal/models"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func DbConn(ctx context.Context, connString string) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			return
		}

		pgInstance = &postgres{db}
	})

}

func GetDbConn() *postgres {
	return pgInstance
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetCustomers(ctx context.Context) ([]models.Customer, error) {
	query := `SELECT * FROM customers`

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Customer])
}

func (pg *postgres) GetCustomer(ctx context.Context, customerId string) (models.Customer, error) {
	query := `SELECT * FROM customers WHERE id = $1`

	rows, err := pg.db.Query(ctx, query, customerId)

	if err != nil {
		return models.Customer{}, fmt.Errorf("unable to query users: %w", err)
	}

	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Customer])
}

func (pg *postgres) DeleteCustomer(ctx context.Context, customerId string) error {
	query := `DELETE FROM customers WHERE id = $1`

	_, err := pg.db.Query(ctx, query, customerId)

	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}

func (pg *postgres) CreateCustomer(ctx context.Context, customer models.Customer) error {
	query := `INSERT INTO customers (name, role, email, phone, contacted) 
				VALUES (@name, @role, @email, @phone, @contacted)`

	args := pgx.NamedArgs{
		"name":      customer.Name,
		"role":      customer.Role,
		"email":     customer.Email,
		"phone":     customer.Phone,
		"contacted": customer.Contacted,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (pg *postgres) UpdateCustomer(ctx context.Context, customerId string, customer models.Customer) error {
	query := `UPDATE customers SET name = @name, role = @role, email = @email, phone = @phone, contacted = @contacted WHERE id = @id`

	args := pgx.NamedArgs{
		"name":      customer.Name,
		"role":      customer.Role,
		"email":     customer.Email,
		"phone":     customer.Phone,
		"contacted": customer.Contacted,
		"id":        customerId,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}
