package accessToken

import (
	"github.com/golang-jwt/jwt/v4"
)

type Jwt struct {
	Email string `json:"name"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

type JwtRefreshClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}
