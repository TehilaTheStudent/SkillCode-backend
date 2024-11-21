package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// Middleware defines a function to process middleware
type Middleware func(*gin.Engine)

func SetupMiddlewares(r *gin.Engine, logger *zap.Logger, frontendURLS []string) {
	// Global middlewares
	r.Use(RequestIDMiddleware())          // Add unique request ID
	r.Use(RequestLogger(logger))          // Log request details
	r.Use(ErrorLoggingMiddleware(logger)) // Log errors after processing

	// CORS middleware
	r.Use(func(c *gin.Context) {
		corsMiddleware := cors.New(cors.Options{
			AllowOriginFunc: func(origin string) bool {
				for _, url := range frontendURLS {
					if origin == url {
						// Origin is allowed, no warning needed
						return true
					}
				}
	
				// Log and block the request if the origin is not allowed
				logger.Warn("CORS blocked request",
					zap.String("origin", origin),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				return false
			},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:           int(12 * time.Hour / time.Second),
		})
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})

	// Inject logger into context
	r.Use(func(c *gin.Context) {
		c.Set("logger", logger.With(
			zap.String("request_id", c.GetString("request_id")),
		))
		c.Next()
	})
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
				zap.String("request_id", c.GetString("request_id")), // Log request ID
				zap.String("error", err.Err.Error()),
			)
		}
	}
}
