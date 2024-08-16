package handlers

import (
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StaticHandler struct {
	service service.Static
}

func InitStaticHandler(serv service.Static) StaticHandler {
	return StaticHandler{service: serv}
}

// CreateReview @Summary Create review
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.ReviewCreate true "Review create"
// @Success 200 {object} int "Successfully created review"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/review [post]
func (s *StaticHandler) CreateReview(c *gin.Context) {
	ctx := c.Request.Context()

	var rc models.ReviewCreate

	if err := c.ShouldBindJSON(&rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.service.CreateReview(ctx, rc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// GetReviews @Summary Get reviews
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int true "Page"
// @Param per_page query int true "Reviews per page"
// @Success 200 {object} []models.Review "Array of reviews"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/static/review [get]
func (s *StaticHandler) GetReviews(c *gin.Context) {
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

	reviews, err := s.service.GetReviews(ctx, filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview @Summary Update review
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.Review true "Review"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/review [put]
func (s *StaticHandler) UpdateReview(c *gin.Context) {
	ctx := c.Request.Context()

	var ru models.Review

	if err := c.ShouldBindJSON(&ru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := json.Marshal(ru.Properties); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid json in properties: ": err})
		return
	}

	err := s.service.UpdateReview(ctx, ru)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteReview @Summary Delete review
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id query int true "Review ID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/review [delete]
func (s *StaticHandler) DeleteReview(c *gin.Context) {
	ctx := c.Request.Context()

	var id struct {
		ID int `form:"id"`
	}

	err := c.BindQuery(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = s.service.DeleteReview(ctx, id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateMetrics @Summary Update metrics
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.MetricsSet true "Metrics updated"
// @Success 200 {object} int "Successfully updated metrics"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/metrics [post]
func (s *StaticHandler) UpdateMetrics(c *gin.Context) {
	var ms models.MetricsSet

	if err := c.ShouldBindJSON(&ms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.service.SetMetrics(ms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetMetrics @Summary Get metrics
// @Tags public
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Metrics "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/static/metrics [get]
func (s *StaticHandler) GetMetrics(c *gin.Context) {
	metrics := s.service.GetMetrics()

	c.JSON(http.StatusOK, metrics)
}

// CreateCase @Summary Create case
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.CaseCreate true "Case create"
// @Success 200 {object} int "Successfully created review"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/case [post]
func (s *StaticHandler) CreateCase(c *gin.Context) {
	ctx := c.Request.Context()

	var cc models.CaseCreate

	if err := c.ShouldBindJSON(&cc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.service.CreateCase(ctx, cc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// Get @Summary Get case
// @Tags public
// @Accept  json
// @Produce  json
// @Param id query int true "id"
// @Success 200 {object} []models.Case "Array of cases"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/static/case [get]
func (s *StaticHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		ID int `form:"id"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	cs, err := s.service.GetCase(ctx, filter.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cs)
}

// GetCases @Summary Get cases
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int true "Page"
// @Param per_page query int true "Reviews per page"
// @Success 200 {object} []models.Case ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/static/case/all [get]
func (s *StaticHandler) GetCases(c *gin.Context) {
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

	cases, err := s.service.GetCases(ctx, filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cases)
}

// UpdateCase @Summary Update case
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.Case true "Case"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/case [put]
func (s *StaticHandler) UpdateCase(c *gin.Context) {
	ctx := c.Request.Context()

	var cu models.Case

	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := json.Marshal(cu.Properties); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid json in properties: ": err})
		return
	}

	err := s.service.UpdateCase(ctx, cu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteCase @Summary Delete case
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id query int true "Case ID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/static/case [delete]
func (s *StaticHandler) DeleteCase(c *gin.Context) {
	ctx := c.Request.Context()

	var id struct {
		ID int `form:"id"`
	}

	err := c.BindQuery(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = s.service.DeleteCase(ctx, id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
