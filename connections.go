package main

import (
	"github.com/gocql/gocql"
)

type CachedSession struct {
	session *gocql.Session
}

var sessions *CachedSession

// GetSession returns a open session from a open connection using the
// singleton pattern to reuses open sessions
func GetSession(keyspace string) *gocql.Session {
	// #TODO improve this for multiply keyspaces
	if sessions == nil {
		cluster := gocql.NewCluster("127.0.0.1")
		cluster.Keyspace = keyspace
		session, _ := cluster.CreateSession()
		sessions = &CachedSession{session: session}
	}
	return sessions.session
}
