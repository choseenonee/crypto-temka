package user

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func InitUserRepo(db *sqlx.DB) repository.UserRepository {
	return UserRepo{db: db}
}

func (repo UserRepo) CreateUser(ctx context.Context, user models.UserCreate) (int, error) {
	var id int
	transaction, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	row := transaction.QueryRowContext(ctx, `insert into users (email, hashed_pwd, id_ref, amount) values ($1, $2, $3, $4) returning id`, user.Email, user.Password, user.Ref, user.Amount)
	err = row.Scan(&id)
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return 0, rbErr
		}
		return 0, err
	}
	if err = transaction.Commit(); err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return 0, rbErr
		}
		return 0, err
	}
	return id, nil
}

func (repo UserRepo) GetHashPWD(ctx context.Context, email string) (int, string, error) {
	var pwd string
	var id int
	rows := repo.db.QueryRowContext(ctx, `SELECT id, hashed_pwd from users WHERE email = $1;`, email)
	err := rows.Scan(&id, &pwd)
	if err != nil {
		return 0, "", err
	}
	return id, pwd, nil
}

func (repo UserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	rows := repo.db.QueryRowContext(ctx, `SELECT email, id_ref, amount, status from users WHERE id = $1;`, id)
	err := rows.Scan(&user.Email, &user.Ref, &user.Amount, &user.Status)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return &user, nil
}
func (repo UserRepo) AddPhoto(ctx context.Context, photo models.UserPhoto) error {
	transaction, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	result, err := transaction.ExecContext(ctx, `UPDATE users SET photo=$2 WHERE id=$1;`, photo.ID, photo.Photo)
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if count != 1 {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return fmt.Errorf("no one string")
	}

	if err = transaction.Commit(); err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return nil
}

func (repo UserRepo) DeleteUser(ctx context.Context, id int) error {
	transaction, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	result, err := transaction.ExecContext(ctx, `DELETE FROM users WHERE id=$1;`, id)
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	if count != 1 {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	if err = transaction.Commit(); err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return nil
}

func (repo UserRepo) GetPhoto(ctx context.Context, idUser int) (*models.UserPhoto, error) {
	var user models.UserPhoto
	rows := repo.db.QueryRowContext(ctx, `SELECT photo, status from users WHERE id = $1;`, idUser)
	err := rows.Scan(&user.Photo, &user.Status)
	if err != nil {
		return nil, err
	}
	user.ID = idUser
	return &user, nil
}

func (repo UserRepo) SetStatus(ctx context.Context, status models.SetStatus) error {
	transaction, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	result, err := transaction.ExecContext(ctx, `UPDATE users SET status=$2 WHERE id=$1;`, status.ID, status.Status)
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if count != 1 {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return fmt.Errorf("no one string")
	}

	if err = transaction.Commit(); err != nil {
		if rbErr := transaction.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return nil
}
