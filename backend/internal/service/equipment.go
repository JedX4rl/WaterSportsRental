package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/filter"
	"WaterSportsRental/internal/repository"
	"golang.org/x/net/context"
	"time"
)

type EquipmentService struct {
	repo repository.Equipment
}

func (e EquipmentService) Create(c context.Context, equipment *entity.Equipment) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.Create(ctx, equipment)
}

func (e EquipmentService) Get(c context.Context, id int) (entity.Equipment, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.Get(ctx, id)
}

func (e EquipmentService) Rent(c context.Context, agreementId int, itemIDs entity.EquipmentRequest) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.Rent(ctx, agreementId, itemIDs)
}

func (e EquipmentService) GetAll(c context.Context, options filter.OptionsMap) ([]entity.Equipment, error) {
	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()
	return e.repo.GetAll(ctx, options)
}

func (e EquipmentService) GetAvailable(c context.Context, itemId int) ([]entity.AvailableProducts, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.GetAvailable(ctx, itemId)
}

func (e EquipmentService) Update(c context.Context, item *entity.Equipment) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.Update(ctx, item)
}

func (e EquipmentService) Delete(c context.Context, id int) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.Delete(ctx, id)
}

func (e EquipmentService) GetRentedItems(c context.Context, userId int) ([]entity.EquipmentResponse, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.GetRentedItems(ctx, userId)
}

func (e EquipmentService) GetAvailableDates(c context.Context, request entity.AvailableDatesRequest) ([]entity.AvailableDatesResponse, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.GetAvailableDates(ctx, request)
}

func (e EquipmentService) AddProductToLocation(c context.Context, info entity.CreateAvailableProductsInput) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.AddProductToLocation(ctx, info)
}

func (e EquipmentService) GetLogs(c context.Context) ([]entity.Logs, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return e.repo.GetLogs(ctx)
}

func NewEquipmentService(repo repository.Equipment) *EquipmentService {
	return &EquipmentService{
		repo: repo,
	}
}
