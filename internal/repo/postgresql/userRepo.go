package postgresql

import (
	"account-management/internal/models"
	"account-management/internal/repo/repoerrs"
	"account-management/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type UserRepo struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewUserRepo(db *pgxpool.Pool, logger *logging.Logger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "t", ""), "\n", " ")
}

func (u *UserRepo) CreateUser(ctx context.Context, user models.User) (string, error) {
	query := `
		INSERT INTO public.users (username, password)
		VALUES ($1, $2)
		RETURNING id;
	`

	u.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(query)))

	fmt.Println(user.Username, user.Password)

	var id string
	err := u.db.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", repoerrs.ErrAlreadyExists
			}
			return "", fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
		}
	}

	return id, nil
}

func (u *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (models.User, error) {
	query := `
		SELECT id, username, password, created_at 
		FROM public.users 
		WHERE username = $1 AND password = $2
	`

	u.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(query)))

	var user models.User
	err := u.db.QueryRow(ctx, query, username, password).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, repoerrs.ErrNotFound
		}
		return models.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	fmt.Println(user.ID, user.Username, user.Password, user.CreatedAt)

	return user, nil
}
