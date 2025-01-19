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

type user struct {
	db *sqlx.DB
}

func InitUser(db *sqlx.DB) User {
	return user{db: db}
}

func (u user) Create(ctx context.Context, uc models.UserCreate) (int, error) {
	propertiesJSON, err := json.Marshal(uc.Properties)
	if err != nil {
		return 0, err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO users (email, phone_number, hashed_password, status, refer_id, properties)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		uc.Email, uc.PhoneNumber, uc.Password, "opened", uc.ReferID, propertiesJSON)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
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

func (u user) processUser(ctx context.Context, query string, value interface{}) (models.User, error) {
	row := u.db.QueryRowContext(ctx, query, value)

	var user models.User
	var propertiesRaw []byte
	err := row.Scan(&user.ID, &user.Email, &user.PhoneNumber, &user.Status, &user.ReferID, &propertiesRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("no user with %v", value)
		}
		return models.User{}, err
	}

	_ = json.Unmarshal(propertiesRaw, &user.Properties)

	return user, nil
}

func (u user) GetByID(ctx context.Context, id int) (models.User, error) {
	query := `SELECT id, email, phone_number, status, refer_id, properties FROM users WHERE id = $1;`
	return u.processUser(ctx, query, id)
}

func (u user) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT id, email, phone_number, status, refer_id, properties FROM users WHERE email = $1;`
	return u.processUser(ctx, query, email)
}

func (u user) GetHashedPwd(ctx context.Context, email string) (string, error) {
	row := u.db.QueryRowContext(ctx, `SELECT hashed_password FROM users WHERE email = $1;`, email)

	var hashedPwd string
	err := row.Scan(&hashedPwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no user with email %v", email)
		}
		return "", err
	}

	return hashedPwd, nil
}

func (u user) GetStatus(id int) (string, error) {
	row := u.db.QueryRow(`SELECT status FROM users WHERE id = $1;`, id)

	var status string
	err := row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no user with id %v", id)
		}
		return "", err
	}

	return status, nil
}

func (u user) GetAll(ctx context.Context, page, perPage int, status string) ([]models.User, error) {
	var err error
	var rows *sql.Rows
	if status != "" {
		query := `SELECT id, email, phone_number, status, refer_id, properties FROM users WHERE status = $1 OFFSET $2 LIMIT $3;`
		rows, err = u.db.QueryContext(ctx, query, status, (page-1)*perPage, perPage)
	} else {
		query := `SELECT id, email, phone_number, status, refer_id, properties FROM users OFFSET $1 LIMIT $2;`
		rows, err = u.db.QueryContext(ctx, query, (page-1)*perPage, perPage)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no users with, status: %v", status)
		}
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		var propertiesRaw []byte
		err = rows.Scan(&user.ID, &user.Email, &user.PhoneNumber, &user.Status, &user.ReferID, &propertiesRaw)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(propertiesRaw, &user.Properties)

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u user) UpdateStatus(ctx context.Context, id int, status string) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE users SET status = $2 WHERE id = $1 AND status = 'opened';`, id, status)
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
			return fmt.Errorf("err: %v, rbErr: %v", "no user found with status opened", rbErr)
		}
		return fmt.Errorf("no user found with status opened")
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

func (u user) UpdateProperties(ctx context.Context, id int, properties interface{}) error {
	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE users SET properties = $2 WHERE id = $1;`, id, propertiesJSON)
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
			return fmt.Errorf("err: %v, rbErr: %v", "no user found by id", rbErr)
		}
		return fmt.Errorf("no user found by id")
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
