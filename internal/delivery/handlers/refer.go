package handlers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReferHandler struct {
	service service.Refer
}

func InitReferHandler(serv service.Refer) ReferHandler {
	return ReferHandler{service: serv}
}

// Claim @Summary Claim
// @Tags refer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id query int true "id"
// @Success 200 {object} int "Successfully created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /refer [put]
func (r *ReferHandler) Claim(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get refers handler"})
		return
	}

	err = r.service.Claim(ctx, filter.ID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Get @Summary Get
// @Tags refer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param page query int true "Page"
// @Param per_page query int true "Refers per page"
// @Success 200 {object} []models.Refer "Array"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /refer [get]
func (r *ReferHandler) Get(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get refers handler"})
		return
	}

	refers, err := r.service.Get(ctx, userID.(int), filter.Page, filter.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, refers)
}
