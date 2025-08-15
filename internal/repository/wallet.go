package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

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
	row := w.db.QueryRowContext(ctx, `SELECT id, user_id, token, deposit, is_outcome, outcome FROM wallets WHERE id = $1`, walletID)

	var wallet models.Wallet
	if err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit, &wallet.IsOutcome, &wallet.Outcome); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Wallet{}, fmt.Errorf("no wallet were found by id: %v", walletID)
		}

		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w wallet) Insert(ctx context.Context, userID int, token string, amount float64, isOutcome bool) error {
	tx, err := w.db.Begin()
	if err != nil {
		return errors.Wrap(err, "could not start transaction")
	}

	rows := w.db.QueryRowContext(ctx, `UPDATE wallets SET deposit = deposit + $2 
               WHERE user_id = $1 AND token = $3 AND is_outcome = $4 RETURNING id`, userID, amount, token, isOutcome)

	var walletID sql.NullInt64
	err = rows.Scan(&walletID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return errors.Errorf("could not exec query: %v, also could not rollback transaction: %v", err, rbErr)
			}

			return err
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		rows := w.db.QueryRowContext(ctx, `INSERT INTO wallets (user_id, token, deposit, outcome, is_outcome) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			userID, token, amount, 0, isOutcome)

		err = rows.Scan(&walletID)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return errors.Errorf("could not exec query: %v, also could not rollback transaction: %v", err, rbErr)
			}

			return err
		}
	}

	insertIntoWalletsHistoryQuery := `INSERT INTO wallets_insert_history (wallet_id, ts, amount) VALUES ($1, current_timestamp, $2)`

	_, err = tx.ExecContext(ctx, insertIntoWalletsHistoryQuery, walletID.Int64, amount)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return errors.Errorf("could not exec query: %v, also could not rollback transaction: %v", err, rbErr)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "could not commit transaction")
	}

	return nil
}

func (w wallet) GetByUser(ctx context.Context, userID int) ([]models.Wallet, error) {
	rows, err := w.db.QueryContext(ctx, `SELECT id, user_id, token, deposit, is_outcome, outcome FROM wallets WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no wallets were found by userID: %v", userID)
		}
		return nil, err
	}

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		err = rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Token, &wallet.Deposit, &wallet.IsOutcome, &wallet.Outcome)
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

func (w wallet) Update(ctx context.Context, wallet models.WalletUpdate) error {
	query := `UPDATE wallets SET token = $2, deposit = $3, outcome = $4, is_outcome = $5 WHERE id = $1;`

	res, err := w.db.ExecContext(ctx, query, wallet.ID, wallet.Token, wallet.Deposit, wallet.Outcome, wallet.IsOutcome)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("no wallet found by id: %v", wallet.ID)
	}

	return nil
}

func (w wallet) GetWalletsInsertHistoryByUserID(ctx context.Context, userID int) ([]models.WalletInsertHistory, error) {
	rows, err := w.db.QueryContext(ctx, `SELECT wih.id, wih.wallet_id, wih.amount, wih.ts FROM wallets_insert_history wih JOIN wallets w ON wih.wallet_id = w.id WHERE w.user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no wallet insert history were found by userID: %v", userID)
		}

		return nil, err
	}

	var wallets []models.WalletInsertHistory
	for rows.Next() {
		var wallet models.WalletInsertHistory
		err = rows.Scan(&wallet.ID, &wallet.WalletID, &wallet.Amount, &wallet.TimeStamp)
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
