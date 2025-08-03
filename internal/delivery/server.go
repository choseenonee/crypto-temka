package delivery

import (
	"fmt"
	"os"

	"crypto-temka/internal/delivery/docs"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/routers"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func Start(db *sqlx.DB, logger *log.Logs, metricsSetFile *os.File) {
	r := gin.Default()

	middlewareStruct := middleware.InitMiddleware(logger)
	initGeneralMiddleware(r, middlewareStruct)

	intiDocs(r)

	routers.InitRouting(r, db, logger, middlewareStruct, metricsSetFile)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
