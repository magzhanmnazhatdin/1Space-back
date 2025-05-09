package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole пропускает дальше только пользователей с одной из переданных ролей
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIf, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
			return
		}
		role := roleIf.(string)
		for _, a := range allowed {
			if role == a {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
	}
}
