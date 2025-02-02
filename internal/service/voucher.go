package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"
)

type voucherService struct {
	logger *log.Logs
	repo   repository.Voucher
}

func InitVoucherService(repo repository.Voucher, logger *log.Logs) Voucher {
	return voucherService{
		logger: logger,
		repo:   repo,
	}
}

func (v voucherService) GetVoucherByID(ctx context.Context, id string) (models.Voucher, error) {
	voucher, err := v.repo.GetVoucherByID(ctx, id)
	if err != nil {
		v.logger.Error(err.Error())
		return voucher, err
	}

	return voucher, nil
}

func (v voucherService) CreateVoucher(ctx context.Context, id, voucherType string, properties []byte) error {
	err := v.repo.CreateVoucher(ctx, id, voucherType, properties)
	if err != nil {
		v.logger.Error(err.Error())
		return err
	}

	return nil
}

func (v voucherService) GetAllVouchers(ctx context.Context, offset, limit int) ([]models.Voucher, error) {
	vouchers, err := v.repo.GetAllVouchers(ctx, offset, limit)
	if err != nil {
		v.logger.Error(err.Error())
		return nil, err
	}

	return vouchers, nil
}

func (v voucherService) UpdateVoucher(ctx context.Context, id, voucherType string, properties []byte) error {
	err := v.repo.UpdateVoucher(ctx, id, voucherType, properties)
	if err != nil {
		v.logger.Error(err.Error())
		return err
	}

	return nil
}

func (v voucherService) DeleteVoucher(ctx context.Context, id string) error {
	err := v.repo.DeleteVoucher(ctx, id)
	if err != nil {
		v.logger.Error(err.Error())
		return err
	}

	return nil
}
