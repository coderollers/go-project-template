package response

import (
	"fmt"
	"math"
	"net/http"

	"github.com/coderollers/go-utils"
	"github.com/gin-gonic/gin"

	"my-microservice/api/models"
	"my-microservice/configuration"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.JSONSuccessResult{
		Code:          http.StatusOK,
		Data:          data,
		Message:       "Success",
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}

func AcceptedResponse(c *gin.Context, id string, data interface{}) {
	c.JSON(http.StatusAccepted, models.JSONAcceptedResult{
		Code:          http.StatusAccepted,
		Id:            id,
		Data:          data,
		Message:       "Accepted",
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}

func FailureResponse(c *gin.Context, data interface{}, err utils.HttpError) {
	if err.Err == nil {
		err = utils.HttpError{Code: int(math.Max(float64(err.Code), 500)), Err: fmt.Errorf("FailureResponse was called with a nil error (%s)", err.Message)}
	}
	var errorString, stackString string
	conf := configuration.AppConfig()
	if conf.Development {
		errorString = err.Error()
		stackString = err.StackTrace()
	}
	c.JSON(err.Code, models.JSONFailureResult{
		Code:          err.Code,
		Data:          data,
		Error:         errorString,
		Stack:         stackString,
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}

func NotFoundResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusNotFound, models.JSONNotFoundResult{
		Code:          http.StatusNotFound,
		Data:          data,
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}
