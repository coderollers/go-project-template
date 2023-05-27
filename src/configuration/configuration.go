package configuration

import (
	"github.com/coderollers/go-utils"
)

type Configuration struct {
	Swagger CSwagger

	// Dependencies section

	// JaegerEndpoint of the Jaeger instance where you want to send telemetry data.
	// Set to "stdout" for activating standard output telemetry or leave empty to
	// disable telemetry.
	JaegerEndpoint string

	// Internal settings section

	// CleanupTimeoutSec sets how long the microservice will wait for goroutines to
	// end before forcibly exiting when it receives a termination signal from the
	// orchestrator or OS.
	CleanupTimeoutSec int32
	// Environment is a string representing the environment where the microservice is
	// deployed, such as "staging" or "production". Optional.
	Environment string
	// Development, if true, will activate development features and DEBUG level logs. Do not activate in production!
	Development bool
	// GinLogger, if true, will activate Gin's internal logger. Use for debugging
	// purposes. Will break structured (json) logging. Do not activate in production!
	GinLogger bool
	// UseSwagger, if true, will activate the swagger endpoint. Do not use in production!
	UseSwagger bool
	// Initialized controls if the configuration has been loaded or not.
	Initialized bool

	// Microservice configuration section

	// HttpPort controls the TCP port that Gin will be listening on for HTTP
	// connections. Defaults to 8080.
	HttpPort int32
	// GrpcPort controls the TCP port that the GRPC services will be available on.
	// Defaults to 9000. If GrpcPort and HttpPort are the same, then grpc-web
	// compatibility will be enabled. This will allow you to call the GRPC services
	// from web clients such as JavaScript and WebAssembly
	GrpcPort int32
	// IngressPrefix must match the path which routes requests to this microservice
	// in your Ingress configuration. Only affects the HTTP server. Note that this
	// will break your grpc-web endpoints, if grpc-web is enabled! See the README for
	// more information.
	IngressPrefix string

	// TEMPLATE: Add more configuration data here or to the above sections to be
	// available throughout the microservice code
}

var appConfig Configuration

// AppConfig returns an instance of the global configuration. Use it to access
// the configuration from anywhere in the code.
func AppConfig() *Configuration {
	if appConfig.Initialized == false {
		loadEnvironmentVariables()
		appConfig.Initialized = true
	}
	return &appConfig
}

// loadEnvironmentVariables is used to set the configuration options based on
// environment variables.
//
// Caution! This only runs once and will not pick up runtime environment variable
// changes!
func loadEnvironmentVariables() {
	appConfig.Environment = utils.EnvOrDefault("ENVIRONMENT", "local")
	appConfig.JaegerEndpoint = utils.EnvOrDefault("JAEGER_ENDPOINT", "")
	appConfig.IngressPrefix = utils.EnvOrDefault("INGRESS_PREFIX", "")
	appConfig.CleanupTimeoutSec = utils.EnvOrDefaultInt32("SHUTDOWN_TIMEOUT", 300)
	appConfig.HttpPort = utils.EnvOrDefaultInt32("HTTP_PORT", 8080)
	appConfig.GrpcPort = utils.EnvOrDefaultInt32("GRPC_PORT", 9000)
}
