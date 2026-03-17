package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/auth"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/response"
)

const (
	requestIDKey = "request_id"
	claimsKey    = "auth_claims"
)

type rateBucket struct {
	Remaining int
	ResetAt   time.Time
}

var (
	rateLimitMu sync.Mutex
	rateLimits  = map[string]*rateBucket{}
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = time.Now().Format("20060102150405.000000")
		}
		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	value, _ := c.Get(requestIDKey)
	if requestID, ok := value.(string); ok {
		return requestID
	}
	return ""
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", "request_id", GetRequestID(c), "panic", recovered)
		response.Error(c, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", nil)
	})
}

func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		logger.Info("http request",
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(started).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func Authenticate(manager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "missing bearer token", "UNAUTHORIZED", nil)
			c.Abort()
			return
		}
		claims, err := manager.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid access token", "INVALID_TOKEN", nil)
			c.Abort()
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) auth.Claims {
	value, _ := c.Get(claimsKey)
	claims, _ := value.(auth.Claims)
	return claims
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		claims := GetClaims(c)
		for _, role := range claims.Roles {
			if _, ok := allowed[role]; ok {
				c.Next()
				return
			}
		}
		response.Error(c, http.StatusForbidden, "insufficient permissions", "FORBIDDEN", nil)
		c.Abort()
	}
}

func RateLimit(namespace string, limit int, perMinutes int) gin.HandlerFunc {
	window := time.Duration(perMinutes) * time.Minute
	return func(c *gin.Context) {
		key := namespace + ":" + c.ClientIP()
		now := time.Now()
		rateLimitMu.Lock()
		bucket, ok := rateLimits[key]
		if !ok || now.After(bucket.ResetAt) {
			bucket = &rateBucket{Remaining: limit, ResetAt: now.Add(window)}
			rateLimits[key] = bucket
		}
		if bucket.Remaining <= 0 {
			rateLimitMu.Unlock()
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded", "RATE_LIMITED", nil)
			c.Abort()
			return
		}
		bucket.Remaining--
		rateLimitMu.Unlock()
		c.Next()
	}
}
