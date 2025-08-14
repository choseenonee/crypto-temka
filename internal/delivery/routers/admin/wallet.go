package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminWalletRouter(r *gin.RouterGroup, db *sqlx.DB, _ *log.Logs) *gin.RouterGroup {
	router := r.Group("/wallet")

	walletRepo := repository.InitWallet(db)

	handler := handlers.InitWalletHandler(walletRepo)

	router.PUT("", handler.Update)

	return router
}
