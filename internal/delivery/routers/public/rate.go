package public

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterPublicRateRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/rate")

	rateRepo := repository.InitRate(db)

	rateService := service.InitRate(rateRepo, logger)

	handler := handlers.InitRateHandler(rateService)

	router.GET("", handler.GetRates)

	return router
}
