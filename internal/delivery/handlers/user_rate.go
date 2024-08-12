package handlers

import (
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
// @Param page query int true "Page"
// @Param per_page query int false "Per page"
// @Param user_id query int true "User id"
// @Success 200 {object} []models.UserRate "Array"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate/user [get]
func (s *UserRateHandler) GetUserRates(c *gin.Context) {
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

	rates, err := s.service.GetByUser(ctx, filter.UserID, filter.Page, filter.PerPage)
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

	rate, err := s.service.Get(ctx, filter.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// ClaimOutcome @Summary Claim from outcome pool
// @Tags rate
// @Accept  json
// @Produce  json
// @Param amount query int true "amount"
// @Param user_rate_id query int true "userRateID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate/claim/outcome [put]
func (s *UserRateHandler) ClaimOutcome(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		UserRateID int `form:"user_rate_id"`
		Amount     int `form:"amount"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = s.service.ClaimOutcome(ctx, filter.UserRateID, filter.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ClaimDeposit @Summary Claim from deposit
// @Tags rate
// @Accept  json
// @Produce  json
// @Param amount query int true "amount"
// @Param user_rate_id query int true "userRateID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rate/claim/deposit [put]
func (s *UserRateHandler) ClaimDeposit(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		UserRateID int `form:"user_rate_id"`
		Amount     int `form:"amount"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = s.service.ClaimDeposit(ctx, filter.UserRateID, filter.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
