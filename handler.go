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

// Result jsonapi style
type Result struct {
	Data []map[string]interface{}
	Meta ResultMeta
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
	iter := session.Query(query.Data.Attributes.Query).Iter()
	sliceMap, _ := iter.SliceMap()
	var result []map[string]interface{}
	for _, row := range sliceMap {
		result = append(result, row)
	}

	// build and send the answer back
	timeNeeded := time.Now().Sub(start)
	jsonApi := Result{
		Data: result,
		Meta: ResultMeta{
			Query: query.Data.Attributes.Query,
			Time:  timeNeeded.Nanoseconds()}}
	jsonResult, err := json.Marshal(jsonApi)
	fmt.Fprintln(w, string(jsonResult))
}
