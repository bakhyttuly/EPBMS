package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIDValue := session.Get("user_id")
		roleValue := session.Get("role")

		if userIDValue == nil || roleValue == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		var userID uint

		switch v := userIDValue.(type) {
		case uint:
			userID = v
		case int:
			userID = uint(v)
		case int64:
			userID = uint(v)
		case float64:
			userID = uint(v)
		default:
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
