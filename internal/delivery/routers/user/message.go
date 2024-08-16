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

func RegisterMessageRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, getStatus func(id int) (string, error),
	mdw middleware.Middleware) *gin.RouterGroup {
	router := r.Group("/message")

	router.Use(mdw.Authorization(false, getStatus))

	repo := repository.InitMessage(db)
	serv := service.InitMessage(repo, logger)

	handler := handlers.InitMessageHandler(serv)

	router.GET("", handler.Get)
	router.GET("/user", handler.GetByUser)

	return router
}
