package repository

import (
	"context"
	"crypto-temka/internal/models"
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// todo: test jsonb convert
	res, err := tx.ExecContext(ctx, `INSERT INTO reviews (tittle, text, properties) VALUES ($1, $2, $3) RETURNING id`,
		rc.Tittle, rc.Text, rc.Tittle)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return 0, err
	}

	return int(id), nil
}

func (s static) GetReviews(ctx context.Context, page, reviewsPerPage int) ([]models.Review, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, tittle, text, properties FROM reviews OFFSET $1 LIMIT $2`,
		(page-1)*reviewsPerPage, reviewsPerPage)
	if err != nil {
		return nil, err
	}

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		// todo: test jsonb convert
		err = rows.Scan(&review.ID, &review.Tittle, &review.Text, &review.Properties)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s static) UpdateReview(ctx context.Context, r models.Review) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// todo: test jsonb convert
	_, err = tx.ExecContext(ctx, `UPDATE reviews
											SET tittle = $2, text = $3, properties = $4 
											WHERE id = $1`,
		r.ID, r.Tittle, r.Text, r.Tittle)
	if err != nil {
		return err
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

	_, err = tx.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	if err != nil {
		return err
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
