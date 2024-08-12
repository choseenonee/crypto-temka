package repository

import (
	"context"
	"crypto-temka/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type wallet struct {
	db *sqlx.DB
}

func (w wallet) GetByUser(ctx context.Context, userID int) ([]models.Wallet, error) {
	rows, err := w.db.QueryContext(ctx, `SELECT id, user_id, token, deposit FROM wallets WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no wallets were found by userID: %v", userID)
		}
		return nil, err
	}

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		err = rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit)
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

func (w wallet) GetByToken(ctx context.Context, userID int, token string) (models.Wallet, error) {
	row := w.db.QueryRowContext(ctx, `SELECT id, user_id, token, deposit FROM wallets WHERE user_id = $1 AND token = $2`,
		userID, token)

	var wallet models.Wallet
	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Wallet{}, fmt.Errorf("no wallet were found, userID: %v, token: %v", userID, token)
		}
		return models.Wallet{}, err
	}

	return wallet, nil
}

func InitWallet(db *sqlx.DB) Wallet {
	return wallet{db: db}
}
