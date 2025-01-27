package handlers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.User
}

func InitUserHandler(serv service.User) UserHandler {
	return UserHandler{service: serv}
}

// Register @Summary Register
// @Tags public
// @Accept  json
// @Produce  json
// @Param data body models.UserCreate true "data"
// @Success 200 {object} string "Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/user [post]
func (u *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var uc models.UserCreate
	if err := c.ShouldBindJSON(&uc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := u.service.Register(ctx, uc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

type payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth @Summary Auth
// @Tags public
// @Accept  json
// @Produce  json
// @Param data body payload true "data"
// @Success 200 {object} string "Token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/user [put]
func (u *UserHandler) Auth(c *gin.Context) {
	ctx := c.Request.Context()

	var pld payload
	if err := c.ShouldBindJSON(&pld); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := u.service.Auth(ctx, pld.Email, pld.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetMe @Summary Get me
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.User "User"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user [get]
func (u *UserHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CUserID was not found in get me handler"})
		return
	}

	user, err := u.service.Get(ctx, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAll @Summary Get all users
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
// @Router /admin/user/all [get]
func (u *UserHandler) GetAll(c *gin.Context) {
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

	users, err := u.service.GetAll(ctx, filter.Page, filter.PerPage, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateStatus @Summary UpdateStatus
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param user_id query int true "user_id"
// @Param status query string true "status"
// @Success 200 {object} string ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/user/status [put]
func (u *UserHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		UserID int    `form:"user_id"`
		Status string `form:"status"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	err = u.service.UpdateStatus(ctx, filter.UserID, filter.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type updatePropertiesInput struct {
	Properties interface{} `json:"properties"`
}

// UpdateProperties @Summary UpdateProperties
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param properties body updatePropertiesInput true "new properties"
// @Param start-verify query bool true "if true passed, will change user status to pending. if false, won't change status."
// @Success 200 {object} string ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/properties [put]
func (u *UserHandler) UpdateProperties(c *gin.Context) {
	ctx := c.Request.Context()

	var filter struct {
		StartVerify bool `form:"start-verify"`
	}

	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad query: ": err.Error()})
		return
	}

	userID, ok := c.Get(middleware.CUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "constant CUserID was not found in update properties handler"})
		return
	}

	p := updatePropertiesInput{}
	err = c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad input body: ": err.Error()})
		return
	}

	err = u.service.UpdateProperties(ctx, userID.(int), p.Properties, filter.StartVerify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
