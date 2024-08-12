package admin

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func InitAdminRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File) {
	admin := r.Group("/admin")
	_ = RegisterStaticRouter(admin, db, logger, metricsSetFile)
	_ = RegisterRateRouter(admin, db, logger)
}
