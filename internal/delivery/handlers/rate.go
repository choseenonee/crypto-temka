package handlers

import (
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RateHandler struct {
	service service.Rate
}

func InitRateHandler(serv service.Rate) RateHandler {
	return RateHandler{service: serv}
}

// CreateRate @Summary Create
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.RateCreate true "data"
// @Success 200 {object} int "Successfully created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/rate [post]
func (s *RateHandler) CreateRate(c *gin.Context) {
	ctx := c.Request.Context()

	var rc models.RateCreate

	if err := c.ShouldBindJSON(&rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.service.CreateRate(ctx, rc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// GetRates @Summary Get
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int true "Page"
// @Param per_page query int true "Reviews per page"
// @Success 200 {object} []models.Rate "Array"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/rate [get]
func (s *RateHandler) GetRates(c *gin.Context) {
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

	rates, err := s.service.GetRates(ctx, filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// UpdateRate @Summary Update
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.Rate true "data"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/rate [put]
func (s *RateHandler) UpdateRate(c *gin.Context) {
	ctx := c.Request.Context()

	var ru models.Rate

	if err := c.ShouldBindJSON(&ru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := json.Marshal(ru.Properties); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid json in properties: ": err})
		return
	}

	err := s.service.UpdateRate(ctx, ru)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
