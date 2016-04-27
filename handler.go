package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AttributeEntry the json api attribute entry for a query
type AttributeEntry struct {
	Query string
}

// DataEntry the jsonapi data entry dict
type DataEntry struct {
	TypeDef    string `json:"type"`
	Keyspace   string
	Attributes AttributeEntry
}

// JsonApi type of every JsonApi structure
type JsonApi struct {
	Data DataEntry
}

// Query handler
func Query(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var query JsonApi
	err := json.Unmarshal(body, &query)
	if err != nil {
		// TODO: Add better error handling
		fmt.Fprintln(w, "Upps error")
	}
	session := GetSession(query.Data.Keyspace)
	iter := session.Query(query.Data.Attributes.Query).Iter()
	defer iter.Close()
	fmt.Fprintln(w, iter.NumRows())
}
