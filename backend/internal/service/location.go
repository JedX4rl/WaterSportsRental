package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"context"
	"time"
)

type LocationService struct {
	repo repository.Location
}

func (l LocationService) Create(c context.Context, location *entity.Location) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return l.repo.Create(ctx, location)
}

func (l LocationService) GetAll(c context.Context) ([]entity.Location, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return l.repo.GetAll(ctx)
}

func (l LocationService) Update(c context.Context, location *entity.Location) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return l.repo.Update(ctx, location)
}

func (l LocationService) Delete(c context.Context, locationId int) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return l.repo.Delete(ctx, locationId)
}

func (l LocationService) DeleteAll(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return l.repo.DeleteAll(ctx)
}

func NewLocationService(repo repository.Location) *LocationService {
	return &LocationService{
		repo: repo,
	}
}
