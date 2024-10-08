package client

import (
	"account-management/config"
	"account-management/pkg/logging"
	"account-management/pkg/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func GetURL(db *config.DB) string {
	if db.URL == "" {
		db.URL =
			"postgres://" +
				db.Username + ":" +
				db.Password + "@" +
				db.Host + ":" +
				db.Port + "/" +
				db.Database +
				"?sslmode=disable"
	}
	return db.URL
}

func NewClient(ctx context.Context, maxAttempts int, db config.DB) (pool *pgxpool.Pool, err error) {
	logger := logging.GetLogger()
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, GetURL(&db))
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logger.Fatal("error to connection to postqresql")
	}

	logger.Info("Success connection to db")

	return pool, nil
}
