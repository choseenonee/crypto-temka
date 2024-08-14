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

	admin.Use(mdw.Authorization(true))

	_ = RegisterAdminStaticRouter(admin, db, logger, metricsSetFile)

	_ = RegisterAdminRateRouter(admin, db, logger)

	_ = RegisterPublicUserRouter(admin, db, logger)
}
