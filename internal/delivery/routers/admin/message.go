package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminMessageRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/message")

	repo := repository.InitMessage(db)

	serv := service.InitMessage(repo, logger)

	handler := handlers.InitMessageHandler(serv)

	router.POST("", handler.Create)

	return router
}
