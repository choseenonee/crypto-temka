package admin

import (
	"os"

	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitAdminRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File) {
	admin := r.Group("/admin")

	admin.Use(mdw.Authorization(true, func(id int) (string, error) {
		panic("wtf")
	}, nil))

	_ = RegisterAdminStaticRouter(admin, db, logger, metricsSetFile)

	_ = RegisterAdminRateRouter(admin, db, logger)

	_ = RegisterAdminUserRouter(admin, db, logger)

	_ = RegisterAdminWithdrawRouter(admin, db, logger)

	_ = RegisterAdminMessageRouter(admin, db, logger)

	_ = RegisterAdminUserRateRouter(admin, db, logger)

	_ = RegisterAdminVoucherRouter(admin, db, logger)

	_ = RegisterAdminWalletRouter(admin, db, logger)
}
