package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/spf13/pflag"

	"my-microservice/api"
	"my-microservice/configuration"
	"my-microservice/docs"
)

func main() {
	// Main context and cancellation tokens
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	// Initialize configuration
	appConfig := configuration.AppConfig()

	// Configure command-line parameters
	pflag.Int32VarP(&appConfig.CleanupTimeoutSec, "timeout", "t", 60, "Time to wait for graceful shutdown on SIGTERM/SIGINT in seconds. Default: 60")
	pflag.Int32VarP(&appConfig.HttpPort, "port", "p", 8080, "TCP port for the HTTP listener to bind to. Default: 8080")
	pflag.BoolVarP(&appConfig.UseSwagger, "swagger", "s", false, "Activate swagger. Do not use this in Production!")
	pflag.BoolVarP(&appConfig.Development, "devel", "d", false, "Start in development mode. Implies --swagger. Do not use this in Production!")
	pflag.BoolVarP(&appConfig.GinLogger, "gin-logger", "g", false, "Activate Gin's logger, for debugging. Do not use this in Production!")
	pflag.StringVarP(&appConfig.UseTelemetry, "telemetry", "r", "", "Activate telemetry local or remote/jaeger")
	pflag.Parse()

	// Initialize main context and set up cancellation token for SIGINT/SIGQUIT
	ctx = context.Background()
	ctx, cancel = context.WithCancel(ctx)
	cSignal := make(chan os.Signal)
	signal.Notify(cSignal, os.Interrupt, syscall.SIGTERM)

	// Initialize logger
	logger.Init(ctx, true, appConfig.Development)
	logger.SetCorrelationIdFieldKey(configuration.CorrelationIdKey)
	logger.SetCorrelationIdContextKey(configuration.CorrelationIdKey)
	log := logger.SugaredLogger()
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()
	defer logger.PanicLogger()

	// Sanity checks
	if !appConfig.Development {
		if appConfig.CleanupTimeoutSec < 120 {
			log.Warnf("Cleanup timeout is set to %d seconds which might be too small for production mode!", appConfig.CleanupTimeoutSec)
		}

		// TEMPLATE: Add more sanity checks here
	}

	if appConfig.Development {
		appConfig.UseSwagger = true
	}

	if appConfig.UseSwagger {
		// TEMPLATE: Modify `swagger.yaml` with your project data
		// Remember to always run `swag init --parseDependency` after changing swagger comments on handlers
		appConfig.LoadSwaggerConf()
		docs.SwaggerInfo.Title = appConfig.Swagger.Title
		docs.SwaggerInfo.Version = appConfig.Swagger.Version
		docs.SwaggerInfo.BasePath = appConfig.IngressPrefix + appConfig.Swagger.BasePath
		docs.SwaggerInfo.Description = appConfig.Swagger.Description
	}
	log.Infof(docs.SwaggerInfo.BasePath)

	// TEMPLATE: Further initialization goes here (kms, database, etc)

	// Trigger context cancellation token on SIGINT/SIGTERM
	go func() {
		<-cSignal
		log.Warnf("SIGTERM received, attempting graceful exit.")
		cancel()
	}()

	// Start the API HTTP Server
	log.Info("Starting webapi handler")
	go api.StartGin(ctx)

	// Block until cancellation signal is received
	<-ctx.Done()

	// Clean up and attempt graceful exit
	log.Infof("Graceful shutdown initiated. Waiting for %d seconds before forced exit.", appConfig.CleanupTimeoutSec)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(appConfig.CleanupTimeoutSec))
	go func() {
		// Eventual clean-up logic would go in this block
		concurrency.GlobalWaitGroup.Wait()
		log.Infof("Cleanup done.")
		cancel()
	}()
	<-ctx.Done()
	log.Info("Exiting.")

}
