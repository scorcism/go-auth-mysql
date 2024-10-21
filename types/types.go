package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type SecretsStore interface {
	GetSecrets() ([]Secret, error)
	AddSecret(Secret) error
}

type Secret struct {
	ID        int       `json:"id"`
	SecretKey string    `json:"secret_key"`
	Label     string    `json:"label"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"createdAt"`
}

type AddSecretPayload struct {
	Label     string `json:"label" validate:"required,min=3"`
	SecretKey string `json:"secret_key" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
