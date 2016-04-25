package main

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/query", Query)

	return router
}
