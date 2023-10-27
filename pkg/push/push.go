package push

import (
	"fmt"
	"hdruk/federated-metadata/pkg/routes"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run Runs our Push API service
func Run() bool {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("FMA_PORT")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Defines routes and handlers for REST interface
	router.GET("/ping", routes.PingHandler)
	router.POST("/test", routes.TestFederationHandler)
	router.POST("/federation", routes.CreateFederationHandler)
	router.DELETE("/federation", routes.DeleteFederationHandler)

	server.ListenAndServe()
	return true
}
