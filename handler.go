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

// ResultMeta
type ResultMeta struct {
	Query string
	Time  int64
}

// Error
type Error struct {
	Title  string
	Detail string
}

// Result jsonapi style
type Result struct {
	Data   []map[string]interface{} `json:"data,omitempty"`
	Errors Error
	Meta   ResultMeta
}

// Query handler
func Query(w http.ResponseWriter, r *http.Request) {
	// parse the body
	var query JsonApi
	start := time.Now()
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &query)
	if err != nil {
		// TODO: Add better error handling
		fmt.Fprintln(w, "Upps error")
	}

	// get session to cassandra
	session := GetSession(query.Data.Keyspace)

	// query the database
	var jsonAPI Result
	iter := session.Query(query.Data.Attributes.Query).Iter()
	sliceMap, err := iter.SliceMap()
	if err != nil {
		jsonAPI = Result{
			Errors: Error{
				Title:  "CQL Query Error",
				Detail: err.Error()}}
	} else {
		var result []map[string]interface{}
		for _, row := range sliceMap {
			result = append(result, row)
		}
		jsonAPI = Result{
			Data: result,
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
