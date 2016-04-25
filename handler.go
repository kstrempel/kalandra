package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Query handler
func Query(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintln(w, string(body))
}
