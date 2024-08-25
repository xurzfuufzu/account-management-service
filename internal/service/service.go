package service

import (
	"account-management/internal/repo"
	"account-management/pkg/hasher"
	"context"
	"time"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (string, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(tokenString string) (string, error)
}

type Services struct {
	Auth
}

type ServiceDependencies struct {
	Repos     *repo.Repositories
	Hasher    hasher.PasswordHasher
	SecretKey string
	TokenTTL  time.Duration
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User, deps.Hasher, deps.SecretKey, deps.TokenTTL),
	}
}
