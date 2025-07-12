package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClaimTokenDto struct {
	Email  string `json:"email"`
	UserId int64  `json:"user_id"`
	jwt.RegisteredClaims
}
