package db

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/auth-service/core/domain/dto"
	"backend/internal/auth-service/core/domain/models"
	"backend/internal/mylogger"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepo struct {
	ctx   context.Context
	mylog mylogger.Logger
	DB    *DB
}

func NewAuthRepo(ctx context.Context, db *DB, mylog mylogger.Logger) *AuthRepo {
	return &AuthRepo{
		ctx:   ctx,
		mylog: mylog,
		DB:    db,
	}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *dto.SignupUser) (string, error) {
	log := r.mylog.With().Str("repo", "CreateUser").Logger()
	q := `
	INSERT INTO users (email, password_hash, firstname, lastname, role, phone)
	VALUES ($1, crypt($2, gen_salt('bf')), $3, $4, $5, $6)
	RETURNING user_id	
	`
	row := r.DB.conn.QueryRow(ctx, q, user.Email, user.Password, user.FirstName, user.LastName, user.Role, user.Phone)

	var id string
	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				switch pgErr.ConstraintName {
				case "users_email_key":
					log.Error().Str("email", *user.Email).Msg("email already exists")
					return "", fmt.Errorf("email already exists")
				case "users_phone_key":
					log.Error().Str("phone", *user.Phone).Msg("phone already exists")
					return "", fmt.Errorf("phone already exists")
				default:
					log.Error().Err(err).Msg("unique constraint violation")
				}
			}
		}
		log.Err(err).Msg("failed to create user")
		return "", err

	}

	return id, nil
}

func (r *AuthRepo) CheckUserPassword(ctx context.Context, email, password string) (*models.User, error) {
	log := r.mylog.With().Str("repo", "CheckUserPassword").Logger()

	var user models.User
	q := `SELECT 
					user_id, 
					role, 
					firstname, 
					lastname, 
					email,
					phone
				FROM users 
				WHERE email = $1 AND password_hash = crypt($2, password_hash)`
	err := r.DB.conn.QueryRow(ctx, q, email, password).Scan(
		&user.UserID,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
	)
	if err != nil {
		log.Err(err).Msg("failed to check user password")
		return nil, err
	}
	log.Info().Str("user_id", user.UserID).Msg("user authenticated successfully")
	return &user, nil
}
