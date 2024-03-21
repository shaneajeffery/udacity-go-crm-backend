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

func GetDbConn(ctx context.Context) *postgres {
	return pgInstance
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetCustomers(ctx context.Context) ([]models.Customer, error) {
	query := `SELECT * FROM customers`

	rows, err := pg.db.Query(ctx, query)

	fmt.Print(err)

	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Customer])
}
