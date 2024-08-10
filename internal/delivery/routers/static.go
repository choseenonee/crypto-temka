package routers

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterStaticRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	staticRouter := r.Group("/static")

	staticRepo := repository.InitStatic(db)

	staticService := service.InitStatic(staticRepo, logger)

	staticHandler := handlers.InitStaticHandler(staticService)

	staticRouter.POST("/review", staticHandler.CreateReview)
	staticRouter.GET("/review", staticHandler.GetReviews)
	staticRouter.PUT("/review", staticHandler.UpdateReview)
	staticRouter.DELETE("/review", staticHandler.DeleteReview)

	return staticRouter
}
