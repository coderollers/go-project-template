package api

import (
	"context"
	"time"

	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	handlersV1 "my-microservice/api/handlers/v1"
	"my-microservice/api/middleware"
	"my-microservice/configuration"
	"my-microservice/tracer"
)

func StartGin(ctx context.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Set up telemetry
	otTimeout := conf.CleanupTimeoutSec / 2
	if conf.CleanupTimeoutSec > 10 {
		otTimeout = 10
	}
	if conf.UseTelemetry == "remote" {
		log.Infof("Remote Telemetry enabled")
		tp, err := tracer.InitTracerJaeger(ctx, conf.JaegerEndpoint, configuration.OTName)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			concurrency.GlobalWaitGroup.Add(1)
			defer concurrency.GlobalWaitGroup.Done()
			localCtx, localCancel := context.WithTimeout(context.Background(), time.Duration(otTimeout)*time.Second)
			defer localCancel()
			if err := tp.Shutdown(localCtx); err != nil {
				log.Errorf("Error shutting down tracer provider: %s", err.Error())
			}
		}()
	}

	if conf.UseTelemetry == "local" {
		log.Infof("Stdout Telemetry enabled")
		tp, err := tracer.InitTracerStdout(ctx)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			concurrency.GlobalWaitGroup.Add(1)
			defer concurrency.GlobalWaitGroup.Done()
			localCtx, localCancel := context.WithTimeout(context.Background(), time.Duration(otTimeout)*time.Second)
			defer localCancel()
			if err := tp.Shutdown(localCtx); err != nil {
				log.Errorf("Error shutting down tracer provider: %s", err.Error())
			}
		}()
	}

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
}
