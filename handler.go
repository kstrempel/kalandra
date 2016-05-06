package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// AttributeEntry the json api attribute entry for a query
type AttributeEntry struct {
	Query string
}

// DataEntry the JSON:API data entry dict
type DataEntry struct {
	TypeDef    string `json:"type"`
	Keyspace   string
	Attributes AttributeEntry
}

// JSONAPI type of every JSONAPI structure
type JSONAPI struct {
	Data DataEntry
}

// ResultMeta in JSON:API definition
type ResultMeta struct {
	Query string `json:"query,omitempty"`
	Time  int64
}

// Error structure
type Error struct {
	Title  string
	Detail string
}

// Result JSON:API style
type Result struct {
	Data   []map[string]interface{} `json:"data,omitempty"`
	Errors Error
	Meta   ResultMeta
}

// Query handler
func Query(w http.ResponseWriter, r *http.Request) {
	var query JSONAPI
	var jsonAPI Result

	// parse the body
	start := time.Now()
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &query)
	if err != nil {
		jsonAPI.Errors = Error{
			Title:  "json error",
			Detail: err.Error()}
	} else {
		// get session to cassandra
		session := GetSession(query.Data.Keyspace)

		// query the database
		iter := session.Query(query.Data.Attributes.Query).Iter()
		sliceMap, err := iter.SliceMap()
		if err != nil {
			jsonAPI.Errors = Error{
				Title:  "CQL Query Error",
				Detail: err.Error()}
		} else {
			jsonAPI.Data = sliceMap
		}
	}

	// add meta data
	timeNeeded := time.Now().Sub(start)
	jsonAPI.Meta = ResultMeta{
		Query: query.Data.Attributes.Query,
		Time:  timeNeeded.Nanoseconds()}

	jsonResult, err := json.Marshal(jsonAPI)
	fmt.Fprintln(w, string(jsonResult))
}
