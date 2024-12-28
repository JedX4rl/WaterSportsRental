package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"context"
	"time"
)

type ReviewService struct {
	repo repository.Review
}

func (r ReviewService) Create(c context.Context, review *entity.Review) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return r.repo.Create(ctx, review)
}

func (r ReviewService) GetAll(c context.Context) ([]entity.Review, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return r.repo.GetAll(ctx)
}

func (r ReviewService) GetByItemID(c context.Context, itemId int) ([]entity.Review, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return r.repo.GetByItemID(ctx, itemId)
}

func (r ReviewService) DeleteByID(c context.Context, reviewId int) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return r.repo.DeleteByID(ctx, reviewId)
}

func (r ReviewService) GetByUserID(ctx context.Context, userId int) ([]entity.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return r.repo.GetByUserID(ctx, userId)
}

func NewReviewService(repo repository.Review) *ReviewService {
	return &ReviewService{repo: repo}
}
