package api

import (
	"context"
	"time"

	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"

	"my-microservice/configuration"
	"my-microservice/tracer"
)

func StartTelemetry(ctx context.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Set up telemetry
	otTimeout := conf.CleanupTimeoutSec / 2
	if conf.CleanupTimeoutSec > 10 {
		otTimeout = 10
	}
	if conf.JaegerEndpoint == "stdout" {
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
	} else if conf.JaegerEndpoint != "" {
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
}
