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

func RegisterUserRateRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware) *gin.RouterGroup {
	router := r.Group("/rate")

	router.Use(mdw.Authorization(false))

	userRateRepo := repository.InitUsersRate(db)
	walletRepo := repository.InitWallet(db)
	rateRepo := repository.InitRate(db)

	userRateService := service.InitUserRate(userRateRepo, walletRepo, rateRepo, logger)

	handler := handlers.InitUserRateHandler(userRateService)

	router.POST("", handler.CreateUserRate)
	router.GET("/user", handler.GetUserRates)
	router.GET("", handler.GetUserRate)
	router.PUT("/claim", handler.Claim)
	//router.PUT("/claim/deposit", handler.ClaimDeposit)

	return router
}
