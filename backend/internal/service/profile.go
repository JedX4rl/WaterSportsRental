package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"golang.org/x/net/context"
	"time"
)

type ProfileService struct {
	repo repository.Profile
}

func NewProfileService(repo repository.Profile) *ProfileService {
	return &ProfileService{repo: repo}
}

func (p ProfileService) Update(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return p.repo.Update(ctx, user)
}

func (p ProfileService) Get(c context.Context, userId int) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return p.repo.Get(ctx, userId)
}
