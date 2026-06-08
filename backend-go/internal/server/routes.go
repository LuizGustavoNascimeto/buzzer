package server

import (
	"log"
	"net/http"
	"os"
	"time"

	activityHandler "backend-go/internal/services/activity/http"
	userHandler "backend-go/internal/services/user/http"

	"backend-go/pkg/auth"
	"backend-go/pkg/logger"
	"backend-go/pkg/observability"

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
	validator, err := auth.New(os.Getenv("AWS_DEFAULT_REGION"), os.Getenv("AWS_USER_POOLS_ID"))
	if err != nil {
		log.Fatal(err)
	}

	// r.Use(logger.GinCloudWatchMiddleware(cwLogger))

	api := r.Group("/api")
	activityHandler.RegisterRoutes(api.Group("/activities"), cwLogger, validator)
	userHandler.RegisterRoutes(api.Group("/users"), cwLogger, validator)
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
