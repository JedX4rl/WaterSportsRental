package env

import (
	"log"
	"os"
	"strconv"
)

type Tokens struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  int64
	RefreshTokenExpiry int64
}

func NewTokens() Tokens {
	tokens := Tokens{}

	tokens.AccessToken = os.Getenv("ACCESS_TOKEN_SECRET")
	if tokens.AccessToken == "" {
		log.Fatalf("ACCESS_TOKEN_SECRET environment variable not set")
	}
	tokens.RefreshToken = os.Getenv("REFRESH_TOKEN_SECRET")
	if tokens.RefreshToken == "" {
		log.Fatalf("REFRESH_TOKEN_SECRET environment variable not set")
	}
	tokens.AccessTokenExpiry, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"), 10, 64)
	if tokens.AccessTokenExpiry == 0 {
		log.Fatalf("ACCESS_TOKEN_EXPIRY_HOUR environment variable not set")
	}
	tokens.RefreshTokenExpiry, _ = strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"), 10, 64)
	if tokens.AccessTokenExpiry == 0 {
		log.Fatalf("ACCESS_TOKEN_EXPIRY_HOUR environment variable not set")
	}
	return tokens
}
