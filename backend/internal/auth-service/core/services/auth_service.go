package services

import (
	"context"
	"time"

	"backend/internal/auth-service/core/domain/dto"
	ports "backend/internal/auth-service/core/ports/driven"
	"backend/internal/mylogger"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	mylog     mylogger.Logger
	AuthRepo  ports.IAuthRepo
	ctx       context.Context
	jwtSecret string
}

func NewAuthService(ctx context.Context, authRepo ports.IAuthRepo, jwtSecret string, mylog mylogger.Logger) *AuthService {
	return &AuthService{
		ctx:       ctx,
		AuthRepo:  authRepo,
		jwtSecret: jwtSecret,
		mylog:     mylog,
	}
}

func (s *AuthService) Signup(ctx context.Context, user *dto.SignupUser) (dto.Resp, error) {
	// Implement signup logic
	log := s.mylog.With().Str("service", "Signup").Logger()
	id, err := s.AuthRepo.CreateUser(ctx, user)
	if err != nil {
		log.Err(err).Msg("Failed to create user")
		return dto.Resp{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": *user.FirstName + " " + *user.LastName,
		"role":     *user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Err(err).Msg("Failed to sign JWT token")
		return dto.Resp{}, err
	}
	log.Info().Str("user_id", id).Msg("User signed up successfully")
	log.Debug().Str("token", tokenString).Msg("Generated JWT token")

	resp := dto.Resp{
		Token: tokenString,
		User: dto.UserInfo{
			Role:      *user.Role,
			FirstName: *user.FirstName,
			LastName:  *user.LastName,
			Email:     *user.Email,
			Phone:     *user.Phone,
		},
	}
	return resp, nil
}

func (s *AuthService) Login(ctx context.Context, user *dto.LoginUser) (dto.Resp, error) {
	// Implement login logic
	log := s.mylog.With().Str("service", "Login").Logger()

	userInfo, err := s.AuthRepo.CheckUserPassword(ctx, *user.Email, *user.Password)
	if err != nil {
		log.Err(err).Msg("Failed to get user by email and password")
		return dto.Resp{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userInfo.UserID,
		"username": *userInfo.FirstName + " " + *userInfo.LastName,
		"role":     *userInfo.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Err(err).Msg("Failed to sign JWT token")
		return dto.Resp{}, err
	}
	log.Info().Str("user_id", userInfo.UserID).Msg("User logged in successfully")
	log.Debug().Str("token", tokenString).Msg("Generated JWT token")

	resp := dto.Resp{
		Token: tokenString,
		User: dto.UserInfo{
			Role:      *userInfo.Role,
			FirstName: *userInfo.FirstName,
			LastName:  *userInfo.LastName,
			Email:     *userInfo.Email,
			Phone:     *userInfo.Phone,
		},
	}

	return resp, nil
}
