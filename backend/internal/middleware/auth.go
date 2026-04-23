package middleware

import (
	"strings"

	"epbms/internal/domain"
	"epbms/pkg/response"
	"epbms/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

// JWTAuth extracts and validates the Bearer token from the Authorization header,
// then injects the user_id and role into the Gin context for downstream handlers.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			response.Error(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

// RequireRoles checks that the authenticated user has one of the allowed roles.
// Must be used after JWTAuth middleware.
func RequireRoles(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextKeyRole)
		if !exists {
			response.Error(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		role, ok := roleVal.(domain.Role)
		if !ok {
			response.Error(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}

		if _, permitted := allowed[role]; !permitted {
			response.Error(c, domain.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCallerID extracts the authenticated user's ID from the Gin context.
func GetCallerID(c *gin.Context) uint {
	v, _ := c.Get(ContextKeyUserID)
	id, _ := v.(uint)
	return id
}

// GetCallerRole extracts the authenticated user's role from the Gin context.
func GetCallerRole(c *gin.Context) domain.Role {
	v, _ := c.Get(ContextKeyRole)
	role, _ := v.(domain.Role)
	return role
}
