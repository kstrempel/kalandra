package main

import (
	"bytes"
	"encoding/json"
	"github.com/gocql/gocql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setUp() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "dev"
	session, _ := cluster.CreateSession()
	defer session.Close()

	session.Query("CREATE KEYSPACE mykeyspace WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };").Exec()
	session.Query("CREATE TABLE mykeyspace.users (user_id int PRIMARY KEY, fname text, lname text);").Exec()
	session.Query("INSERT INTO mykeyspace.users (user_id,  fname, lname) VALUES (1745, 'john', 'smith');").Exec()
}

func tearDown() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = ""
	session, _ := cluster.CreateSession()
	defer session.Close()

	session.Query("drop keyspace mykeyspace;").Exec()
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestEmptyArrayQuery(t *testing.T) {
	var query = []byte(`{}`)
	req, _ := http.NewRequest("POST", "/query", bytes.NewBuffer(query))
	queryHandler := NewRouter()
	w := httptest.NewRecorder()
	queryHandler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("query returned with wrong status code")
	}
}

func TestBrokenJsonQuery(t *testing.T) {
	var query = []byte(`{{}}`)
	req, _ := http.NewRequest("POST", "/query", bytes.NewBuffer(query))
	queryHandler := NewRouter()
	w := httptest.NewRecorder()
	queryHandler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("query returned with wrong status code")
	}
}

func TestQueryJsonQuery(t *testing.T) {
	var query = []byte(`
   {
   "data": {
        "type": "query",
        "keyspace": "mykeyspace",
        "attributes": {
            "query": "select * from users;"
        }
    }
  }`)
	req, _ := http.NewRequest("POST", "/query", bytes.NewBuffer(query))
	queryHandler := NewRouter()
	w := httptest.NewRecorder()
	queryHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("home page returned with wrong status code")
	}
	var result Result
	body, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(body, &result)
	data := result.Data
	if len(data) != 1 {
		t.Error("data size not 1")
	}
}
