package middleware

import (
	"net/http"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/dependencies"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// Middleware defines a function to process middleware
type Middleware func(*gin.Engine)

// setupMiddlewares adds required middlewares to the Gin router
func SetupMiddlewares(r *gin.Engine, logger *zap.Logger, frontendURL string) {
	middlewares := []Middleware{
		injectLogger(logger),
		addCORS(frontendURL),
	}
	r.Use(HealthMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(RequestLogger(logger))

	for _, middleware := range middlewares {
		middleware(r)
	}
}

// injectLogger returns a middleware to inject logger into Gin context
func injectLogger(logger *zap.Logger) Middleware {
	return func(r *gin.Engine) {
		r.Use(func(c *gin.Context) {
			c.Set("logger", logger.With(
				zap.String("request_id", c.Request.Header.Get("X-Request-ID")), // Add request ID context
			))
			c.Next()
		})
	}
}


func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Log request and response details
		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("request_id", c.GetString("request_id")),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// addCORS returns a middleware to add CORS with custom logic to allow the frontend
func addCORS(frontendURL string) Middleware {
	return func(r *gin.Engine) {
		corsMiddleware := cors.New(cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return origin == frontendURL
			},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:           int(12 * time.Hour / time.Second),
		})
		r.Use(func(c *gin.Context) {
			corsMiddleware.HandlerFunc(c.Writer, c.Request)
			c.Next()
		})
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String() // Generate a new UUID
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID) // Echo back to client
		c.Next()
	}
}

// Middleware to block requests when unhealthy
func HealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !dependencies.IsMongoDBHealthy {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service is temporarily unavailable"})
			c.Abort()
			return
		}
		c.Next()
	}
}


func ErrorLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request
		c.Next()

		// Log errors if any occurred
		err := c.Errors.Last()
		if err != nil {
			logger.Error("Request failed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
				zap.String("error", err.Err.Error()),
			)
		}
	}
}
