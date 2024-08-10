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
// @Tags review
// @Accept  json
// @Produce  json
// @Param data body models.ReviewCreate true "Review create"
// @Success 200 {object} int "Successfully created review"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /static/review [post]
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
// @Tags review
// @Accept  json
// @Produce  json
// @Param page query int true "Page"
// @Param per_page query int true "Reviews per page"
// @Success 200 {object} []models.Review "Array of reviews"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /static/review [get]
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
// @Tags review
// @Accept  json
// @Produce  json
// @Param data body models.Review true "Review"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /static/review [put]
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
// @Tags review
// @Accept  json
// @Produce  json
// @Param id query int true "Review ID"
// @Success 200 {object} nil ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /static/review [delete]
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