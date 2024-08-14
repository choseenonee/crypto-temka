package handlers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WithdrawHandler struct {
	service service.Withdraw
}

func InitWithdrawHandler(serv service.Withdraw) WithdrawHandler {
	return WithdrawHandler{service: serv}
}

// Create @Summary Create
// @Tags withdraw
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.WithdrawCreate true "data"
// @Success 200 {object} int "id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /withdraw [post]
func (w *WithdrawHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var wc models.WithdrawCreate
	if err := c.ShouldBindJSON(&wc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get me handler"})
		return
	}

	wc.UserID = userID.(int)

	id, err := w.service.Create(ctx, wc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// GetByUserID @Summary GetByUserID
// @Tags withdraw
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param page query int true "page"
// @Param per_page query int true "per_page"
// @Success 200 {object} []models.Withdraw "Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /withdraw/user [get]
func (w *WithdrawHandler) GetByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		Page    int `form:"page"`
		PerPage int `form:"per_page"`
	}

	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get me handler"})
		return
	}

	withdrawals, err := w.service.GetByUserID(ctx, filter.Page, filter.PerPage, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

// GetByID @Summary GetByID
// @Tags withdraw
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param withdraw_id query int true "withdraw id"
// @Success 200 {object} models.Withdraw "withdraw"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /withdraw [get]
func (w *WithdrawHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		WithdrawID int `form:"withdraw_id"`
	}

	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get me handler"})
		return
	}

	withdraw, err := w.service.GetByID(ctx, filter.WithdrawID, userID.(int))
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdraw)
}

// TODO: доделать withdraws, а также прописать 403 для forbidden

// GetAll @Summary Get all
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param page query int true "page"
// @Param per_page query int true "per page"
// @Param status query string false "filter by status"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/withdraw/all [get]
func (w *WithdrawHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		Page    int     `form:"page"`
		PerPage int     `form:"per_page"`
		Status  *string `form:"status"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	var status string
	if filter.Status != nil {
		status = *filter.Status
	}

	withdrawals, err := w.service.GetAll(ctx, filter.Page, filter.PerPage, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

type payloadStatus struct {
	WithdrawID int         `form:"withdraw_id" json:"withdraw_id"`
	Status     string      `form:"status"`
	Properties interface{} `form:"properties"`
}

// UpdateStatus @Summary UpdateStatus
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body payloadStatus true "data"
// @Success 200 {object} string ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/withdraw/status [put]
func (w *WithdrawHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var pld payloadStatus
	if err := c.ShouldBindJSON(&pld); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := w.service.UpdateStatus(ctx, pld.WithdrawID, pld.Status, pld.Properties)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
