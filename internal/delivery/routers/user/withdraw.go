package user

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterWithdrawRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, getStatus func(id int) (string, error),
	mdw middleware.Middleware) *gin.RouterGroup {

	router := r.Group("/withdraw")

	router.Use(mdw.Authorization(false, getStatus))

	repo := repository.InitWithdraw(db)
	walletRepo := repository.InitWallet(db)

	serv := service.InitWithdraw(repo, walletRepo, logger)

	handler := handlers.InitWithdrawHandler(serv)

	router.POST("", handler.Create)
	router.GET("/user", handler.GetByUserID)
	router.GET("", handler.GetByID)

	return router
}
