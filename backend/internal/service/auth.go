package service

import (
	"WaterSportsRental/internal/entity"
	accessToken "WaterSportsRental/internal/jwt"
	"WaterSportsRental/internal/repository"
	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/net/context"
	"strconv"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

func (a AuthService) CreateAccessToken(user *entity.User, secret string, expiry int) (accToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &accessToken.Jwt{
		Email: user.Email,
		Id:    strconv.Itoa(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func (a AuthService) Create(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, 6*time.Second) //TODO: change
	defer cancel()
	return a.repo.Create(ctx, user)
}

func (a AuthService) GetByEmail(c context.Context, email string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()
	return a.repo.GetByEmail(ctx, email)
}

func (a AuthService) GetById(c context.Context, id int) (entity.User, error) {
	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()
	return a.repo.GetById(ctx, id)
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
