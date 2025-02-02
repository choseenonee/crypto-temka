package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminVoucherRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/voucher")

	voucherRepo := repository.InitVoucherRepo(db)

	voucherService := service.InitVoucherService(voucherRepo, logger)

	handler := handlers.InitVoucherHandler(voucherService)

	router.POST("", handler.CreateVoucher)
	router.GET("", handler.GetAllVouchers)
	router.PUT("", handler.UpdateVoucher)
	router.DELETE("", handler.DeleteVoucher)

	return router
}
