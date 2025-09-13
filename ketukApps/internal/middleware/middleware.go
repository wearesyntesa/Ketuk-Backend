package middleware

import (
	"fmt"
	"log"
	"time"

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
