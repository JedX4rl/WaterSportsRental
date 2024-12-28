package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"golang.org/x/net/context"
	"time"
)

type AgreementService struct {
	repo repository.Agreement
}

func (a AgreementService) Create(c context.Context, userId int) (entity.Agreement, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return a.repo.Create(ctx, userId)
}

func (a AgreementService) Delete(c context.Context, agreementId int) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return a.repo.Delete(ctx, agreementId)
}

func (a AgreementService) GetAll(c context.Context, userId int) ([]entity.AgreementResponse, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return a.repo.GetAll(ctx, userId)
}

func NewAgreementService(repo repository.Agreement) *AgreementService {
	return &AgreementService{repo: repo}
}
