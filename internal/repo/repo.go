package repo

import (
	"account-management/internal/models"
	"account-management/internal/repo/postgresql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (models.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(pgx *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: postgresql.NewUserRepo(pgx),
	}
}
