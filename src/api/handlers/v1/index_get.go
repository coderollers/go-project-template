package v1

import (
	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/coderollers/go-utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"my-microservice/api/response"
	"my-microservice/tracer"
)

// IndexGet godoc
// @Summary Sample GET handler
// @Description Sample GET handler
// @ID index-get
// @Accept json
// @Produce json
// @Success 200 {object} models.JSONSuccessResult "Positive response"
// @Failure 400 {object} models.JSONFailureResult "The request data could not be processed"
// @Failure 404 {object} models.JSONNotFoundResult "The object was not found"
// @Failure 500 {object} models.JSONFailureResult "An internal error has occurred, most likely due to an uncaught exception"
// @Failure 503 {object} models.JSONFailureResult "An error has occurred, most likely due to an unavailable dependency"
// @Router /v1/ [get]
func IndexGet(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()
	var (
		log           = logger.SugaredLogger().WithContextCorrelationId(c).With("package", "handlers", "action", "GetTask")
		correlationId = c.MustGet("correlation_id").(string)
		r             interface{}
	)

	// Create tracer span
	_, span := tracer.Tracer.Start(c.Request.Context(), "IndexGet")
	defer span.End()

	// Add tracer event to span
	span.AddEvent("Index Get", trace.WithAttributes(attribute.String("CorrelationId", correlationId)))

	// Log debug (will only show in Development mode)
	log.Debugf("Correlation ID for request: %s", correlationId)

	// Do some work and get the response data you want to send back to the client
	responseData := map[string]string{"Motto": "Hello world!"}

	// Example not found response
	if responseData == nil {
		response.NotFoundResponse(c, r)
		return // Always return after responding to client!
	}

	// Example failure response
	// err := json.Unmarshal([]byte{1, 0, 1}, r) // Fails
	var err error // Works
	if err != nil {
		response.FailureResponse(c, nil, utils.HttpError{
			Code:    400,
			Err:     err,
			Message: "There was an unexpected error processing the request",
		})
		return // Always return after responding to client!
	}

	// Example positive response
	response.SuccessResponse(c, responseData)
}
