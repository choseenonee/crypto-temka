package handlers

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	service service.UserService
}

func InitUserHandler(serv service.UserService) UserHandler {
	return UserHandler{service: serv}
}

// @Summary Create user
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body models.UserCreate true "user create"
// @Success 200 {object} int "Successfully created user"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/create [post]
func (handler UserHandler) Create(g *gin.Context) {
	var newUser models.UserCreate

	if err := g.ShouldBindJSON(&newUser); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	id, token, err := handler.service.CreateUser(ctx, newUser)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"id": id, "token": token})
}

// @Summary Login user
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body models.UserLogin true "user login"
// @Success 200 {object} int "Successfully login user"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/login [post]
func (handler UserHandler) Login(g *gin.Context) {
	var user models.UserLogin

	if err := g.ShouldBindJSON(&user); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := g.Request.Context()

	id, token, err := handler.service.Login(ctx, user)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"userID": id, "token": token})
}

// @Summary Get user using ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query int true "UserID"
// @Success 200 {object} int "Successfully get user"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/{id} [get]
func (handler UserHandler) Get(c *gin.Context) {
	id := c.Query("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	user, err := handler.service.GetUserByID(ctx, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary Delete user
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query int true "UserID"
// @Success 200 {object} int "Successfully deleted"
// @Failure 400 {object} map[string]string "Invalid id"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/{id} [delete]
func (handler UserHandler) Delete(g *gin.Context) {
	userID := g.Query("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = handler.service.Delete(ctx, id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"delete": id})
}

// @Summary Add photo
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body models.UserPhoto true "add photo"
// @Success 200 {object} int "Success changing"
// @Failure 400 {object} map[string]string "Invalid id"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/photo/add [put]
func (handler UserHandler) AddPhoto(g *gin.Context) {
	var user models.UserPhoto
	if err := g.ShouldBindJSON(&user); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := handler.service.AddPhoto(ctx, user)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"add": "success", "id": user.ID})
}

// @Summary Get user photo
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query int true "UserID"
// @Success 200 {object} int "Successfully get user photo"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/photo/{id} [get]
func (handler UserHandler) GetPhoto(c *gin.Context) {
	id := c.Query("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	photo, err := handler.service.GetPhoto(ctx, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photo": photo})
}

// @Summary Set Status
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body models.SetStatus true "Status"
// @Success 200 {object} map[string]string "Successfully set status"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/status [put]
func (handler UserHandler) SetStatus(g *gin.Context) {
	var status models.SetStatus
	if err := g.ShouldBindJSON(&status); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := handler.service.SetStatus(ctx, status)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"set": "success", "id": status.ID})
}

// @Summary Refresh tokens
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query int true "UserID"
// @Success 200 {object} map[string]string "Successfully authorized, returning JWT and new session_id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/refresh [post]
func (handler UserHandler) Refresh(g *gin.Context) {
	id := g.Query("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userToken, err := handler.service.Refresh(ctx, aid)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"JWT": userToken})
}
