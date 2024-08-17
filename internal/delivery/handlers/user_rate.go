package handlers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRateHandler struct {
	service service.UserRate
}

func InitUserRateHandler(serv service.UserRate) UserRateHandler {
	return UserRateHandler{service: serv}
}

// CreateUserRate @Summary Create
// @Tags rate
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.UserRateCreate true "data"
// @Success 200 {object} int "Successfully created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate [post]
func (s *UserRateHandler) CreateUserRate(c *gin.Context) {
	ctx := c.Request.Context()

	var urc models.UserRateCreate
	if err := c.ShouldBindJSON(&urc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in create user rate handler"})
		return
	}

	urc.UserID = userID.(int)

	id, err := s.service.Create(ctx, urc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// GetUserRates @Summary Get
// @Tags rate
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param page query int true "Page"
// @Param per_page query int false "Per page"
// @Success 200 {object} []models.UserRate "Array"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate/user [get]
func (s *UserRateHandler) GetUserRates(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		Page    int `form:"page"`
		PerPage int `form:"per_page"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get user rates handler"})
		return
	}

	rates, err := s.service.GetByUser(ctx, userID.(int), filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// GetUserRate @Summary Get
// @Tags rate
// @Accept  json
// @Produce  json
// @Param id query int true "id"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.UserRate "User Rate"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate [get]
func (s *UserRateHandler) GetUserRate(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		ID int `form:"id"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get user rate handler"})
		return
	}

	rate, err := s.service.Get(ctx, filter.ID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// Claim @Summary Claim from outcome pool
// @Tags rate
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param amount query float64 true "amount"
// @Param user_rate_id query int true "userRateID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate/claim [put]
func (s *UserRateHandler) Claim(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		UserRateID int     `form:"user_rate_id"`
		Amount     float64 `form:"amount"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get user rate handler"})
		return
	}

	err = s.service.Claim(ctx, filter.UserRateID, userID.(int), filter.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetAll @Summary GetAll
// @Tags admin
// @Accept  json
// @Produce  json
// @Param user_id query int false "userID as filter, can be null"
// @Param page query int true "page"
// @Param per_page query int true "perPage"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} []models.UserRateAdmin "User Rate admin"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/user_rate/all [get]
func (s *UserRateHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		Page    int `form:"page"`
		PerPage int `form:"per_page"`
		UserID  int `form:"user_id"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	rates, err := s.service.GetAll(ctx, filter.UserID, filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// UpdateNextDayCharge @Summary UpdateNextDayCharge
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param amount query float64 true "amount"
// @Param user_rate_id query int true "userRateID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/user_rate/next_day_charge [put]
func (s *UserRateHandler) UpdateNextDayCharge(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		UserRateID int     `form:"user_rate_id"`
		Amount     float64 `form:"amount"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = s.service.UpdateNextDayCharge(ctx, filter.UserRateID, filter.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
