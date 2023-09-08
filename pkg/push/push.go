package push

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Run Runs our Push API service
func Run() bool {
	finished := make(chan bool)

	router := gin.Default()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("FMA_PORT")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.ListenAndServe()

	finished <- true
	return true
}
