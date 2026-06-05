package server

import (
	"net/http"
	"os"
	"time"

	"backend-go/internal/domain/activities"
	"backend-go/internal/logger"
	"backend-go/internal/observability"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(otelgin.Middleware(observability.ServiceName()))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"*"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false, // Enable cookies/auth
	}))

	cwLogger := logger.MustGetInstance(&logger.LoggerConfig{
		Region:        os.Getenv("AWS_DEFAULT_REGION"),
		LogGroupName:  "/buzzer/backend-go",
		LogStreamName: "buzzer-" + time.Now().Format("2006-01-02"),
	})

	// r.Use(logger.GinCloudWatchMiddleware(cwLogger))

	api := r.Group("/api")
	activities.RegisterRoutes(api.Group("/activities"), cwLogger)
	//users.RegisterRoutes(api.Group("/users"), cwLogger)

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.HealthHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
