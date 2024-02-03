package t

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtUserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
