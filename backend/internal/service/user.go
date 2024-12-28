package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"context"
	"time"
)

type UserService struct {
	repo repository.User
}

func (u UserService) Create(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.Create(ctx, user)
}

func (u UserService) UpdateWithRole(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.UpdateWithRole(ctx, user)
}

func (u UserService) GiveRole(c context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.UpdateWithRole(ctx, user)
}

func (u UserService) Get(c context.Context, userId int) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.Get(ctx, userId)
}

func (u UserService) GetAll(c context.Context) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.GetAll(ctx)
}

func (u UserService) Delete(c context.Context, userId int) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.Delete(ctx, userId)
}

func (u UserService) GetAllAgreementsById(c context.Context, agreementId int) ([]entity.AgreementResponse, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return u.repo.GetAllAgreementsById(ctx, agreementId)
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}
