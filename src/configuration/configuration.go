package configuration

import (
	"github.com/coderollers/go-utils"
)

type Configuration struct {
	Swagger CSwagger

	// Dependencies section

	// JaegerEndpoint of the Jaeger instance where you want to send telemetry data. Optional, see UseTelemetry.
	JaegerEndpoint string

	// Internal settings section

	// CleanupTimeoutSec sets how long the microservice will wait for goroutines to
	// end before forcibly exiting when it receives a termination signal from the
	// orchestrator or OS.
	CleanupTimeoutSec int32
	// Environment is a string representing the environment where the microservice is
	// deployed, such as "staging" or "production". Optional.
	Environment string
	// UseTelemetry sets the behavior of OpenTelemetry.
	// "remote" will push telemetry data to a OT-compatible server, such as Jaeger. See JaegerEndpoint.
	// "local" will activate telemetry output on standard output. Do not use in production!
	// Any other value will disable OpenTelemetry.
	UseTelemetry string
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

	// HttpPort controls the TCP port that Gin will be listening on for HTTP connections.
	HttpPort int32
	// IngressPrefix must match the path which routes requests to this microservice
	// in your Ingress configuration. See the README for more information.
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
}
