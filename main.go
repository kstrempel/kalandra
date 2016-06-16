package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
)

func main() {
	router := NewRouter()
	log.SetLevel(log.DebugLevel)
	log.Fatal(http.ListenAndServe(":8080", router))
}
