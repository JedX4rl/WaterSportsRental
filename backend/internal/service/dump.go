package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/repository"
	"context"
	"time"
)

type DumpService struct {
	repo repository.Dump
}

func (d DumpService) InsertDump(c context.Context, filePath string) error {
	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()
	return d.repo.InsertDump(ctx, filePath)
}

func (d DumpService) GetAllDumps(c context.Context) ([]entity.Dump, error) {
	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()
	return d.repo.GetAllDumps(ctx)
}

func NewDumpService(repo repository.Dump) *DumpService {
	return &DumpService{
		repo: repo,
	}
}
