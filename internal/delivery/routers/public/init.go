package public

import (
	"os"

	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitPublicRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File) {
	public := r.Group("/public")
	_ = RegisterPublicStaticRouter(public, db, logger, metricsSetFile)
	_ = RegisterPublicRateRouter(public, db, logger)
	_ = RegisterPublicUserRouter(public, db, logger)
	_ = RegisterPublicWalletRouter(public, db)
}
