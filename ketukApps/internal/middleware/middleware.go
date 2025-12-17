package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"ketukApps/internal/database"
	"ketukApps/internal/models"
	"ketukApps/internal/scheduler"
	"ketukApps/internal/utils"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// CORS middleware handles Cross-Origin Resource Sharing
// Configured for development - allows all origins
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow all origins for development
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow comprehensive headers for development
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, "+
				"Authorization, accept, origin, Cache-Control, X-Requested-With, "+
				"X-HTTP-Method-Override, Accept, Accept-Language, Content-Language")

		// Allow all common HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler middleware handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(500, gin.H{
					"success": false,
					"message": "Internal server error",
					"error":   "Something went wrong on our end",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// AuthRequired middleware validates JWT token and extracts user information
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user from database to ensure user still exists and get latest role
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "User not found",
			})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		c.Set("user", user)

		c.Next()
	}
}

// RequireRole middleware checks if the authenticated user has one of the required roles
// Must be used after AuthRequired middleware
// Usage: router.GET("/endpoint", middleware.AuthRequired(), middleware.RequireRole("admin"), handler)
// Usage: router.GET("/endpoint", middleware.AuthRequired(), middleware.RequireRole("admin", "user"), handler)
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthRequired middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Forbidden",
				Error:   "You do not have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a shorthand for RequireRole("admin")
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// IsOwnerOrAdmin checks if the user is the owner of a resource or an admin
// userIDParam is the URL parameter name containing the resource owner's ID
func IsOwnerOrAdmin(userIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user ID from context
		currentUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "User role not found",
			})
			c.Abort()
			return
		}

		// Admin can access everything
		if userRole == "admin" {
			c.Next()
			return
		}

		// Check if user is the owner
		resourceUserID := c.Param(userIDParam)
		if fmt.Sprintf("%v", currentUserID) != resourceUserID {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Forbidden",
				Error:   "You can only access your own resources",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Check State of unblock In current system
func CheckUnblockState() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !scheduler.IsUnblockEnabled() {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "This feature is currently disabled",
				Error:   "Please contact the administrator",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckUnblockStateReverseTechnique() gin.HandlerFunc {
	return func(c *gin.Context) {
		if scheduler.IsUnblockEnabled() {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "This feature is currently enabled, cant",
				Error:   "Please contact the administrator",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
