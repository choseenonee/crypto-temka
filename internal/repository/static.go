package repository

import (
	"context"
	"crypto-temka/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type static struct {
	db *sqlx.DB
}

func InitStatic(db *sqlx.DB) Static {
	return static{db: db}
}

func (s static) CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error) {
	propertiesJSON, err := json.Marshal(rc.Properties)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO reviews (title, text, properties) VALUES ($1, $2, $3) RETURNING id`,
		rc.Title, rc.Text, propertiesJSON)
	if err != nil {
		return 0, err
	}

	var id int
	err = row.Scan(&id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return 0, err
	}

	return id, nil
}

func (s static) GetReviews(ctx context.Context, page, perPage int) ([]models.Review, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, title, text, properties FROM reviews OFFSET $1 LIMIT $2`,
		(page-1)*perPage, perPage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no reviews were found, page: %v, perPage: %v", page, perPage)
		}
		return nil, err
	}

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var propertiesRaw []byte
		err = rows.Scan(&review.ID, &review.Title, &review.Text, &propertiesRaw)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(propertiesRaw, &review.Properties)

		reviews = append(reviews, review)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s static) UpdateReview(ctx context.Context, r models.Review) error {
	propertiesJSON, err := json.Marshal(r.Properties)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE reviews
											SET title = $2, text = $3, properties = $4 
											WHERE id = $1`,
		r.ID, r.Title, r.Text, propertiesJSON)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return fmt.Errorf("update cancelled, rows affected = %v", rowsAffected)
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return nil
}

func (s static) DeleteReview(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return fmt.Errorf("update cancelled, rows affected = %v", rowsAffected)
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return nil
}
