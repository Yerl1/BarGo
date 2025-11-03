package dto

type SignupUser struct {
	Role      *string `json:"role" validate:"required,oneof=user admin"`
	FirstName *string `json:"firstname" validate:"required,min=2,max=100"`
	LastName  *string `json:"lastname" validate:"required,min=2,max=100"`
	Email     *string `json:"email" validate:"required,email"`
	Phone     *string `json:"phone" validate:"required,min=10,max=12"`
	Password  *string `json:"password" validate:"required,min=10,max=100"`
}

type LoginUser struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=10,max=100"`
}

type Resp struct {
	Token string   `json:"token"` // JWT token
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	Role      string `json:"role"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
