package repository

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/filter"
	"database/sql"
	"golang.org/x/net/context"
)

type Authorization interface {
	Create(c context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetById(c context.Context, id int) (entity.User, error)
}

type Equipment interface {
	Create(c context.Context, equipment *entity.Equipment) error
	Get(c context.Context, id int) (entity.Equipment, error)
	GetAll(c context.Context, options filter.OptionsMap) ([]entity.Equipment, error)
	Rent(c context.Context, agreementId int, itemIDs entity.EquipmentRequest) error
	GetAvailable(c context.Context, itemId int) ([]entity.AvailableProducts, error)
	GetRentedItems(c context.Context, userId int) ([]entity.EquipmentResponse, error)
	Update(c context.Context, item *entity.Equipment) error
	Delete(c context.Context, itemId int) error
	GetAvailableDates(c context.Context, request entity.AvailableDatesRequest) ([]entity.AvailableDatesResponse, error)
	AddProductToLocation(c context.Context, info entity.CreateAvailableProductsInput) error
	GetLogs(c context.Context) ([]entity.Logs, error)
}

type Agreement interface {
	Create(c context.Context, userId int) (entity.Agreement, error)
	Delete(c context.Context, agreementId int) error
	GetAll(c context.Context, userId int) ([]entity.AgreementResponse, error)
}

type Payment interface {
	Create(c context.Context, payment *entity.Payment) error
}

type Profile interface {
	Update(c context.Context, user *entity.User) error
	Get(c context.Context, userId int) (*entity.User, error)
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

type User interface {
	Create(c context.Context, user *entity.User) error
	Get(c context.Context, userId int) (*entity.User, error)
	UpdateWithRole(context context.Context, user *entity.User) error
	GiveRole(context context.Context, user *entity.User) error
	GetAll(c context.Context) ([]entity.User, error)
	Delete(c context.Context, userId int) error
	GetAllAgreementsById(c context.Context, agreementId int) ([]entity.AgreementResponse, error)
}

type Dump interface {
	InsertDump(c context.Context, filePath string) error
	GetAllDumps(c context.Context) ([]entity.Dump, error)
}

type Repository struct {
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

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Agreement:     NewAgreementPostgres(db),
		Equipment:     NewEquipmentPostgres(db),
		Payment:       NewPaymentPostgres(db),
		Profile:       NewProfilePostgres(db),
		Review:        NewReviewPostgres(db),
		Location:      NewLocationPostgres(db),
		User:          NewUserRepository(db),
		Dump:          NewDumpPostgres(db),
	}
}
