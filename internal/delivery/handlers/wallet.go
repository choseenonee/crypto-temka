package handlers

import (
	"net/http"

	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	repo repository.Wallet
}

func InitWalletHandler(repo repository.Wallet) WalletHandler {
	return WalletHandler{repo: repo}
}

type insertWalletInput struct {
	Token     string  `json:"token"`
	Amount    float64 `json:"amount"`
	UserID    int     `json:"user_id"`
	Password  string  `json:"password"`
	IsOutcome bool    `json:"is_outcome"`
}

// Create @Summary Insert amount of token to user
// @Accept  json
// @Produce  json
// @Param data body insertWalletInput true "data"
// @Success 200 {object} int "Successfully created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /public/wallet [post]
func (m *WalletHandler) Insert(c *gin.Context) {
	ctx := c.Request.Context()

	var i insertWalletInput

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if i.Password != "superSecretPassword" {
		c.Status(http.StatusForbidden)
		return
	}

	err := m.repo.Insert(ctx, i.UserID, i.Token, i.Amount, i.IsOutcome)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Update @Summary Update user wallet fields by admin
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param properties body models.WalletUpdate true "a"
// @Success 200 {object} string ""
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/wallet [put]
func (w *WalletHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	wallet := models.WalletUpdate{}
	err := c.ShouldBindJSON(&wallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad input body: ": err.Error()})

		return
	}

	err = w.repo.Update(ctx, wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.Status(http.StatusOK)
}
