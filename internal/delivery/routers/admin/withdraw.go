package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminWithdrawRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/withdraw")

	repo := repository.InitWithdraw(db)

	walletRepo := repository.InitWallet(db)

	serv := service.InitWithdraw(repo, walletRepo, logger)

	handler := handlers.InitWithdrawHandler(serv)

	router.GET("/all", handler.GetAll)
	router.PUT("/status", handler.UpdateStatus)

	return router
}
