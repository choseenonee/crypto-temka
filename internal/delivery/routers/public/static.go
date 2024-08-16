package public

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func RegisterPublicStaticRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs, metricsSetFile *os.File) *gin.RouterGroup {
	staticRouter := r.Group("/static")

	staticRepo := repository.InitStatic(db)

	staticService := service.InitStatic(staticRepo, logger, metricsSetFile)

	staticHandler := handlers.InitStaticHandler(staticService)

	staticRouter.GET("/review", staticHandler.GetReviews)

	staticRouter.GET("/case", staticHandler.Get)
	staticRouter.GET("/case/all", staticHandler.GetCases)

	staticRouter.GET("/metrics", staticHandler.GetMetrics)

	return staticRouter
}
