package routers

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware, metricsSetFile *os.File, jwtUtil auth.JWTUtil) {
	_ = RegisterStaticRouter(r, db, logger, metricsSetFile)
	_ = RegisterUserRouter(r, db, jwtUtil, logger)
}
