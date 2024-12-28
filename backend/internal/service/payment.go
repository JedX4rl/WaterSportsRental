package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"golang.org/x/net/context"
	"time"
)

type PaymentService struct {
	repo repository.Payment
}

func (p PaymentService) Create(c context.Context, payment *entity.Payment) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return p.repo.Create(ctx, payment)
}

func NewPaymentService(repo repository.Payment) *PaymentService {
	return &PaymentService{repo: repo}
}
