package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gluke77/mbapi/srv"
)

func main() {
	apiKey := os.Getenv("MBAPI_KEY")

	if apiKey == "" {
		log.Fatalln("MBAPI_KEY environment variable is not set")
	}

	port := os.Getenv("MBAPI_PORT")
	if port == "" {
		port = "8888"
	}

	s := srv.New(apiKey)
	log.Print("Listening at port " + port)
	log.Fatalln(http.ListenAndServe(":"+port, s))
}
