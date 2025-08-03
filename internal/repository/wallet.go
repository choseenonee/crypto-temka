package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crypto-temka/internal/models"

	"github.com/jmoiron/sqlx"
)

type wallet struct {
	db *sqlx.DB
}

func InitWallet(db *sqlx.DB) Wallet {
	return wallet{db: db}
}

func (w wallet) Get(ctx context.Context, walletID int) (models.Wallet, error) {
	row := w.db.QueryRowContext(ctx, `SELECT id, user_id, token, deposit, is_outcome FROM wallets WHERE id = $1`, walletID)

	var wallet models.Wallet
	if err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit, &wallet.IsOutcome); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Wallet{}, fmt.Errorf("no wallet were found by id: %v", walletID)
		}

		return models.Wallet{}, err
	}

	return wallet, nil
}

// TODO: pass is_outcome flag
// TODO: if there is no wallet found, then create one
func (w wallet) Insert(ctx context.Context, userID int, token string, amount float64, isOutcome bool) error {
	res, err := w.db.ExecContext(ctx, `UPDATE wallets SET deposit = deposit + $2 
               WHERE user_id = $1 AND token = $3 AND is_outcome = $4`, userID, amount, token, isOutcome)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no wallets were found by userID: %v and token: %v", userID, token)
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err := w.db.ExecContext(ctx, `INSERT INTO wallets (user_id, token, deposit, outcome, is_outcome) VALUES ($1, $2, $3, $4, $5)`,
			userID, token, amount, 0, isOutcome)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w wallet) GetByUser(ctx context.Context, userID int) ([]models.Wallet, error) {
	rows, err := w.db.QueryContext(ctx, `SELECT id, user_id, token, deposit, is_outcome FROM wallets WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no wallets were found by userID: %v", userID)
		}
		return nil, err
	}

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		err = rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit, &wallet.IsOutcome)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

// deprecated
//func (w wallet) GetByToken(ctx context.Context, userID int, token string) (models.Wallet, error) {
//	row := w.db.QueryRowContext(ctx, `SELECT id, user_id, token, deposit, is_outcome FROM wallets WHERE user_id = $1 AND token = $2`,
//		userID, token)
//
//	var wallet models.Wallet
//	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit, &wallet.IsOutcome)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return models.Wallet{}, fmt.Errorf("no wallet were found, userID: %v, token: %v", userID, token)
//		}
//		return models.Wallet{}, err
//	}
//
//	return wallet, nil
//}
