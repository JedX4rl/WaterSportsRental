package repository

import (
	"WaterSportsRental/internal/entity"
	"WaterSportsRental/internal/filter"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"log/slog"
)

type EquipmentPostgres struct {
	db *sql.DB
}

func (e EquipmentPostgres) Create(c context.Context, equipment *entity.Equipment) error {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return err
	}

	query := "SELECT add_item($1, $2, $3, $4, $5, $6)"
	row := tx.QueryRowContext(c, query, equipment.Type, equipment.Brand, equipment.Model, equipment.Year, equipment.Price, equipment.Image)
	if err = row.Err(); err != nil {
		slog.Info(err.Error())
		tx.Rollback()
		return err
	}
	if err = row.Scan(&equipment.Id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}

func (e EquipmentPostgres) Get(ctx context.Context, id int) (entity.Equipment, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Equipment{}, err
	}
	query := "SELECT * FROM get_item($1)"
	row := tx.QueryRowContext(ctx, query, id)

	if err = row.Err(); err != nil {
		tx.Rollback()
		slog.Info(err.Error())
		return entity.Equipment{}, err
	}

	var equipment entity.Equipment
	var tYear sql.NullInt64
	var tType, tBrand, tModel, image sql.NullString
	var tPrice sql.NullFloat64

	if err = row.Scan(&equipment.Id, &tType, &tBrand, &tModel, &tYear, &tPrice, &image); err != nil {
		tx.Rollback()
		return entity.Equipment{}, err
	}

	equipment.Type = tType.String
	equipment.Brand = tBrand.String
	equipment.Model = tModel.String
	equipment.Year = int(tYear.Int64)
	equipment.Price = tPrice.Float64
	equipment.Image = image.String

	return equipment, tx.Commit()

}

func (e EquipmentPostgres) GetAll(ctx context.Context, options filter.OptionsMap) ([]entity.Equipment, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	fieldsToFilter := options.Fields()
	filters := make(map[string]interface{})

	for key, conditions := range fieldsToFilter {
		var values []interface{}
		for _, condition := range conditions {
			values = append(values, condition.Value)
		}
		filters[key] = values
	}

	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to serialize filters: %v", err)
	}

	var query string

	if len(filters) == 0 {
		query = "SELECT * FROM get_all_items($1)"
	} else {
		query = "SELECT * FROM get_filtered_items($1)"
	}

	rows, err := tx.QueryContext(ctx, query, filtersJSON)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	var equipments []entity.Equipment

	for rows.Next() {
		var equipment entity.Equipment
		var image sql.NullString

		err := rows.Scan(
			&equipment.Id,
			&equipment.Type,
			&equipment.Brand,
			&equipment.Model,
			&equipment.Year,
			&equipment.Price,
			&image,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		equipment.Image = image.String
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("rows iteration failed: %v", err)
	}

	return equipments, tx.Commit()
}

func (e EquipmentPostgres) Rent(ctx context.Context, agreementId int, request entity.EquipmentRequest) error {

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := "SELECT rent_equipment($1, $2, $3, $4, $5, $6)"

	row := tx.QueryRowContext(ctx, query, agreementId, pq.Array(request.IDs), request.LocationId, request.StartDate, request.EndDate, request.PaymentMethod)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (e EquipmentPostgres) GetAvailable(c context.Context, itemId int) ([]entity.AvailableProducts, error) {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_available_products($1)"
	rows, err := tx.QueryContext(c, query, itemId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var availableProducts []entity.AvailableProducts

	for rows.Next() {
		var currentProduct entity.AvailableProducts
		err := rows.Scan(&currentProduct.ProductId, &currentProduct.LocationId, &currentProduct.Location, &currentProduct.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		availableProducts = append(availableProducts, currentProduct)
	}
	return availableProducts, tx.Commit()
}

func (e EquipmentPostgres) GetAvailableDates(c context.Context, request entity.AvailableDatesRequest) ([]entity.AvailableDatesResponse, error) {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_available_periods_by_location_inventory($1, $2)"
	rows, err := tx.QueryContext(c, query, request.LocationId, request.ItemId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var availableDates []entity.AvailableDatesResponse
	for rows.Next() {
		var currentPeriod entity.AvailableDatesResponse
		if err = rows.Scan(&currentPeriod.StartDate, &currentPeriod.EndDate); err != nil {
			tx.Rollback()
			return nil, err
		}
		availableDates = append(availableDates, currentPeriod)
	}
	return availableDates, tx.Commit()
}

func (e EquipmentPostgres) Update(c context.Context, item *entity.Equipment) error {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT update_item($1, $2, $3, $4, $5, $6, $7)"
	row := tx.QueryRowContext(c, query, item.Id, item.Type, item.Brand, item.Model, item.Year, item.Price, item.Image)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (e EquipmentPostgres) Delete(c context.Context, itemId int) error {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT delete_item($1)"
	row := tx.QueryRowContext(c, query, itemId)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (e EquipmentPostgres) GetRentedItems(c context.Context, userId int) ([]entity.EquipmentResponse, error) {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_rented_items($1)"
	rows, err := tx.QueryContext(c, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var rentedItems []entity.EquipmentResponse
	for rows.Next() {
		var rentedItem entity.EquipmentResponse
		if err = rows.Scan(&rentedItem.Id, &rentedItem.Type, &rentedItem.Brand, &rentedItem.Model, &rentedItem.Year, &rentedItem.TotalPrice, &rentedItem.Image,
			&rentedItem.StartDate, &rentedItem.EndDate); err != nil {
			tx.Rollback()
			return nil, err
		}
		rentedItems = append(rentedItems, rentedItem)
	}
	return rentedItems, tx.Commit()
}

func (e EquipmentPostgres) AddProductToLocation(c context.Context, info entity.CreateAvailableProductsInput) error {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	query := "SELECT add_product_availability($1, $2, $3)"
	row := tx.QueryRowContext(c, query, info.ProductId, info.LocationId, info.Number)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (e EquipmentPostgres) GetLogs(c context.Context) ([]entity.Logs, error) {
	tx, err := e.db.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_all_rented_equipment_logs()"
	rows, err := tx.QueryContext(c, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var logs []entity.Logs
	for rows.Next() {
		var log entity.Logs
		if err = rows.Scan(&log.ID, &log.ItemId, &log.StartDate, &log.EndDate); err != nil {
			tx.Rollback()
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, tx.Commit()
}

func NewEquipmentPostgres(db *sql.DB) *EquipmentPostgres {
	return &EquipmentPostgres{
		db: db,
	}
}
