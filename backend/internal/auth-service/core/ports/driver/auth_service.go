package ports

import (
	"context"

	"backend/internal/auth-service/core/domain/dto"
)

type IAuthService interface {
	Signup(ctx context.Context, user *dto.SignupUser) (dto.Resp, error)
	Login(ctx context.Context, user *dto.LoginUser) (dto.Resp, error)
}
