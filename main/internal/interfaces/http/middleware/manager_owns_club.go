// internal/interfaces/http/middleware/owner.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
)

// ManagerOwnsClub проверяет, что при роли "manager" текущий пользователь — владелец клуба.
// Если роль не "manager", просто пропускает дальше.
func ManagerOwnsClub(clubUC usecase.ClubUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIf, _ := c.Get("role")
		role := roleIf.(string)
		if role != "manager" {
			c.Next()
			return
		}

		uidIf, _ := c.Get("uid")
		uid := uidIf.(string)

		clubID := c.Param("id")
		club, err := clubUC.GetByID(c.Request.Context(), clubID)
		if err != nil || club.ManagerID != uid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}
