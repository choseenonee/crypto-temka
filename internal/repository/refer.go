package repository

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/pkg/config"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type refer struct {
	db *sqlx.DB
}

func InitRefer(db *sqlx.DB) Refer {
	return refer{db: db}
}

func (r refer) Get(ctx context.Context, userID, page, perPage int) ([]models.Refer, error) {
	query := `SELECT id, parent_id, child_id, amount, token FROM refers WHERE parent_id = $1 OFFSET $2 LIMIT $3;`
	rows, err := r.db.QueryContext(ctx, query, userID, (page-1)*perPage, perPage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no refers were found, userID: %v, page: %v, perPage: %v", userID, page, perPage)
		}
		return nil, err
	}

	var refers []models.Refer
	for rows.Next() {
		var refer models.Refer
		err = rows.Scan(&refer.ID, &refer.ParentID, &refer.ChildID, &refer.Amount, &refer.Token)
		if err != nil {
			return nil, err
		}

		refers = append(refers, refer)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return refers, nil
}

func (r refer) Claim(ctx context.Context, id, userID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `SELECT amount, token FROM refers WHERE id = $1 AND parent_id = $2 FOR UPDATE;`

	row := tx.QueryRowContext(ctx, query, id, userID)

	var amount float64
	var token string
	err = row.Scan(&amount, &token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("no refer were found, id: %v, userID %v", id, userID)
		}
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE refers SET amount = amount - $1 WHERE id = $2`, amount, id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	amount = amount * viper.GetFloat64(config.ReferPercent) / 100

	res, err := tx.ExecContext(ctx, `UPDATE wallets SET deposit = deposit + $1, outcome = outcome + $1 WHERE user_id = $2 AND token = $3`,
		amount, userID, token)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		_, err = tx.ExecContext(ctx, `INSERT INTO wallets (user_id, token, deposit, outcome) VALUES ($1, $2, $3, $3)`,
			userID, token, amount)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
			}
			return err
		}
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
