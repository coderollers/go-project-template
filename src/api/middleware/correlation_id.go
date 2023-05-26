package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var correlationId string
		if correlationId = c.Request.Header.Get("X-Correlation-ID"); correlationId == "" {
			correlationId = uuid.New().String()
		}
		c.Set("correlation_id", correlationId)
		c.Next()
	}
}
