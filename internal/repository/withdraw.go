package repository

import (
	"context"
	"crypto-temka/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

type withdraw struct {
	db *sqlx.DB
}

func InitWithdraw(db *sqlx.DB) Withdraw {
	return withdraw{db: db}
}

func (w withdraw) Create(ctx context.Context, wu models.WithdrawCreate, walletID int) (int, error) {
	propertiesRaw, err := json.Marshal(wu.Properties)
	if err != nil {
		return 0, err
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE wallets SET deposit = deposit - $2 WHERE id = $1;`,
		walletID, wu.Amount)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO withdrawals (user_id, amount, token, status, properties) 
			VALUES ($1, $2, $3, $4, $5) RETURNING id;`, wu.UserID, wu.Amount, wu.Token, "opened", propertiesRaw)

	var id int
	err = row.Scan(&id)
	if err != nil {
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

	return id, err
}

func (w withdraw) GetByUserID(ctx context.Context, page, perPage, userID int) ([]models.Withdraw, error) {
	rows, err := w.db.QueryContext(ctx, `SELECT id, user_id, amount, token, status, properties FROM withdrawals 
                                                      WHERE user_id = $3 OFFSET $1 LIMIT $2;`,
		(page-1)*perPage, perPage, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no withdrawals were found, userID: %v, page: %v, perPage: %v", userID, page, perPage)
		}
		return nil, err
	}

	var withdrawals []models.Withdraw
	for rows.Next() {
		var withdraw models.Withdraw
		var propertiesRaw []byte
		err = rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Amount, &withdraw.Token, &withdraw.Status, &propertiesRaw)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(propertiesRaw, &withdraw.Properties)

		withdrawals = append(withdrawals, withdraw)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (w withdraw) GetByID(ctx context.Context, withdrawID int) (models.Withdraw, error) {
	row := w.db.QueryRowContext(ctx, `SELECT id, user_id, amount, token, status, properties FROM withdrawals 
                                                      WHERE id = $1;`, withdrawID)

	var withdraw models.Withdraw
	var propertiesRaw []byte
	err := row.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Amount, &withdraw.Token, &withdraw.Status, &propertiesRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Withdraw{}, fmt.Errorf("no withdrawals were found, withdrawID: %v", withdrawID)
		}
		return models.Withdraw{}, err
	}

	_ = json.Unmarshal(propertiesRaw, &withdraw.Properties)

	err = row.Err()
	if err != nil {
		return models.Withdraw{}, err
	}

	return withdraw, nil
}

func (w withdraw) GetAll(ctx context.Context, page, perPage int, status string) ([]models.Withdraw, error) {
	var err error
	var rows *sql.Rows
	if status != "" {
		query := `SELECT id, user_id, amount, token, status, properties FROM withdrawals WHERE status = $1 OFFSET $2 LIMIT $3;`
		rows, err = w.db.QueryContext(ctx, query, status, (page-1)*perPage, perPage)
	} else {
		query := `SELECT id, user_id, amount, token, status, properties FROM withdrawals OFFSET $1 LIMIT $2;`
		rows, err = w.db.QueryContext(ctx, query, (page-1)*perPage, perPage)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no withdrawals were found, status: %v, page: %v, perPage: %v", status, page, perPage)
		}
		return nil, err
	}

	var withdrawals []models.Withdraw
	for rows.Next() {
		var withdraw models.Withdraw
		var propertiesRaw []byte
		err = rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Amount, &withdraw.Token, &withdraw.Status, &propertiesRaw)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func processRefer(tx *sql.Tx, ctx context.Context, withdrawID int) error {
	row := tx.QueryRowContext(ctx, `SELECT u.id, u.refer_id, w.token, w.amount  FROM withdrawals w JOIN users u on u.id = w.user_id 
                  WHERE w.id = $1;`, withdrawID)

	var childID null.Int
	var parentID null.Int
	var withdrawalToken null.String
	var withdrawalAmount null.Int
	err := row.Scan(&childID, &parentID, &withdrawalToken, &withdrawalAmount)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("error while selecting refer data err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	if parentID.Valid {
		res, err := tx.ExecContext(ctx, `UPDATE refers SET amount = amount + $1 WHERE token = $2 AND child_id = $3`,
			withdrawalAmount.Int64, withdrawalToken.String, childID.Int64)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("error while updating refer data err: %v, rbErr: %v", err, rbErr)
			}
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("error while rows affected refer data err: %v, rbErr: %v", err, rbErr)
			}
			return err
		}

		if rowsAffected == 0 {
			_, err = tx.ExecContext(ctx, `INSERT INTO refers (parent_id, child_id, amount, token) VALUES ($1, $2, $3, $4)`,
				parentID.Int64, childID.Int64, withdrawalAmount.Int64, withdrawalToken.String)
			if err != nil {
				rbErr := tx.Rollback()
				if rbErr != nil {
					return fmt.Errorf("error while updating refer data err: %v, rbErr: %v", err, rbErr)
				}
				return err
			}
		}
	}

	return nil
}

func (w withdraw) UpdateStatus(ctx context.Context, withdrawID int, status string, properties interface{}) error {
	propertiesRaw, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE withdrawals SET status = $2, properties = $3 WHERE id = $1;`,
		withdrawID, status, propertiesRaw)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	if status == "verified" {
		err = processRefer(tx, ctx, withdrawID)
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
