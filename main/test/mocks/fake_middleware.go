package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FakeAuthMiddlewareAlways401 всегда возвращает 401 Unauthorized
func FakeAuthMiddlewareAlways401() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
