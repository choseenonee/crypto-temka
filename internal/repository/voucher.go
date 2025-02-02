package repository

import (
	"context"
	"crypto-temka/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type voucherRepo struct {
	db *sqlx.DB
}

func InitVoucherRepo(db *sqlx.DB) Voucher {
	return &voucherRepo{
		db: db,
	}
}

func (v voucherRepo) GetVoucherByID(ctx context.Context, id string) (models.Voucher, error) {
	var voucher models.Voucher
	var propertiesRaw []byte
	err := v.db.QueryRowContext(ctx, `SELECT id, type, properties FROM vouchers WHERE id = $1`).Scan(
		&voucher.Id, &voucher.VoucherType, &propertiesRaw)
	if err != nil {
		return models.Voucher{}, fmt.Errorf("voucherRepo.GetVoucherByID - v.db.QueryRowContext: %v", err)
	}

	err = json.Unmarshal(propertiesRaw, &voucher.Properties)
	if err != nil {
		return models.Voucher{}, fmt.Errorf("voucherRepo.GetVoucherByID - json.Unmarshal: %v", err)
	}

	return voucher, nil
}

func (v voucherRepo) CreateVoucher(ctx context.Context, id, voucherType string, properties []byte) error {
	_, err := v.db.ExecContext(ctx, `INSERT INTO vouchers VALUES ($1, $2, $3)`, id, voucherType, properties)
	if err != nil {
		return fmt.Errorf("voucherRepo.CreateVoucher - v.db.ExecContext: %v", err)
	}

	return nil
}

func (v voucherRepo) GetAllVouchers(ctx context.Context, offset, limit int) ([]models.Voucher, error) {
	rows, err := v.db.QueryContext(ctx, `SELECT id, type, properties FROM vouchers OFFSET $1 LIMIT $2`,
		offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no vouchers were found, offset: %v, limit: %v", offset, limit)
		}
		return nil, err
	}

	var vouchers []models.Voucher
	for rows.Next() {
		var voucher models.Voucher
		var propertiesRaw []byte
		err = rows.Scan(&voucher.Id, &voucher.VoucherType, &propertiesRaw)
		if err != nil {
			return nil, fmt.Errorf("voucherRepo.GetAllVouchers - rows.Scan: %v", err)
		}

		err := json.Unmarshal(propertiesRaw, &voucher.Properties)
		if err != nil {
			return nil, fmt.Errorf("voucherRepo.GetAllVouchers - json.Unmarshal: %v", err)
		}

		vouchers = append(vouchers, voucher)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}

func (v voucherRepo) UpdateVoucher(ctx context.Context, id string, voucherType string, properties []byte) error {
	res, err := v.db.ExecContext(ctx, `UPDATE vouchers SET type = $2, properties = $3 WHERE id = $1;`,
		id, voucherType, properties)
	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("no voucherRepo found with id %v", id)
	}

	return nil
}

func (v voucherRepo) DeleteVoucher(ctx context.Context, id string) error {
	res, err := v.db.ExecContext(ctx, `DELETE FROM vouchers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("voucherRepo.DeleteVoucher - tx.ExecContext: %v", err)
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("no voucherRepo found with id %v", id)
	}

	return nil
}
