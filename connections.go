package main

import (
	"github.com/gocql/gocql"
	"sync"
)

var sessions map[string]*gocql.Session
var mutex sync.Mutex

// GetSession returns a open session from a open connection using the
// singleton pattern to reuses open sessions
func GetSession(keyspace string) *gocql.Session {
	mutex.Lock()
	defer mutex.Unlock()

	session, ok := sessions[keyspace]
	if ok == false {
		if sessions == nil {
			sessions = make(map[string]*gocql.Session)
		}
		cluster := gocql.NewCluster("127.0.0.1")
		cluster.Keyspace = keyspace
		session, _ := cluster.CreateSession()
		sessions[keyspace] = session
		return session
	}

	return session
}
