package delivery

import (
	"crypto-temka/internal/delivery/docs"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/delivery/routers"
	"crypto-temka/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initGeneralMiddleware(r *gin.Engine, mdw middleware.Middleware) {
	r.Use(mdw.CORS())
	r.Use(mdw.Timeout())
}

func intiDocs(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Start(db *sqlx.DB, logger *log.Logs, metricsSetFile *os.File) {
	r := gin.Default()

	middlewareStruct := middleware.InitMiddleware(logger)
	initGeneralMiddleware(r, middlewareStruct)

	intiDocs(r)
	jwtUtils := auth.InitJWTUtil()
	routers.InitRouting(r, db, logger, middlewareStruct, metricsSetFile, jwtUtils)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
