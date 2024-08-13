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

type rate struct {
	db *sqlx.DB
}

func InitRate(db *sqlx.DB) Rate {
	return rate{db: db}
}

func (r rate) CreateRate(ctx context.Context, rc models.RateCreate) (int, error) {
	propertiesJSON, err := json.Marshal(rc.Properties)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO rates (title, profit, min_lock_days, commission, properties) 
VALUES ($1, $2, $3, $4) RETURNING id`,
		rc.Title, rc.Profit, rc.MinLockDays, rc.Commission, propertiesJSON)
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

func (r rate) GetRates(ctx context.Context, page, perPage int) ([]models.Rate, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, profit, min_lock_days, commission, properties 
FROM rates OFFSET $1 LIMIT $2`,
		(page-1)*perPage, perPage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no rates were found, page: %v, perPage: %v", page, perPage)
		}
		return nil, err
	}

	var rates []models.Rate
	for rows.Next() {
		var rate models.Rate
		var propertiesRaw []byte
		err = rows.Scan(&rate.ID, &rate.Title, &rate.Profit, &rate.MinLockDays, &rate.Commission, &propertiesRaw)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(propertiesRaw, &rate.Properties)

		rates = append(rates, rate)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (r rate) UpdateRate(ctx context.Context, ru models.Rate) error {
	propertiesJSON, err := json.Marshal(ru.Properties)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE rates
											SET title = $2, profit = $3, min_lock_days = $4, commission = $5, properties = $6
											WHERE id = $1`,
		ru.ID, ru.Title, ru.Profit, ru.MinLockDays, ru.Commission, propertiesJSON)
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

func (r rate) GetRate(ctx context.Context, id int) (models.Rate, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, profit, min_lock_days, properties FROM rates WHERE id = $1;`, id)

	var rate models.Rate
	var propertiesRaw []byte
	err := row.Scan(&rate.ID, &rate.Title, &rate.Profit, &rate.MinLockDays, &rate.Commission, &propertiesRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Rate{}, fmt.Errorf("rate not found with id %v", id)
		}
		return models.Rate{}, err
	}

	_ = json.Unmarshal(propertiesRaw, &rate.Properties)

	return rate, nil
}
