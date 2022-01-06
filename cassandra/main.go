package main

import "github.com/gocql/gocql"

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = ""
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.Query()
}
