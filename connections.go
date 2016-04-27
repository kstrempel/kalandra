package main

import (
	"github.com/gocql/gocql"
)

// GetSession returns a open session from a open connection using the
// singleton pattern to reuses open sessions
func GetSession(keyspace string) *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	session, _ := cluster.CreateSession()

	return session
}
