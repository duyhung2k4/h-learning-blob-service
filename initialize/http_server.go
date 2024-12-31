package initialize

import (
	httpapp "app/cmd/blob-service/delivery/http"
	"app/internal/connection"
	"fmt"
	"log"
	"net/http"
	"time"
)

func runHttpSrver() {
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
			connection.GetConnectionInfo().BlobService.Host,
			connection.GetConnectionInfo().BlobService.Port,
		),
		Handler:      httpapp.Register(),
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	log.Fatalln(s.ListenAndServe())
}
