package routers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/routers/admin"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File) {
	admin.InitAdminRouters(r, db, logger, mdw, metricsSetFile)
	_ = RegisterUserRateRouter(r, db, logger)
}
