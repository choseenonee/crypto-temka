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

type message struct {
	db *sqlx.DB
}

func InitMessage(db *sqlx.DB) Message {
	return message{db: db}
}

func (m message) Create(ctx context.Context, mc models.MessageCreate) (int, error) {
	propertiesRaw, err := json.Marshal(mc.Properties)
	if err != nil {
		return 0, err
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// todo: convert timestamp to utc
	row := tx.QueryRowContext(ctx, `INSERT INTO messages (user_id, properties, timestamp, is_read) VALUES ($1, $2, $3, false) 
                                                               RETURNING id`, mc.UserID, propertiesRaw, mc.Timestamp)

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

func (m message) GetByID(ctx context.Context, id, userID int) (models.Message, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Message{}, err
	}

	row := tx.QueryRowContext(ctx, `SELECT id, user_id, properties, timestamp FROM messages WHERE id = $1 AND 
                                                              user_id = $2 FOR UPDATE`, id, userID)

	var msg models.Message
	var propertiesRaw []byte
	err = row.Scan(&msg.ID, &msg.UserID, &propertiesRaw, &msg.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("no message were found, id: %v, userID: %v", id, userID)
		}
		rbErr := tx.Rollback()
		if rbErr != nil {
			return models.Message{}, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return models.Message{}, err
	}

	_ = json.Unmarshal(propertiesRaw, &msg.Properties)

	_, err = tx.ExecContext(ctx, `UPDATE messages SET is_read = true WHERE id = $1`, id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return models.Message{}, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return models.Message{}, err
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return models.Message{}, fmt.Errorf("err: %v, rbErr: %v", err, rbErr)
		}
		return models.Message{}, err
	}

	return msg, nil
}

func (m message) GetByUser(ctx context.Context, userID int) ([]models.Message, error) {
	rows, err := m.db.QueryContext(ctx, `SELECT id, user_id, properties, timestamp, is_read FROM messages WHERE user_id = $1`,
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("no messages were found, userID: %v", userID)
		}
		return nil, err
	}

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var propertiesRaw []byte
		err = rows.Scan(&msg.ID, &msg.UserID, &propertiesRaw, &msg.Timestamp, &msg.IsRead)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = fmt.Errorf("no messages were found, userID: %v", userID)
			}
			return nil, err
		}

		_ = json.Unmarshal(propertiesRaw, &msg.Properties)

		messages = append(messages, msg)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil
}
