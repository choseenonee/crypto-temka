package user

import (
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitUserRouters(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware) {
	userRepo := repository.InitUser(db)

	jwt := auth.InitJWTUtil()

	walletRepo := repository.InitWallet(db)

	userService := service.InitUser(userRepo, walletRepo, jwt, logger)

	_ = RegisterUserRouter(r, userService, mdw)

	_ = RegisterUserRateRouter(r, db, logger, userService.GetStatus, mdw)

	_ = RegisterWithdrawRouter(r, db, logger, userService.GetStatus, mdw)

	_ = RegisterReferRouter(r, db, logger, userService.GetStatus, mdw)

	_ = RegisterMessageRouter(r, db, logger, userService.GetStatus, mdw)
}
