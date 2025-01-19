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

func RegisterReferRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, getStatus func(id int) (string, error),
	mdw middleware.Middleware) *gin.RouterGroup {
	router := r.Group("/refer")

	router.Use(mdw.Authorization(false, getStatus, []string{"verified"}))

	repo := repository.InitRefer(db)
	serv := service.InitRefer(repo, logger)

	handler := handlers.InitReferHandler(serv)

	router.GET("", handler.Get)
	router.PUT("", handler.Claim)

	return router
}
