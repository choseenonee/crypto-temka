package public

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterPublicWalletRouter(r *gin.RouterGroup, db *sqlx.DB) *gin.RouterGroup {
	router := r.Group("/wallet")

	walletRepo := repository.InitWallet(db)

	handler := handlers.InitWalletHandler(walletRepo)

	router.POST("", handler.Insert)

	return router
}
