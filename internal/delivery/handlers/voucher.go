package handlers

import (
	"crypto-temka/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	service service.Voucher
}

func InitVoucherHandler(serv service.Voucher) VoucherHandler {
	return VoucherHandler{service: serv}
}

type voucherCreateInput struct {
	Id          string      `json:"id" binding:"required"`
	VoucherType string      `json:"voucher_type" binding:"required"`
	Properties  interface{} `json:"properties" binding:"required"`
}

// CreateVoucher @Summary CreateVoucher
// @Tags voucher
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body voucherCreateInput true "input"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/voucher [post]
func (u *VoucherHandler) CreateVoucher(c *gin.Context) {
	var input voucherCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := u.service.CreateVoucher(c.Request.Context(), input.Id, input.VoucherType, input.Properties)
	if err != nil {
		if errors.Is(err, service.ErrInvalidJSON) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type getAllVouchersInput struct {
	Offset *int `form:"offset" binding:"required"`
	Limit  *int `form:"limit" binding:"required"`
}

// GetAllVouchers @Summary GetAllVouchers
// @Tags voucher
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data query getAllVouchersInput true "data"
// @Success 200 {object} []models.Voucher "Vouchers"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/voucher [get]
func (u *VoucherHandler) GetAllVouchers(c *gin.Context) {
	var input getAllVouchersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vouchers, err := u.service.GetAllVouchers(c.Request.Context(), *input.Offset, *input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vouchers)
}

type updateVoucherInput struct {
	Id          string      `json:"id" binding:"required"`
	VoucherType string      `json:"voucher_type" binding:"required"`
	Properties  interface{} `json:"properties" binding:"required"`
}

// UpdateVoucher @Summary UpdateVoucher
// @Tags voucher
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body updateVoucherInput true "data"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/voucher [put]
func (u *VoucherHandler) UpdateVoucher(c *gin.Context) {
	var input updateVoucherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.service.UpdateVoucher(c.Request.Context(), input.Id, input.VoucherType, input.Properties)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type deleteVoucherInput struct {
	Id string `form:"id" binding:"required"`
}

// DeleteVoucher @Summary DeleteVoucher
// @Tags voucher
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data query deleteVoucherInput true "data"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/voucher [delete]
func (u *VoucherHandler) DeleteVoucher(c *gin.Context) {
	var input deleteVoucherInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.service.DeleteVoucher(c.Request.Context(), input.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
