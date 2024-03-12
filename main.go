package main

import (
	"blog/db"
	"blog/transport"
	_ "github.com/lib/pq"
)

func main() {
	db.InitDB()
	transport.RoutingHandler()
}
