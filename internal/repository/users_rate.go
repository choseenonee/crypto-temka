package repository

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type usersRate struct {
	db *sqlx.DB
}

func InitUsersRate(db *sqlx.DB) UsersRate {
	return usersRate{db: db}
}

func (u usersRate) Create(ctx context.Context, urc models.UserRateCreate, walletID int) (int, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	timestamp := time.Now()

	_, err = tx.ExecContext(ctx, `UPDATE wallets SET deposit = deposit - $2 WHERE wallets.id = $1`, walletID, urc.Deposit)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit, token) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		urc.UserID, urc.RateID, urc.Lock, timestamp, timestamp, urc.Deposit, urc.Token)
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

func (u usersRate) Get(ctx context.Context, id int) (models.UserRate, error) {
	row := u.db.QueryRowContext(ctx, `SELECT id, user_id, rate_id, lock, opened, deposit, 
       earned_pool, outcome_pool, token FROM users_rates WHERE id = $1;`, id)

	var userRate models.UserRate
	err := row.Scan(&userRate.ID, &userRate.UserID, &userRate.RateID, &userRate.Lock, &userRate.Opened, &userRate.Deposit,
		&userRate.EarnedPool, &userRate.OutcomePool, &userRate.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserRate{}, fmt.Errorf("no userRate found with id %v", id)
		}
		return models.UserRate{}, err
	}

	return userRate, nil
}

func (u usersRate) GetByUser(ctx context.Context, userID, page, perPage int) ([]models.UserRate, error) {
	rows, err := u.db.QueryContext(ctx, `SELECT id, user_id, rate_id, lock, opened, deposit, 
       earned_pool, outcome_pool, token FROM users_rates WHERE user_id = $3 OFFSET $1 LIMIT $2`,
		(page-1)*perPage, perPage, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no userRate were found, userID: %v, page: %v, perPage: %v", userID, page, perPage)
		}
		return nil, err
	}

	var userRates []models.UserRate
	for rows.Next() {
		var userRate models.UserRate
		err = rows.Scan(&userRate.ID, &userRate.UserID, &userRate.RateID, &userRate.Lock, &userRate.Opened, &userRate.Deposit,
			&userRate.EarnedPool, &userRate.OutcomePool, &userRate.Token)
		if err != nil {
			return nil, err
		}

		userRates = append(userRates, userRate)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userRates, nil
}

func (u usersRate) Claim(ctx context.Context, userRateID, amount, walletID int) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users_rates SET outcome_pool = outcome_pool - $2 WHERE id = $1;`,
		userRateID, amount)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE wallets SET deposit = deposit + $2 WHERE id = $1`, walletID, amount)
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

func (u usersRate) ClaimDeposit(ctx context.Context, userRateID, amount, walletID int) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, `SELECT lock FROM users_rates WHERE id = $1;`, userRateID)

	var lock time.Time
	err = row.Scan(&lock)
	if err != nil {
		return err
	}

	lock = utils.DateOnly(lock)
	today := utils.DateOnly(time.Now())
	if today.Before(lock) {
		return fmt.Errorf("can't claim deposit before lock date: %v", lock)
	}

	_, err = tx.ExecContext(ctx, `UPDATE users_rates SET deposit = users_rates.deposit - $2 WHERE id = $1;`,
		userRateID, amount)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE wallets SET deposit = deposit + $2 WHERE id = $1`, walletID, amount)
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
