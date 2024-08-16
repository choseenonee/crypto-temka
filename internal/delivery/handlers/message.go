package handlers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageHandler struct {
	service service.Message
}

func InitMessageHandler(serv service.Message) MessageHandler {
	return MessageHandler{service: serv}
}

// Create @Summary Create message
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.MessageCreate true "data"
// @Success 200 {object} int "Successfully created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/message [post]
func (m *MessageHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var mc models.MessageCreate

	if err := c.ShouldBindJSON(&mc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := m.service.Create(ctx, mc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// Get @Summary Get message by id
// @Tags message
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id query int true "message id"
// @Success 200 {object} models.Message "message"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /message [get]
func (m *MessageHandler) Get(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get message handler"})
		return
	}

	message, err := m.service.GetByID(ctx, filter.ID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// GetByUser @Summary GetByUser messages by user
// @Tags message
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} []models.Message "message"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /message/user [get]
func (m *MessageHandler) GetByUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in create user rate handler"})
		return
	}

	messages, err := m.service.GetByUser(ctx, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
