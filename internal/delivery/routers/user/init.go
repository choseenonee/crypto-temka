package user

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitUserRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware) {
	_ = RegisterUserRouter(r, db, logger, mdw)

	_ = RegisterUserRateRouter(r, db, logger, mdw)
}
