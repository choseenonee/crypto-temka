package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminUserRateRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/user_rate")

	repo := repository.InitUsersRate(db)
	walletRepo := repository.InitWallet(db)
	rateRepo := repository.InitRate(db)

	serv := service.InitUserRate(repo, walletRepo, rateRepo, logger)

	handler := handlers.InitUserRateHandler(serv)

	router.GET("/all", handler.GetAll)
	router.PUT("/next_day_charge", handler.UpdateNextDayCharge)

	return router
}
