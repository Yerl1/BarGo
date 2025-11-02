package ports

import (
	"context"

	"backend/internal/auth-service/core/domain/dto"
	"backend/internal/auth-service/core/domain/models"
)

type IAuthRepo interface {
	CreateUser(ctx context.Context, user *dto.SignupUser) (string, error)
	CheckUserPassword(ctx context.Context, email, password string) (*models.User, error)
}
