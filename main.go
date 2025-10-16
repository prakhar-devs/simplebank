package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/prakhar-devs/simplebank/api"
	db "github.com/prakhar-devs/simplebank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	// Creaate a store using DB connection
	store := db.NewStore(conn)

	// Create a server using store
	server := api.NewServer(store)

	// Start the server
	err = server.StartServer(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
