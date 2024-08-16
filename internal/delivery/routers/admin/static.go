package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func RegisterAdminStaticRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs, metricsSetFile *os.File) *gin.RouterGroup {
	staticRouter := r.Group("/static")

	staticRepo := repository.InitStatic(db)

	staticService := service.InitStatic(staticRepo, logger, metricsSetFile)

	staticHandler := handlers.InitStaticHandler(staticService)

	staticRouter.POST("/review", staticHandler.CreateReview)
	staticRouter.PUT("/review", staticHandler.UpdateReview)
	staticRouter.DELETE("/review", staticHandler.DeleteReview)

	staticRouter.POST("/metrics", staticHandler.UpdateMetrics)

	staticRouter.POST("/case", staticHandler.CreateCase)
	staticRouter.PUT("/case", staticHandler.UpdateCase)
	staticRouter.DELETE("/case", staticHandler.DeleteCase)

	return staticRouter
}
