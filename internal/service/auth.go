package service

import (
	"account-management/internal/models"
	"account-management/internal/repo"
	"account-management/internal/repo/repoerrs"
	"account-management/pkg/hasher"
	"account-management/pkg/logging"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo       repo.User
	secretKey      string
	tokenTTL       time.Duration
	passwordHasher hasher.PasswordHasher
}

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, secretKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		secretKey:      secretKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (string, error) {
	user := models.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return "", ErrUserAlreadyExists
		}
		return "", ErrCannotCreateUser
	}

	return userID, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserNotFound
		}
		logger := logging.GetLogger()
		logger.Errorf("AuthService.GenerateToken: cannot get user: %v", err)
		return "", ErrCannotGetUser
	}

	claims := TokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (s *AuthService) ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", jwt.NewValidationError("invalid claims", jwt.ValidationErrorClaimsInvalid)
	}

	return claims.UserID, nil
}
