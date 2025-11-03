package models

import "time"

type User struct {
	UserID           string    `json:"user_id"`
	Role             *string   `json:"role" validate:"required,oneof=user admin"`
	FirstName        *string   `json:"firstname" validate:"required,min=2,max=100"`
	LastName         *string   `json:"lastname" validate:"required,min=2,max=100"`
	Email            *string   `json:"email" validate:"required, email"`
	Phone            *string   `json:"phone" validate:"required"`
	PasswordHash     *string   `json:"password_hash" validate:"required"`
	SubscriptionPlan *string   `json:"subscription_plan" validate:"required,oneof=none basic premium"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
