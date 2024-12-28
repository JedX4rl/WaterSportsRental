package repository

import (
	"WaterSportsRental/internal/entity"
	"context"
	"database/sql"
)

type ReviewPostgres struct {
	db *sql.DB
}

func (r ReviewPostgres) Create(c context.Context, review *entity.Review) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT create_review($1, $2, $3, $4, $5, $6)"
	row := tx.QueryRowContext(c, query, review.UserId, review.ItemId, review.Name, review.Rating, review.Comment, review.ReviewDate)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r ReviewPostgres) GetAll(c context.Context) ([]entity.Review, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	var reviews []entity.Review

	query := "SELECT * FROM get_all_reviews()"
	rows, err := tx.QueryContext(c, query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		var review entity.Review
		if err = rows.Scan(&review.Id, &review.Name, &review.Rating, &review.Comment, &review.ReviewDate, &review.ItemId); err != nil {
			tx.Rollback()
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, tx.Commit()
}

func (r ReviewPostgres) GetByItemID(c context.Context, itemId int) ([]entity.Review, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_reviews_by_item_id($1)"
	rows, err := tx.QueryContext(c, query, itemId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var reviews []entity.Review
	for rows.Next() {
		var review entity.Review
		if err = rows.Scan(&review.Id, &review.Name, &review.Rating, &review.Comment, &review.ReviewDate); err != nil {
			tx.Rollback()
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, tx.Commit()
}

func (r ReviewPostgres) GetByUserID(c context.Context, userId int) ([]entity.Review, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	query := "SELECT * FROM get_reviews_by_user_id($1)"
	rows, err := tx.QueryContext(c, query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var reviews []entity.Review
	for rows.Next() {
		var review entity.Review
		if err = rows.Scan(&review.Id, &review.Name, &review.Rating, &review.Comment, &review.ReviewDate); err != nil {
			tx.Rollback()
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, tx.Commit()
}

func (r ReviewPostgres) DeleteByID(c context.Context, reviewId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := "SELECT delete_review_by_id($1)"
	row := tx.QueryRowContext(c, query, reviewId)
	if err = row.Err(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewReviewPostgres(db *sql.DB) *ReviewPostgres {
	return &ReviewPostgres{
		db: db,
	}
}
