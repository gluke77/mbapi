package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gluke77/mbapi/srv"
)

func main() {
	apiKey := os.Getenv("MBAPI_KEY")

	if apiKey == "" {
		log.Fatalln("MBAPI_KEY environment variable is not set")
	}

	queueSizeStr := os.Getenv("MBAPI_QUEUE_SIZE")
	queueSize, err := strconv.Atoi(queueSizeStr)

	if err != nil || queueSize < 0 {
		queueSize = 10
	}

	log.Printf("Using queue size %d", queueSize)

	s := srv.New(apiKey, queueSize)

	port := os.Getenv("MBAPI_PORT")
	if port == "" {
		port = "8888"
	}
	log.Print("Listening at port " + port)

	log.Fatalln(http.ListenAndServe(":"+port, s))
}
