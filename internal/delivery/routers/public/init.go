package public

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func InitPublicRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File) {
	admin := r.Group("/public")
	_ = RegisterPublicStaticRouter(admin, db, logger, metricsSetFile)
	_ = RegisterPublicRateRouter(admin, db, logger)
	_ = RegisterPublicUserRouter(admin, db, logger)
}
