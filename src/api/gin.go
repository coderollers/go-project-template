package api

import (
	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	handlersV1 "my-microservice/api/handlers/v1"
	"my-microservice/api/middleware"
	"my-microservice/configuration"
)

func SetupGin() *gin.Engine {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Set up gin
	log.Debugf("Setting up Gin")
	if !conf.GinLogger {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Set up the middleware
	if conf.GinLogger {
		log.Warnf("Gin's logger is active! Logs will be unstructured!")
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())
	router.Use(otelgin.Middleware(configuration.OTName))
	// TEMPLATE: Add more middleware

	userAPI := router.Group("/v1")
	{
		userAPI.GET("/", handlersV1.IndexGet)
		// TEMPLATE: Add more handlers
	}
	// Activate swagger if configured
	if conf.UseSwagger {
		log.Infof("Swagger is active, enabling endpoints")
		url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	return router
}
