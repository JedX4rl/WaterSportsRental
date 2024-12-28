package service

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/filter"
	"WaterSportsRental/internal/repository"
	"golang.org/x/net/context"
)

type Authorization interface {
	Create(c context.Context, user *entity.User) error
	GetByEmail(c context.Context, email string) (entity.User, error)
	GetById(c context.Context, id int) (entity.User, error)
	CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error)
}

type Agreement interface {
	Create(c context.Context, userId int) (entity.Agreement, error)
	Delete(c context.Context, agreementId int) error
	GetAll(c context.Context, userId int) ([]entity.AgreementResponse, error)
}

type Payment interface {
	Create(c context.Context, payment *entity.Payment) error
}

type Equipment interface {
	Create(c context.Context, equipment *entity.Equipment) error
	Get(c context.Context, id int) (entity.Equipment, error)
	Rent(c context.Context, agreementId int, itemIDs entity.EquipmentRequest) error
	GetAll(c context.Context, options filter.OptionsMap) ([]entity.Equipment, error)
	GetAvailable(c context.Context, itemId int) ([]entity.AvailableProducts, error)
	GetRentedItems(c context.Context, userId int) ([]entity.EquipmentResponse, error)
	Update(c context.Context, item *entity.Equipment) error
	Delete(c context.Context, itemId int) error
	GetAvailableDates(c context.Context, request entity.AvailableDatesRequest) ([]entity.AvailableDatesResponse, error)
	AddProductToLocation(c context.Context, info entity.CreateAvailableProductsInput) error
	GetLogs(c context.Context) ([]entity.Logs, error)
}

type Profile interface {
	Update(context context.Context, user *entity.User) error
	Get(c context.Context, userId int) (*entity.User, error)
}

type User interface {
	Create(c context.Context, user *entity.User) error
	Get(c context.Context, userId int) (*entity.User, error)
	UpdateWithRole(context context.Context, user *entity.User) error
	GiveRole(context context.Context, user *entity.User) error
	GetAll(c context.Context) ([]entity.User, error)
	Delete(c context.Context, userId int) error
	GetAllAgreementsById(c context.Context, agreementId int) ([]entity.AgreementResponse, error)
}

type Review interface {
	Create(c context.Context, review *entity.Review) error
	GetAll(c context.Context) ([]entity.Review, error)
	GetByItemID(c context.Context, itemId int) ([]entity.Review, error)
	GetByUserID(ctx context.Context, userId int) ([]entity.Review, error)
	DeleteByID(ctx context.Context, reviewId int) error
}

type Location interface {
	Create(c context.Context, location *entity.Location) error
	GetAll(c context.Context) ([]entity.Location, error)
	Update(c context.Context, location *entity.Location) error
	Delete(c context.Context, locationId int) error
	DeleteAll(c context.Context) error
}

type Dump interface {
	InsertDump(c context.Context, filePath string) error
	GetAllDumps(c context.Context) ([]entity.Dump, error)
}

type Service struct {
	Authorization
	Agreement
	Equipment
	Payment
	Profile
	Review
	Location
	User
	Dump
}

func NewServices(repositories *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repositories.Authorization),
		Equipment:     NewEquipmentService(repositories.Equipment),
		Agreement:     NewAgreementService(repositories.Agreement),
		Payment:       NewPaymentService(repositories.Payment),
		Profile:       NewProfileService(repositories.Profile),
		Review:        NewReviewService(repositories.Review),
		Location:      NewLocationService(repositories.Location),
		User:          NewUserService(repositories.User),
		Dump:          NewDumpService(repositories.Dump),
	}
}
