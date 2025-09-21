package guards

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Guard interface {
	CanActivate(ctx *gin.Context) bool
}

type AuthGuard struct {
	secretKey string
}

func NewAuthGuard(secretKey string) *AuthGuard {
	return &AuthGuard{secretKey: secretKey}
}

func (g *AuthGuard) CanActivate(ctx *gin.Context) bool {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		ctx.Abort()
		return false
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		ctx.Abort()
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if !g.validateToken(token) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return false
	}

	return true
}

func (g *AuthGuard) validateToken(token string) bool {
	return token != ""
}

type RoleGuard struct {
	requiredRoles []string
}

func NewRoleGuard(roles ...string) *RoleGuard {
	return &RoleGuard{requiredRoles: roles}
}

func (g *RoleGuard) CanActivate(ctx *gin.Context) bool {
	userRoles, exists := ctx.Get("user_roles")
	if !exists {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User roles not found"})
		ctx.Abort()
		return false
	}

	roles, ok := userRoles.([]string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user roles format"})
		ctx.Abort()
		return false
	}

	for _, requiredRole := range g.requiredRoles {
		for _, userRole := range roles {
			if userRole == requiredRole {
				return true
			}
		}
	}

	ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	ctx.Abort()
	return false
}

func GuardMiddleware(guards ...Guard) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, guard := range guards {
			if !guard.CanActivate(ctx) {
				return
			}
		}
		ctx.Next()
	}
}

type ThrottleGuard struct {
	requests map[string]int
	limit    int
}

func NewThrottleGuard(limit int) *ThrottleGuard {
	return &ThrottleGuard{
		requests: make(map[string]int),
		limit:    limit,
	}
}

func (g *ThrottleGuard) CanActivate(ctx *gin.Context) bool {
	clientIP := ctx.ClientIP()
	g.requests[clientIP]++

	if g.requests[clientIP] > g.limit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		ctx.Abort()
		return false
	}

	return true
}