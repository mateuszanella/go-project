package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = ":80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s", webPort),
		Handler: routes(&app),
	}

	error := server.ListenAndServe()
	if error != nil {
		log.Fatal(error)
	}

}
