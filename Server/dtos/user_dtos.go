package dtos

import (
	"github.com/BouhairieMcKnight/MovieStream/Server/domain"
)

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UserId         string         `json:"user_id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	Role           string         `json:"role"`
	RefreshToken   string         `json:"refresh_token"`
	Token          string         `json:"token"`
	FavoriteGenres []domain.Genre `json:"favorite_genres"`
}

type UserLogout struct {
	UserId string `json:"user_id"`
}
