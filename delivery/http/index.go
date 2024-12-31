package httpapp

import (
	sendblobhandle "app/cmd/blob-service/delivery/http/send_blob"
	"app/internal/connection"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Register() http.Handler {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	sendblobhandle.Register(r)

	log.Printf(
		"Server h-learning-blob starting success! URL: http://%s:%s",
		connection.GetConnect().BlobService.Host,
		connection.GetConnect().BlobService.Port,
	)

	return r
}
