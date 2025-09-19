package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Context keys
const (
	ContextUserIDKey = "auth_user_id"
	ContextRoleKey   = "auth_role"
)

// AuthRequired verifies token and stores claims in context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}
		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "auth misconfigured"})
			return
		}
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// ensure algorithm
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token", "error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token claims"})
			return
		}

		// Extract sub (user id) - might be numeric or string
		var userID int64
		if sub, exists := claims["sub"]; exists {
			switch v := sub.(type) {
			case float64:
				userID = int64(v)
			case string:
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					userID = parsed
				}
			}
		}

		// Extract role string if present
		var role string
		if r, ok := claims["role"].(string); ok {
			role = r
		}

		// Set in context
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, role)
		c.Next()
	}
}

// OwnerOrAdmin middleware checks that the :user_id param equals token sub or role contains "ADMIN" / "SUPERADMIN"
func OwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		uidStr := c.Param("user_id")
		if uidStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user_id required"})
			return
		}
		paramID, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid user_id param"})
			return
		}
		uidVal, exists := c.Get(ContextUserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthenticated"})
			return
		}
		tokenUID := uidVal.(int64)
		roleVal, _ := c.Get(ContextRoleKey)
		role := ""
		if roleVal != nil {
			role = roleVal.(string)
		}

		if tokenUID != paramID && !(strings.Contains(strings.ToUpper(role), "ADMIN")) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}
		c.Next()
	}
}

// OwnerOnly ensures :user_id matches token sub (no admin bypass)
func OwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		uidStr := c.Param("id")
		if uidStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user_id required"})
			return
		}
		paramID, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid user_id param"})
			return
		}
		uidVal, exists := c.Get(ContextUserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthenticated"})
			return
		}
		tokenUID := uidVal.(int64)
		if tokenUID != paramID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}
		c.Next()
	}
}
