package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shaneajeffery/udacity-go-crm-backend/db/models"
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

	rows, err := pg.db.Query(ctx, query, customerId)

	fmt.Println(rows)

	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}
