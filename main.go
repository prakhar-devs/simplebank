package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/prakhar-devs/simplebank/api"
	db "github.com/prakhar-devs/simplebank/db/sqlc"
	"github.com/prakhar-devs/simplebank/util"
)

func main() {
	// we used "." because the file where we are loading config and the app.env file are present in same directory
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	// Creaate a store using DB connection
	store := db.NewStore(conn)

	// Create a server using store
	server := api.NewServer(store)

	// Start the server
	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
