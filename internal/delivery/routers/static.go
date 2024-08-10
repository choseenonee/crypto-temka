package routers

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func RegisterStaticRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, metricsSetFile *os.File) *gin.RouterGroup {
	staticRouter := r.Group("/static")

	staticRepo := repository.InitStatic(db)

	staticService := service.InitStatic(staticRepo, logger, metricsSetFile)

	staticHandler := handlers.InitStaticHandler(staticService)

	staticRouter.POST("/review", staticHandler.CreateReview)
	staticRouter.GET("/review", staticHandler.GetReviews)
	staticRouter.PUT("/review", staticHandler.UpdateReview)
	staticRouter.DELETE("/review", staticHandler.DeleteReview)

	staticRouter.POST("/metrics", staticHandler.UpdateMetrics)
	staticRouter.GET("/metrics", staticHandler.GetMetrics)

	return staticRouter
}
